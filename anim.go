package pctk

import (
	"time"
)

const (
	DefaultAnimationDelay = 100 * time.Millisecond
)

// Animation represents a sequence of images that can be played.
type Animation struct {
	frames []animationFrame
	flip   bool

	currentFrame int
	lastFrame    time.Time
}

// NewAnimation creates a new animation.
func NewAnimation() *Animation {
	return &Animation{}
}

// WithFrame adds a frame to the animation. The frame is located at the i-th row and j-th column of
// the sprite sheet. The delay is the time to wait before moving to the next frame.
func (a *Animation) WithFrame(i, j uint, delay time.Duration) *Animation {
	a.frames = append(a.frames, animationFrame{i, j, delay})
	return a
}

// WithFramesInRow adds a sequence of frames to the animation. The frames are located at the row-th
// row and fromCol-th to toCol-th columns of the sprite sheet. The delay is the time to wait before
// moving to the next frame.
func (a *Animation) WithFramesInRow(row uint, delay time.Duration, cols ...uint) *Animation {
	for _, col := range cols {
		a.WithFrame(col, row, delay)
	}
	return a
}

// Flip sets the flip flag for the animation.
func (a *Animation) Flip(flip bool) *Animation {
	a.flip = flip
	return a
}

func (a *Animation) draw(sprites *SpriteSheet, pos Position) {
	if a.frames[a.currentFrame].delay < time.Since(a.lastFrame) {
		a.lastFrame = time.Now()
		a.currentFrame++
		if a.currentFrame >= len(a.frames) {
			a.currentFrame = 0
		}
	}

	sprites.DrawSprite(
		a.frames[a.currentFrame].i,
		a.frames[a.currentFrame].j,
		pos,
		a.flip,
	)
}

type animationFrame struct {
	i, j  uint
	delay time.Duration
}
