package pctk

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Scene represents a scene in the game.
type Scene struct {
	bg            *Image
	dialogs       []Dialog
	sceneViewPort Rectangle
}

// NewScene creates a new scene with the given background image.
func NewScene(bg *Image) *Scene {
	if bg.Width() < ScreenWidth || bg.Height() < ScreenHeightScene {
		log.Fatal("Background image is too small")
	}
	return &Scene{
		bg:            bg,
		sceneViewPort: NewRect(0, 0, int(bg.Width()), int(bg.Height())),
	}
}

// PlayScene is a command that will play the scene with the given resource locator.
type PlayScene struct {
	SceneResource ResourceLocator
}

func (cmd PlayScene) Execute(app *App, done Promise) {
	// TODO: dispose the previous scene if any
	app.scene = app.res.LoadScene(cmd.SceneResource)
	done.Complete()
}

func (a *App) drawBackgroud() {
	if a.scene != nil {
		rl.DrawTexture(a.scene.bg.tex, 0, 0, rl.White)
	}
}
