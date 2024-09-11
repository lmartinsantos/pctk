package pctk

import (
	"time"
)

const (
	DefaultAnimationDelay = 100 * time.Millisecond
)

// Animation represents a sequence of images that can be played.
type Animation struct {
	sprites ResourceLocator
	frames  []animationFrame

	currentFrame int
	lastFrame    time.Time
}

// NewAnimation creates a new animation.
func NewAnimation(sprites ResourceLocator) *Animation {
	return &Animation{
		sprites: sprites,
	}
}

// WithFrame adds a frame to the animation. The frame is located at the i-th row and j-th column of
// the sprite sheet. The delay is the time to wait before moving to the next frame.
func (a *Animation) WithFrame(i, j uint, delay time.Duration) *Animation {
	a.frames = append(a.frames, animationFrame{i, j, delay})
	return a
}

func (a *Animation) draw(app *App, pos Position, flip bool) {
	if a.frames[a.currentFrame].delay < time.Since(a.lastFrame) {
		a.lastFrame = time.Now()
		a.currentFrame++
		if a.currentFrame >= len(a.frames) {
			a.currentFrame = 0
		}
	}

	sprites := app.res.LoadSpriteSheet(a.sprites)
	sprites.DrawSprite(
		a.frames[a.currentFrame].i,
		a.frames[a.currentFrame].j,
		pos,
		flip,
	)
}

type animationFrame struct {
	i, j  uint
	delay time.Duration
}
