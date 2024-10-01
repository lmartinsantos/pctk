package pctk

import (
	"io"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Image represents an graphic image.
type Image struct {
	raw *rl.Image
	tex rl.Texture2D
}

// LoadImageFromFile loads an image from a file.
func LoadImageFromFile(path string) *Image {
	raw := rl.LoadImage(path)
	if !rl.IsImageReady(raw) {
		log.Fatalf("Failed to load image from file %s", path)
	}
	return &Image{raw: raw}
}

// Texture returns the texture of the image. If the texture is not ready, it will be loaded.
func (i *Image) Texture() rl.Texture2D {
	if !rl.IsTextureReady(i.tex) {
		i.tex = rl.LoadTextureFromImage(i.raw)
	}
	return i.tex
}

// Release the resources used by the image.
func (i *Image) Release() {
	rl.UnloadTexture(i.tex)
	rl.UnloadImage(i.raw)
}

// Width returns the width of the image.
func (i *Image) Width() int32 {
	return i.raw.Width
}

// Height returns the height of the image.
func (i *Image) Height() int32 {
	return i.raw.Height
}

// BinaryEncode encodes the image to a binary format. The encoded format is:
// - [0..3] uint32: the length of the image bytes.
// - [4..n] []byte: the image bytes in PNG format.
func (i *Image) BinaryEncode(w io.Writer) (int, error) {
	bytes := rl.ExportImageToMemory(*i.raw, ".png")
	return BinaryEncode(w, uint32(len(bytes)), bytes)
}

// BinaryDecode decodes the image from a binary format. See Image.BinaryEncode for the format.
func (i *Image) BinaryDecode(r io.Reader) error {
	var length uint32
	if err := BinaryDecode(r, &length); err != nil {
		return err
	}

	bytes := make([]byte, length)
	if err := BinaryDecode(r, &bytes); err != nil {
		return err
	}

	i.raw = rl.LoadImageFromMemory(".png", bytes, int32(len(bytes)))
	return nil
}

// Draw the image in the image on the screen.
func (i *Image) Draw(pos Position, tint rl.Color) {
	if i == nil {
		return
	}
	rl.DrawTexture(i.Texture(), int32(pos.X), int32(pos.Y), tint)
}
