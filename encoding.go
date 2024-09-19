package pctk

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"io"
)

const (
	// ResourceIndexMagic is the magic number for the resource index.
	ResourceIndexMagic = "PCTK:IDX"

	// ResourceDataMagic is the magic number for the resource data.
	ResourceDataMagic = "PCTK:DAT"

	// ResourceFormatVersion
	ResourceFormatVersion uint16 = 0x0001
)

// BinaryEncode encodes objects to a writer using the binary format. If the object implements the
// BinaryEncoder interface, it will be used to encode the object. If the object is a string, it will
// be encoded with its bytes prefixed with a word indicating its size. Otherwise, the object will be
// encoded using the binary package assuming it has a fixed size.
//
// Use this function to encode only values that implement BinaryEncoder, are strings, or are
// fixed-size primitive types. Passing a custom struct might work, but it will fail if it has
// variable-size fields.
func BinaryEncode(w io.Writer, o ...any) (n int, err error) {
	for _, obj := range o {
		nn, err := binaryEncode(w, obj)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func binaryEncode(w io.Writer, obj any) (n int, err error) {
	if m, ok := obj.(BinaryEncoder); ok {
		return m.BinaryEncode(w)
	}
	if str, ok := obj.(string); ok {
		return BinaryEncode(w, uint16(len(str)), []byte(str))
	}
	err = binary.Write(w, binary.LittleEndian, obj)
	n = int(binary.Size(obj))
	return
}

// BinaryEncoder is a value that can encode itself to a binary format.
type BinaryEncoder interface {
	BinaryEncode(w io.Writer) (int, error)
}

// ResourceEncoder is a value that can encode resources to a writer.
type ResourceEncoder struct {
	index io.Writer
	data  io.Writer

	next int
}

// NewResourceEncoder creates a new resource encoder that writes the index and data to the given
// writers.
func NewResourceEncoder(index, data io.Writer) (*ResourceEncoder, error) {
	enc := &ResourceEncoder{
		index: index,
		data:  data,
	}
	err := enc.encodeHeaders()
	return enc, err
}

// EncodeCostume encodes a costume using the resource encoder.
func (e *ResourceEncoder) EncodeCostume(loc ResourceLocator, c *Costume, comp ResourceCompression) error {
	return e.encodeResource(loc, c, resourceHeader{
		Type:        ResourceTypeCostume,
		Compression: comp,
	})
}

// EncodeMusic encodes a music using the resource encoder.
func (e *ResourceEncoder) EncodeMusic(loc ResourceLocator, m *Music, comp ResourceCompression) error {
	return e.encodeResource(loc, m, resourceHeader{
		Type:        ResourceTypeMusic,
		Compression: comp,
	})
}

// EncodeScene encodes a scene using the resource encoder.
func (e *ResourceEncoder) EncodeScene(loc ResourceLocator, s *Scene, comp ResourceCompression) error {
	return e.encodeResource(loc, s, resourceHeader{
		Type:        ResourceTypeScene,
		Compression: comp,
	})
}

// EncodeScript encodes a script using the resource encoder.
func (e *ResourceEncoder) EncodeScript(loc ResourceLocator, s *Script, comp ResourceCompression) error {
	return e.encodeResource(loc, s, resourceHeader{
		Type:        ResourceTypeScript,
		Compression: comp,
	})
}

// EncodeSound encodes a sound using the resource encoder.
func (e *ResourceEncoder) EncodeSound(loc ResourceLocator, s *Sound, comp ResourceCompression) error {
	return e.encodeResource(loc, s, resourceHeader{
		Type:        ResourceTypeSound,
		Compression: comp,
	})
}

func (e *ResourceEncoder) encodeResource(loc ResourceLocator, res BinaryEncoder, h resourceHeader) error {
	var n int
	var err error

	if n, err = h.BinaryEncode(e.data); err != nil {
		return err
	}

	switch h.Compression {
	case CompressionNone:
		nn, err := res.BinaryEncode(e.data)
		if err != nil {
			return err
		}
		n += nn
	case CompressionGzip:
		var buf bytes.Buffer
		zipper := gzip.NewWriter(&buf)
		if _, err = res.BinaryEncode(zipper); err != nil {
			return err
		}
		if err = zipper.Close(); err != nil {
			return err
		}
		nn, err := e.data.Write(buf.Bytes())
		if err != nil {
			return err
		}
		n += nn
	}

	if err := e.encodeIndexEntry(loc, e.next, n); err != nil {
		return err
	}
	e.next += n
	return nil
}

func (e *ResourceEncoder) encodeHeaders() error {
	if err := e.encodeIndexHeader(); err != nil {
		return err
	}
	if err := e.encodeDataHeader(); err != nil {
		return err
	}
	return nil
}

func (e *ResourceEncoder) encodeIndexHeader() error {
	_, err := BinaryEncode(e.index,
		[]byte(ResourceIndexMagic),
		ResourceFormatVersion,
	)
	return err
}

func (e *ResourceEncoder) encodeDataHeader() error {
	n, err := BinaryEncode(e.data,
		[]byte(ResourceDataMagic),
		ResourceFormatVersion,
	)
	e.next += n
	return err
}

func (e *ResourceEncoder) encodeIndexEntry(loc ResourceLocator, offset, size int) error {
	_, err := BinaryEncode(e.index, loc, uint32(offset), uint32(size))
	return err
}

type resourceHeader struct {
	Type        ResourceType
	Compression ResourceCompression
}

func (h resourceHeader) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, h.Type, h.Compression, [14]byte{})
}

// ResourceType is the type of a resource.
type ResourceType byte

const (
	ResourceTypeUndefined ResourceType = iota
	ResourceTypeCostume
	ResourceTypeMusic
	ResourceTypeScene
	ResourceTypeScript
	ResourceTypeSound
)

// ResourceCompression is the type of compression used for resources while encoding.
type ResourceCompression byte

const (
	// CompressionNone is the no compression.
	CompressionNone ResourceCompression = iota

	// CompressionGzip is the gzip compression.
	CompressionGzip
)
