package pctk_test

import (
	"testing"

	"github.com/apoloval/pctk"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultWalkBoxID = "walkbox"
)

func TestNewWalkBox(t *testing.T) {
	testCases := []struct {
		name        string
		vertices    [4]*pctk.Positionf
		shouldPanic bool
		message     string
	}{
		{
			name:        "Concave polygon should panic",
			vertices:    [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 2, Y: 1}, {X: 4, Y: 4}},
			shouldPanic: true,
			message:     "Expected panic because vertices form a concave polygon!",
		},
		{
			name:        "Collinear vertices should panic",
			vertices:    [4]*pctk.Positionf{{X: 1, Y: 2}, {X: 2, Y: 4}, {X: 3, Y: 6}, {X: 4, Y: 8}},
			shouldPanic: true,
			message:     "Expected panic because vertices are collinear!",
		},
		{
			name:        "Should successfully create a valid WalkBox with a convex polygon",
			vertices:    [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			shouldPanic: false,
			message:     "Expected create a valid WalkBox, vertices form a convex polygon!",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.shouldPanic {
				assert.Panics(t, func() {
					pctk.NewWalkBox(DefaultWalkBoxID, testCase.vertices)
				}, testCase.message)
			} else {
				assert.NotPanics(t, func() {
					pctk.NewWalkBox(DefaultWalkBoxID, testCase.vertices)
				}, testCase.message)
			}
		})
	}

}

func TestContainsPoint(t *testing.T) {
	testCases := []struct {
		name       string
		vertices   [4]*pctk.Positionf
		point      *pctk.Positionf
		assertFunc func(t *testing.T, isInside bool)
	}{
		{
			name:     "The point should be considered inside the polygon when it is on the edge",
			vertices: [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Positionf{X: 2, Y: 0}, // On the edge
			assertFunc: func(t *testing.T, isInside bool) {
				assert.True(t, isInside)
			},
		},
		{
			name:     "The point should be inside the polygon",
			vertices: [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Positionf{X: 2, Y: 2},
			assertFunc: func(t *testing.T, isInside bool) {
				assert.True(t, isInside)
			},
		},
		{
			name:     "The point should be outside the polygon",
			vertices: [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Positionf{X: 5, Y: 5},
			assertFunc: func(t *testing.T, isInside bool) {
				assert.False(t, isInside)
			},
		},
		{
			name:     "The point should be considered inside the polygon when it is on a vertex",
			vertices: [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Positionf{X: 0, Y: 0}, // On the vertex
			assertFunc: func(t *testing.T, isInside bool) {
				assert.True(t, isInside)
			},
		},
		{
			name:     "The point should be outside when it is far from the polygon",
			vertices: [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Positionf{X: 10, Y: 10}, // Clearly outside the polygon
			assertFunc: func(t *testing.T, isInside bool) {
				assert.False(t, isInside)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			walkBox := pctk.NewWalkBox(DefaultWalkBoxID, testCase.vertices)
			isInside := walkBox.ContainsPoint(testCase.point)
			testCase.assertFunc(t, isInside)
		})
	}
}

func TestWalkBoxMatrixAdd(t *testing.T) {
	testCases := []struct {
		name        string
		walkBoxes   []*pctk.WalkBox
		shouldPanic bool
		message     string
	}{
		{
			name: "Add connected WalkBoxes",
			walkBoxes: []*pctk.WalkBox{
				pctk.NewWalkBox("WalkBoxA", [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}}),
				pctk.NewWalkBox("WalkBoxB", [4]*pctk.Positionf{{X: 1.5, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 2}, {X: 1.5, Y: 2}}),
			},
			shouldPanic: false,
			message:     "Expected WalkBoxes A and B to be connected without panic!",
		},
		{
			name: "Add disconnected WalkBox should panic",
			walkBoxes: []*pctk.WalkBox{
				pctk.NewWalkBox("WalkBoxA", [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}}),
				pctk.NewWalkBox("WalkBoxB", [4]*pctk.Positionf{{X: 5, Y: 5}, {X: 6, Y: 5}, {X: 6, Y: 6}, {X: 5, Y: 6}}),
			},
			shouldPanic: true,
			message:     "Expected panic because WalkBoxB is not connected!",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wm := pctk.NewWalkBoxMatrix()

			if testCase.shouldPanic {
				assert.Panics(t, func() {
					for _, wb := range testCase.walkBoxes {
						wm.Add(wb)
					}
				}, testCase.message)
			} else {
				assert.NotPanics(t, func() {
					for _, wb := range testCase.walkBoxes {
						wm.Add(wb)
					}
				}, testCase.message)

				// Verifying connections of WalkBoxA
				adjacents := wm.Adjacents("WalkBoxA")
				assert.True(t, len(adjacents) == 1, "WalkBoxA should be connected to 1 WalkBox")
			}
		})
	}
}

func TestWalkBoxMatrixAdjacents(t *testing.T) {
	wbA := pctk.NewWalkBox("WalkBoxA", [4]*pctk.Positionf{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 2}, {X: 0, Y: 2}})
	wbB := pctk.NewWalkBox("WalkBoxB", [4]*pctk.Positionf{{X: 1.5, Y: 0}, {X: 3.5, Y: 0}, {X: 3.5, Y: 2}, {X: 1.5, Y: 2}})
	wbC := pctk.NewWalkBox("WalkBoxC", [4]*pctk.Positionf{{X: 3, Y: 0}, {X: 5, Y: 0}, {X: 5, Y: 2}, {X: 3, Y: 2}})
	wbD := pctk.NewWalkBox("WalkBoxD", [4]*pctk.Positionf{{X: 4.5, Y: 0}, {X: 6.5, Y: 0}, {X: 6.5, Y: 2}, {X: 4.5, Y: 2}})

	wm := pctk.NewWalkBoxMatrix()

	wm.Add(wbA)
	wm.Add(wbB)
	wm.Add(wbC)
	wm.Add(wbD)

	t.Run("Check WalkBoxA adjacency", func(t *testing.T) {
		adjacents := wm.Adjacents("WalkBoxA")
		assert.Equal(t, 1, len(adjacents), "WalkBoxA should be adjacent to 1 WalkBox")
		assert.Contains(t, adjacents, wbB, "WalkBoxA should be adjacent to WalkBoxB")
	})

	t.Run("Check WalkBoxB adjacency", func(t *testing.T) {
		adjacents := wm.Adjacents("WalkBoxB")
		assert.Equal(t, 2, len(adjacents), "WalkBoxB should be adjacent to 2 WalkBoxes")
		assert.Contains(t, adjacents, wbA, "WalkBoxB should be adjacent to WalkBoxA")
		assert.Contains(t, adjacents, wbC, "WalkBoxB should be adjacent to WalkBoxC")
	})

	t.Run("Check WalkBoxC adjacency", func(t *testing.T) {
		adjacents := wm.Adjacents("WalkBoxC")
		assert.Equal(t, 2, len(adjacents), "WalkBoxC should be adjacent to 2 WalkBoxes")
		assert.Contains(t, adjacents, wbB, "WalkBoxC should be adjacent to WalkBoxB")
		assert.Contains(t, adjacents, wbD, "WalkBoxC should be adjacent to WalkBoxD")
	})

	t.Run("Check WalkBoxD adjacency", func(t *testing.T) {
		adjacents := wm.Adjacents("WalkBoxD")
		assert.Equal(t, 1, len(adjacents), "WalkBoxD should be adjacent to 1 WalkBox")
		assert.Contains(t, adjacents, wbC, "WalkBoxD should be adjacent to WalkBoxC")
	})
}
