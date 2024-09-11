package pctk

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Scene represents a scene in the game.
type Scene struct {
	bg      *Image
	dialogs []Dialog
}

// NewScene creates a new scene with the given background image.
func NewScene(bg *Image) *Scene {
	if bg.Width() < ScreenWidth || bg.Height() < ScreenHeightScene {
		log.Fatal("Background image is too small")
	}
	return &Scene{
		bg: bg,
	}
}

// PlayScene loads and plays the scene with the given resource locator.
func (a *App) PlayScene(loc ResourceLocator) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// TODO: dispose the previous scene if any
	a.scene = a.cat.LoadScene(loc)
}

func (a *App) drawBackgroud() {
	if a.scene != nil {
		rl.DrawTexture(a.scene.bg.tex, 0, 0, rl.White)
	}
}
