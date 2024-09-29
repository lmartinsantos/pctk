package pctk

import "log"

const (
	// MinPolygonVertices defines the minimum number of vertices to form a polygon.
	MinPolygonVertices = 3
)

// Walkbox refers to a convex polygonal area that defines the walkable space for actors.
type WalkBox struct {
	WalkBoxID string
	Enabled   bool
	Vertices  []*Position
	// TODO: Scale float32
}

// NewWalkBox creates a new WalkBox with the given ID and vertices.
// It ensures the polygon formed by the vertices is convex. If not, it will cause a panic.
// Why convex? Because you can draw a straight line/path between any two vertices inside the polygon
// without needing to implement complex pathfinding algorithms.
func NewWalkBox(id string, vertices []*Position) *WalkBox {
	w := &WalkBox{
		WalkBoxID: id,
		Vertices:  vertices,
		Enabled:   true,
	}

	if !w.isConvex() {
		log.Panicf("walkbox must be a convex polygon: %v", vertices)
	}
	return w
}

// isConvex check if the current WalkBox is a convex poligon.
func (w *WalkBox) isConvex() bool {
	numVertices := len(w.Vertices)

	if numVertices < MinPolygonVertices {
		return false
	}

	var totalCrossProduct int
	var polygonDirection bool // true if clockwise, false if counter-clockwise
	for i := 0; i < numVertices; i++ {
		// Get three consecutive vertices (cyclically)
		p1 := w.Vertices[i]
		p2 := w.Vertices[(i+1)%numVertices]
		p3 := w.Vertices[(i+2)%numVertices]

		cp := crossProduct(p1, p2, p3)

		if cp == 0 {
			continue // Skip collinear vertices
		}

		totalCrossProduct += cp

		if i == 0 {
			polygonDirection = cp > 0
		} else {
			if (cp > 0) != polygonDirection {
				return false // If direction changes, the polygon is not convex
			}
		}
	}
	return totalCrossProduct != 0
}

// Enable sets the enabled state of the WalkBox.
func (w *WalkBox) Enable(enable bool) *WalkBox {
	w.Enabled = enable
	return w
}

// ContainsPoint check if the provided position is in the boundaries defined by the WalkBox.
func (w *WalkBox) ContainsPoint(p *Position) bool {
	numberOfIntersections := 0
	numVertices := len(w.Vertices)

	for i := 0; i < numVertices; i++ {
		p1 := w.Vertices[i]
		p2 := w.Vertices[(i+1)%numVertices]

		if isIntersecting(p, p1, p2) {
			numberOfIntersections++
		}
	}

	return numberOfIntersections%2 == 1 // Odd count means inside
}

// crossProduct calculates the cross product of the vectors formed by three consecutive vertices of a polygon.
func crossProduct(p1, p2, p3 *Position) int {
	return (p2.X-p1.X)*(p3.Y-p2.Y) - (p2.Y-p1.Y)*(p3.X-p2.X)
}

// isIntersecting calculate the x-coordinate of the intersection of the ray with the line segment (Ray-Casting method).
func isIntersecting(p *Position, p1, p2 *Position) bool {
	if (p1.Y > p.Y) != (p2.Y > p.Y) {
		return p.X < (p2.X-p1.X)*(p.Y-p1.Y)/(p2.Y-p1.Y)+p1.X
	}

	return false
}
