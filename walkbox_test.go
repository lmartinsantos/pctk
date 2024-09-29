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
