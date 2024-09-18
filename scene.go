package pctk

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var SceneViewport = NewRect(0, 0, ScreenWidth, ScreenHeightScene)

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

// ScenePlay is a command that will play the scene with the given resource locator.
type ScenePlay struct {
	SceneResource ResourceLocator
}

func (cmd ScenePlay) Execute(app *App, done Promise) {
	// TODO: dispose the previous scene if any
	app.scene = app.res.LoadScene(cmd.SceneResource)
	done.Complete()
}

func (a *App) drawBackgroud() {
	if a.scene != nil {
		rl.DrawTexture(a.scene.bg.tex, 0, 0, rl.White)
	}
}
