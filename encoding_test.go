package pctk_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/apoloval/pctk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryEncodeDecode(t *testing.T) {
	var buf bytes.Buffer
	n, err := pctk.BinaryEncode(&buf, pctk.ResourceID("foo/bar"), uint64(42))
	require.NoError(t, err)
	require.Positive(t, n)

	var ref pctk.ResourceID
	var num uint64
	err = pctk.BinaryDecode(&buf, &ref, &num)
	require.NoError(t, err)

	assert.Equal(t, pctk.ResourceID("foo/bar"), ref)
	assert.Equal(t, uint64(42), num)
}

func TestResourceEncoder_Empty(t *testing.T) {
	var idxBuf bytes.Buffer
	var datBuf bytes.Buffer
	_, err := pctk.NewResourceEncoder(&idxBuf, &datBuf)
	idx := idxBuf.Bytes()
	dat := datBuf.Bytes()

	require.NoError(t, err)
	assert.Equal(t, []byte("PCTK:IDX"), idx[0:8])                         // magic
	assert.Equal(t, byte(pctk.ResourceFormatVersion&0x00FF), idx[8])      // version low byte
	assert.Equal(t, byte((pctk.ResourceFormatVersion&0xFF00)>>8), idx[9]) // version high byte
	assert.Len(t, idx, 10)

	assert.Equal(t, []byte("PCTK:DAT"), dat[0:8])                         // magic
	assert.Equal(t, byte(pctk.ResourceFormatVersion&0x00FF), dat[8])      // version low byte
	assert.Equal(t, byte((pctk.ResourceFormatVersion&0xFF00)>>8), dat[9]) // version high byte
	assert.Len(t, dat, 10)
}

func TestResourceEncoder_EncodeResource(t *testing.T) {
	var idxBuf bytes.Buffer
	var datBuf bytes.Buffer
	enc, err := pctk.NewResourceEncoder(&idxBuf, &datBuf)
	require.NoError(t, err)

	err = enc.EncodeScript(
		pctk.ResourceID("hello"),
		&pctk.Script{
			Language: pctk.ScriptLua,
			Code:     []byte("print('Hello, world!')"),
		},
		pctk.CompressionNone,
	)
	idx := idxBuf.Bytes()
	dat := datBuf.Bytes()

	require.NoError(t, err)
	assert.Equal(t, []byte("\x05\x00hello"), idx[0x0A:0x11])                  // idx entry ref
	assert.Equal(t, uint32(0x0A), binary.LittleEndian.Uint32(idx[0x11:0x15])) // idx entry offset
	assert.Equal(t, uint32(43), binary.LittleEndian.Uint32(idx[0x15:0x19]))   // idx entry size

	assert.Equal(t, append([]byte{
		0x04,                                           // resource type
		0x00,                                           // compression
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ...
		0x01,                   // script type
		0x16, 0x00, 0x00, 0x00, // script size
	}, []byte("print('Hello, world!')")...), dat[10:])
}

func TestResourceEncoder_EncodeResourceWithGZip(t *testing.T) {
	var idxBuf bytes.Buffer
	var datBuf bytes.Buffer
	enc, err := pctk.NewResourceEncoder(&idxBuf, &datBuf)
	require.NoError(t, err)

	err = enc.EncodeScript(
		pctk.ResourceID("hello"),
		&pctk.Script{
			Language: pctk.ScriptLua,
			Code:     []byte("print('Hello, world!')"),
		},
		pctk.CompressionGzip,
	)
	idx := idxBuf.Bytes()
	dat := datBuf.Bytes()

	require.NoError(t, err)
	assert.Equal(t, []byte("\x05\x00hello"), idx[0x0A:0x11])                  // idx entry ref
	assert.Equal(t, uint32(0x0A), binary.LittleEndian.Uint32(idx[0x11:0x15])) // idx entry offset
	assert.Equal(t, uint32(67), binary.LittleEndian.Uint32(idx[0x15:0x19]))   // idx entry size

	assert.Equal(t, byte(0x04), dat[0x0A])            // dat hd resource type
	assert.Equal(t, byte(0x01), dat[0x0B])            // dat hd compression
	assert.Equal(t, make([]byte, 14), dat[0x0C:0x1A]) // dat hd reserved
}
