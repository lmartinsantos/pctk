package pctk

import (
	"image"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// MousePosition returns the current mouse position in the screen.
func (a *App) MousePosition() Position {
	return positionFromRaylib(
		rl.GetScreenToWorld2D(rl.GetMousePosition(), a.cam),
	)
}

// MouseIsInto returns true if the mouse is into the given region.
func (a *App) MouseIsInto(rect Rectangle) bool {
	return rl.CheckCollisionPointRec(a.MousePosition().toRaylib(), rect.toRaylib())
}

func (a *App) initMouse() {
	a.cursorTx = rl.LoadTextureFromImage(
		rl.NewImage(mouseCursorData(), 15, 15, 1, rl.UncompressedR8g8b8a8),
	)
	a.cursorColor = rl.NewColor(0xAA, 0xAA, 0xAA, 0xFF)
	rl.HideCursor()
}

func (a *App) drawMouseCursor() {
	if rl.IsCursorOnScreen() {
		pos := a.MousePosition()
		rl.DrawTexture(a.cursorTx, int32(pos.X-7), int32(pos.Y-7), a.cursorColor)
		a.cursorColor.R = max(0xAA, a.cursorColor.R+6)
		a.cursorColor.G = max(0xAA, a.cursorColor.G+6)
		a.cursorColor.B = max(0xAA, a.cursorColor.B+6)
	}
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
