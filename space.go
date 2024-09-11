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

// Position represents a 2D position.
type Position struct {
	X, Y int
}

// NewPos creates a new position.
func NewPos(x, y int) Position {
	return Position{x, y}
}

func positionFromRaylib(v rl.Vector2) Position {
	return Position{int(v.X), int(v.Y)}
}

// Add adds two positions.
func (p Position) Add(other Position) Position {
	return Position{p.X + other.X, p.Y + other.Y}
}

// Sub subtracts two positions.
func (p Position) Sub(other Position) Position {
	return Position{p.X - other.X, p.Y - other.Y}
}

func (p Position) toRaylib() rl.Vector2 {
	return rl.NewVector2(float32(p.X), float32(p.Y))
}

// Size represents a 2D size.
type Size struct {
	W, H uint
}

func (s Size) toRaylib() rl.Vector2 {
	return rl.NewVector2(float32(s.W), float32(s.H))
}

// Rectangle represents a 2D rectangle.
type Rectangle struct {
	Pos  Position
	Size Size
}

func NewRect(x, y int, w, h uint) Rectangle {
	return Rectangle{
		Pos:  Position{x, y},
		Size: Size{uint(w), uint(h)},
	}
}

func (r Rectangle) toRaylib() rl.Rectangle {
	return rl.NewRectangle(
		float32(r.Pos.X), float32(r.Pos.Y), float32(r.Size.W), float32(r.Size.H),
	)
}
