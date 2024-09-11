package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// SpriteSheet represents a collection of sprites arranged in a grid-shaped sheet.
type SpriteSheet struct {
	tex       rl.Texture2D
	frameSize Size
}

// LoadSpriteSheetFromFile loads a sprite sheet from a image file.
func LoadSpriteSheetFromFile(path string, frameSize Size) *SpriteSheet {
	return &SpriteSheet{
		tex:       rl.LoadTexture(path),
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
	rl.DrawTextureRec(s.tex, src.toRaylib(), pos.toRaylib(), rl.White)
}
