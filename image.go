package pctk

import (
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
	tex := rl.LoadTextureFromImage(raw)
	if !rl.IsTextureReady(tex) {
		log.Fatalf("Failed to load texture from image %s", path)
	}
	return &Image{raw, tex}
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
