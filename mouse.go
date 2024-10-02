package pctk

import (
	"image"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// MouseCursor is a custom mouse cursor.
type MouseCursor struct {
	Enabled bool

	cam *rl.Camera2D
	col Color
	tx  rl.Texture2D
}

// NewMouseCursor creates a new mouse cursor.
func NewMouseCursor(cam *rl.Camera2D) *MouseCursor {
	return &MouseCursor{
		cam: cam,
		col: rl.NewColor(0xAA, 0xAA, 0xAA, 0xFF),
		tx: rl.LoadTextureFromImage(
			rl.NewImage(mouseCursorData(), 15, 15, 1, rl.UncompressedR8g8b8a8),
		),
	}
}

// Draw renders the mouse cursor at the current position.
func (m *MouseCursor) Draw() {
	pos := m.Position()
	if m.Enabled && rl.IsCursorOnScreen() {
		rl.DrawTexture(m.tx, int32(pos.X-7), int32(pos.Y-7), m.col)
		m.col.R = max(0xAA, m.col.R+6)
		m.col.G = max(0xAA, m.col.G+6)
		m.col.B = max(0xAA, m.col.B+6)
	}
}

// OnScreen returns true if the mouse is on the screen.
func (m *MouseCursor) OnScreen() bool {
	return rl.IsCursorOnScreen()
}

// Position returns the current mouse position in the screen.
func (m *MouseCursor) Position() Position {
	if !m.Enabled {
		return Position{-1, -1}
	}
	return positionFromRaylib(
		rl.GetScreenToWorld2D(rl.GetMousePosition(), *m.cam),
	)
}

// MouseIsInto returns true if the mouse is into the given region.
func (m *MouseCursor) IsInto(rect Rectangle) bool {
	return rl.CheckCollisionPointRec(m.Position().toRaylib(), rect.toRaylib())
}

func mouseCursorData() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 15, 15))
	for i := 0; i <= 5; i++ {
		img.Set(i, 7, White)
		img.Set(7, i, White)
	}
	for i := 9; i <= 15; i++ {
		img.Set(i, 7, White)
		img.Set(7, i, White)
	}
	return img.Pix
}
