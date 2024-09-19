package pctk

import (
	"io"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// SpriteSheet represents a collection of sprites arranged in a grid-shaped sheet.
type SpriteSheet struct {
	raw       *rl.Image
	tex       rl.Texture2D
	frameSize Size
}

// LoadSpriteSheetFromFile loads a sprite sheet from a image file.
func LoadSpriteSheetFromFile(path string, frameSize Size) *SpriteSheet {
	return &SpriteSheet{
		raw:       rl.LoadImage(path),
		frameSize: frameSize,
	}
}

// Release releases the resources used by the sprite sheet.
func (s *SpriteSheet) Release() {
	rl.UnloadTexture(s.tex)
}

// DrawSprite draws a sprite from the sprite sheet at the given position.
func (s *SpriteSheet) DrawSprite(i, j uint, pos Position, flip bool) {
	src := Rectangle{
		Pos: Position{
			int(s.frameSize.W) * int(i),
			int(s.frameSize.H) * int(j),
		},
		Size: s.frameSize,
	}
	if flip {
		src.Size = src.Size.FlipH()
	}
	rl.DrawTextureRec(s.texture(), src.toRaylib(), pos.toRaylib(), rl.White)
}

// BinaryEncode encodes the sprite sheet to a binary format. The encoded format is:
// - [0..3] uint32: the length of the image bytes.
// - [4..n] []byte: the image bytes in PNG format.
func (s *SpriteSheet) BinaryEncode(w io.Writer) (int, error) {
	bytes := rl.ExportImageToMemory(*s.raw, ".png")
	return BinaryEncode(w, uint32(len(bytes)), bytes)
}

func (s *SpriteSheet) texture() rl.Texture2D {
	if !rl.IsTextureReady(s.tex) {
		s.tex = rl.LoadTextureFromImage(s.raw)
	}
	return s.tex
}
