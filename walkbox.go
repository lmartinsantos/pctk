package pctk

import "log"

// Walkbox refers to a convex polygonal area that defines the walkable space for actors.
type WalkBox struct {
	walkBoxID string
	enabled   bool
	vertices  [4]*Positionf
}

// NewWalkBox creates a new WalkBox with the given ID and vertices.
// It ensures the polygon formed by the vertices is convex. If not, it will cause a panic.
// Why convex? Because you can draw a straight line/path between any two vertices inside the polygon
// without needing to implement complex pathfinding algorithms.
func NewWalkBox(id string, vertices [4]*Positionf) *WalkBox {
	w := &WalkBox{
		walkBoxID: id,
		vertices:  vertices,
		enabled:   true,
	}

	if !w.isConvex() {
		log.Panicf("walkbox must be a convex polygon: %v", vertices)
	}
	return w
}

// isConvex check if the current WalkBox is a convex poligon.
func (w *WalkBox) isConvex() bool {
	numVertices := len(w.vertices)

	var totalCrossProduct float32
	var polygonDirection bool // true if clockwise, false if counter-clockwise
	for i := 0; i < numVertices; i++ {
		// Get three consecutive vertices (cyclically)
		p1 := w.vertices[i]
		p2 := w.vertices[(i+1)%numVertices]
		p3 := w.vertices[(i+2)%numVertices]

		cp := p1.CrossProduct(p2, p3)

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
	w.enabled = enable
	return w
}

// ContainsPoint check if the provided position is in the boundaries defined by the WalkBox.
func (w *WalkBox) ContainsPoint(p *Positionf) bool {
	numberOfIntersections := 0
	numVertices := len(w.vertices)

	for i := 0; i < numVertices; i++ {
		p1 := w.vertices[i]
		p2 := w.vertices[(i+1)%numVertices]

		if p.IsIntersecting(p1, p2) {
			numberOfIntersections++
		}
	}

	return numberOfIntersections%2 == 1 // Odd count means inside
}

// IsConnected checks if WalkBoxes are connected.
func (w *WalkBox) IsConnected(otherWalkBox *WalkBox) bool {
	for _, vertex := range otherWalkBox.vertices {
		if w.ContainsPoint(vertex) {
			return true
		}
	}
	return false
}

// WalkBoxMatrix represents a collection of WalkBoxes and their adjacency relationships.
type WalkBoxMatrix struct {
	walkBoxes     map[string]*WalkBox
	adjacencyList map[string][]*WalkBox
}

// NewWalkBoxMatrix creates and returns a new WalkBoxMatrix instance
func NewWalkBoxMatrix() *WalkBoxMatrix {
	return &WalkBoxMatrix{
		walkBoxes:     make(map[string]*WalkBox),
		adjacencyList: make(map[string][]*WalkBox),
	}
}

// Add adds a new WalkBox to the WalkBoxMatrix.
// The new WalkBox must be connected to at least one existing WalkBox in the matrix
// ensuring that there are no isolated walkable areas.
func (wm *WalkBoxMatrix) Add(w *WalkBox) {
	if _, exists := wm.adjacencyList[w.walkBoxID]; !exists {
		wm.adjacencyList[w.walkBoxID] = []*WalkBox{}
	}

	for _, otherWalkBox := range wm.walkBoxes {
		if w.IsConnected(otherWalkBox) {
			wm.adjacencyList[w.walkBoxID] = append(wm.adjacencyList[w.walkBoxID], otherWalkBox)
			wm.adjacencyList[otherWalkBox.walkBoxID] = append(wm.adjacencyList[otherWalkBox.walkBoxID], w)
		}
	}

	wm.walkBoxes[w.walkBoxID] = w

	if len(wm.walkBoxes) > 1 && len(wm.adjacencyList[w.walkBoxID]) < 1 {
		log.Panicf("walkbox %s is not connected", w.walkBoxID)
	}
}

// Adjacents returns a slice of WalkBoxes that are adjacent to the WalkBox with the specified walkBoxID.
// If no adjacencies are found, it returns an empty slice.
func (wm *WalkBoxMatrix) Adjacents(walkBoxID string) []*WalkBox {
	if adjacents, ok := wm.adjacencyList[walkBoxID]; ok {
		return adjacents
	}
	return []*WalkBox{}
}
