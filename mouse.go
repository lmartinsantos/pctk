package pctk

import rl "github.com/gen2brain/raylib-go/raylib"

// MousePosition returns the current mouse position in the screen.
func (a *App) MousePosition() ScreenPosition {
	return rl.GetScreenToWorld2D(rl.GetMousePosition(), a.cam)
}

// MouseIsInto returns true if the mouse is into the given region.
func (a *App) MouseIsInto(reg ScreenRegion) bool {
	return rl.CheckCollisionPointRec(a.MousePosition(), reg)
}
