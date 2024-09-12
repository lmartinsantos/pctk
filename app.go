package pctk

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// App is the pctk application. It is the main struct that holds all the context necessary to run
// the application.
type App struct {
	mutex sync.Mutex

	res ResourceLoader

	screenCaption string
	screenZoom    int32

	scene   *Scene
	dialogs []Dialog
	actors  []*Actor

	cam               rl.Camera2D
	fontDefault       rl.Font
	fontDialogSolid   rl.Font
	fontDialogOutline rl.Font
	cursorTx          rl.Texture2D
	cursorColor       Color
	music             Music
}

// New creates a new pctk application.
func New(resources ResourceLoader, opts ...AppOption) *App {
	app := &App{
		res: resources,
	}

	opts = append(defaultAppOptions, opts...)
	for _, opt := range opts {
		opt(app)
	}

	app.init()

	return app
}

func (a *App) init() {
	rl.InitWindow(ScreenWidth*a.screenZoom, ScreenHeight*a.screenZoom, a.screenCaption)
	rl.InitAudioDevice()
	rl.SetTargetFPS(60)

	a.cam.Zoom = float32(a.screenZoom)
	a.initFonts()
	a.initMouse()
}

func (a *App) Close() {
	a.UnloadMusic()
	rl.CloseAudioDevice()
	rl.CloseWindow()
}
