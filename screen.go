package pctk

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	// ScreenWidth is the width of the screen (ignoring zoom).
	ScreenWidth = 320

	// ScreenHeight is the height of the screen (ignoring zoom).
	ScreenHeight = 200

	// ScreenHeightScene is the width of the scene section of the screen.
	ScreenHeightScene = 144

	// ScreenHeightControl is the width of the control box section of the screen.
	ScreenHeightControl = 56
)

// ScreenPosition is a 2D position on the screen.
type ScreenPosition = rl.Vector2

// ScreenRegion is a rectangular region on the screen.
type ScreenRegion = rl.Rectangle
