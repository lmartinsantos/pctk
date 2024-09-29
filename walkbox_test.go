package pctk_test

import (
	"fmt"
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
		vertices    []*pctk.Position
		shouldPanic bool
		message     string
	}{
		{
			name:        "Insufficient vertices for a polygon",
			vertices:    []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}},
			shouldPanic: true,
			message:     fmt.Sprintf("Expected panic because we don't have enough vertices for a polygon (min required is %d)!", pctk.MinPolygonVertices),
		},
		{
			name:        "Concave polygon should panic",
			vertices:    []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 2, Y: 1}, {X: 4, Y: 4}},
			shouldPanic: true,
			message:     "Expected panic because vertices form a concave polygon!",
		},
		{
			name:        "Collinear vertices should panic",
			vertices:    []*pctk.Position{{X: 1, Y: 2}, {X: 2, Y: 4}, {X: 3, Y: 6}, {X: 4, Y: 8}},
			shouldPanic: true,
			message:     "Expected panic because vertices are collinear!",
		},
		{
			name:        "Should successfully create a valid WalkBox with a convex polygon",
			vertices:    []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			shouldPanic: false,
			message:     "Expected create a valid WalkBox, vertices form a convex polygon!",
		},
		{
			name:        "Should successfully create a valid WalkBox with a complex convex polygon (octagon)",
			vertices:    []*pctk.Position{{X: 0, Y: 0}, {X: 1, Y: 2}, {X: 2, Y: 3}, {X: 3, Y: 2}, {X: 4, Y: 0}, {X: 3, Y: -2}, {X: 2, Y: -3}, {X: 1, Y: -2}},
			shouldPanic: false,
			message:     "Expected create a valid WalkBox, vertices form a complex convex octagon!",
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
		vertices   []*pctk.Position
		point      *pctk.Position
		assertFunc func(t *testing.T, isInside bool)
	}{
		{
			name:     "The point should be considered inside the polygon when it is on the edge",
			vertices: []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Position{X: 2, Y: 0}, // On the edge
			assertFunc: func(t *testing.T, isInside bool) {
				assert.True(t, isInside)
			},
		},
		{
			name:     "The point should be inside the polygon",
			vertices: []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Position{X: 2, Y: 2},
			assertFunc: func(t *testing.T, isInside bool) {
				assert.True(t, isInside)
			},
		},
		{
			name:     "The point should be outside the polygon",
			vertices: []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Position{X: 5, Y: 5},
			assertFunc: func(t *testing.T, isInside bool) {
				assert.False(t, isInside)
			},
		},
		{
			name:     "The point should be considered inside the polygon when it is on a vertex",
			vertices: []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Position{X: 0, Y: 0}, // On the vertex
			assertFunc: func(t *testing.T, isInside bool) {
				assert.True(t, isInside)
			},
		},
		{
			name:     "The point should be inside a complex convex polygon with more vertices",
			vertices: []*pctk.Position{{X: 0, Y: 0}, {X: 5, Y: 0}, {X: 6, Y: 3}, {X: 4, Y: 6}, {X: 1, Y: 5}, {X: 0, Y: 2}},
			point:    &pctk.Position{X: 3, Y: 3},
			assertFunc: func(t *testing.T, isInside bool) {
				assert.True(t, isInside)
			},
		},
		{
			name:     "The point should be outside when it is far from the polygon",
			vertices: []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			point:    &pctk.Position{X: 10, Y: 10}, // Clearly outside the polygon
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

func TestEnable(t *testing.T) {
	testCases := []struct {
		name       string
		enable     bool
		assertFunc func(t *testing.T, walkBox *pctk.WalkBox)
	}{
		{
			name:   "WalkBox.Enabled should set it to true",
			enable: true,
			assertFunc: func(t *testing.T, walkBox *pctk.WalkBox) {
				assert.True(t, walkBox.Enabled, "Expected WalkBox to be enabled")
			},
		},
		{
			name:   "WalkBox.Enabled should set it to false (disabled)",
			enable: false,
			assertFunc: func(t *testing.T, walkBox *pctk.WalkBox) {
				assert.False(t, walkBox.Enabled, "Expected WalkBox to be disabled")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			vertices := []*pctk.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}}
			walkBox := pctk.NewWalkBox(DefaultWalkBoxID, vertices)
			walkBox.Enable(testCase.enable)
			testCase.assertFunc(t, walkBox)
		})
	}
}
