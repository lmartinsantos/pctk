package pctk

import (
	"io"
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
// row and the indicated cols of the sprite sheet. The delay is the time to wait before moving to
// the next frame.
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

// BinaryEncode encodes the animation to a binary format. The format is as follows:
// - byte: the flip flag.
// - uint32: the number of frames.
// - for each frame:
//   - byte: the i-th row.
//   - byte: the j-th column.
//   - uint64: the delay.
func (a *Animation) BinaryEncode(w io.Writer) (n int, err error) {
	n, err = BinaryEncode(w, a.flip, uint32(len(a.frames)))
	for _, frame := range a.frames {
		nn, err := BinaryEncode(w, byte(frame.i), byte(frame.j), uint64(frame.delay))
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
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
