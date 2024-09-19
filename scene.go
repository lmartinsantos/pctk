package pctk

import (
	"io"
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

func (s *Scene) BinaryEncode(w io.Writer) (int, error) {
	panic("not implemented")
}

// ScenePlay is a command that will play the scene with the given resource reference.
type ScenePlay struct {
	SceneResource ResourceRef
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
