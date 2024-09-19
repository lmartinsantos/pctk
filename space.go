package pctk

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	// ScreenWidth is the width of the screen (ignoring zoom).
	ScreenWidth = 320

	// ScreenHeight is the height of the screen (ignoring zoom).
	ScreenHeight = 200

	// ViewportHeight is the width of the room section of the screen.
	ViewportHeight = 144

	// ControlBoxHeight is the width of the control box section of the screen.
	ControlBoxHeight = 56
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

func (p Position) String() string {
	return fmt.Sprintf("(X:%d, Y:%d)", p.X, p.Y)
}

// Add adds two positions.
func (p Position) Add(other Position) Position {
	return Position{p.X + other.X, p.Y + other.Y}
}

// Sub subtracts two positions.
func (p Position) Sub(other Position) Position {
	return Position{p.X - other.X, p.Y - other.Y}
}

// Above returns a position above the current one.
func (p Position) Above(h int) Position {
	return Position{p.X, p.Y - h}
}

// Distance returns the distance between two positions.
func (p Position) Distance(other Position) Size {
	s := Size{other.X - p.X, other.Y - p.Y}
	if s.W < 0 {
		s.W = -s.W
	}
	if s.H < 0 {
		s.H = -s.H
	}
	return s
}

// DirectionTo returns the direction from the current position to another.
func (p Position) DirectionTo(other Position) Direction {
	dist := p.Distance(other)
	if dist.W > dist.H {
		// Horizontal direction
		if other.X < p.X {
			return DirLeft
		}
		return DirRight
	} else {
		// Vertical direction
		if other.Y < p.Y {
			return DirUp
		}
		return DirDown
	}
}

func (p Position) toRaylib() rl.Vector2 {
	return rl.NewVector2(float32(p.X), float32(p.Y))
}

// Size represents a 2D size.
type Size struct {
	W, H int
}

func sizeFromRaylib(v rl.Vector2) Size {
	return Size{int(v.X), int(v.Y)}
}

func (s Size) String() string {
	return fmt.Sprintf("(W:%d, H:%d)", s.W, s.H)
}

func (s Size) FlipH() Size {
	return Size{W: -s.W, H: s.H}
}

func (s Size) toRaylib() rl.Vector2 {
	return rl.NewVector2(float32(s.W), float32(s.H))
}

// Rectangle represents a 2D rectangle.
type Rectangle struct {
	Pos  Position
	Size Size
}

func NewRect(x, y, w, h int) Rectangle {
	return Rectangle{
		Pos:  Position{x, y},
		Size: Size{w, h},
	}
}

func (r Rectangle) String() string {
	return fmt.Sprintf("(Pos:%v, Size:%v)", r.Pos, r.Size)
}

func (r Rectangle) toRaylib() rl.Rectangle {
	return rl.NewRectangle(
		float32(r.Pos.X), float32(r.Pos.Y), float32(r.Size.W), float32(r.Size.H),
	)
}

// Direction represents a direction in 2D space.
type Direction byte

func (d Direction) String() string {
	switch d {
	case DirUp:
		return "Up"
	case DirDown:
		return "Down"
	case DirLeft:
		return "Left"
	case DirRight:
		return "Right"
	default:
		return "Unknown"
	}
}

const (
	DirRight Direction = iota
	DirLeft
	DirUp
	DirDown
)
