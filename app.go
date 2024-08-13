package pctk

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const ()

// App is the pctk application. It is the main struct that holds all the context necessary to run
// the application.
type App struct {
	mutex sync.Mutex

	cat *ResourceCatalog

	screenCaption string
	screenZoom    int32

	background *Background
	dialogs    []Dialog

	cam               rl.Camera2D
	fontDefault       rl.Font
	fontDialogSolid   rl.Font
	fontDialogOutline rl.Font
	cursorTx          rl.Texture2D
	cursorColor       Color
}

// New creates a new pctk application.
func New(opts ...AppOption) *App {
	app := &App{
		cat: NewResourceCatalog(),
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
	rl.SetTargetFPS(60)

	a.cam.Zoom = float32(a.screenZoom)
	a.initFonts()
	a.initMouse()
}

// ResourceCatalog returns the resource catalog of the application.
func (a *App) ResourceCatalog() *ResourceCatalog {
	return a.cat
}
