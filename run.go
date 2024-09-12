package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (a *App) Run() {
	defer a.Close()

	for !rl.WindowShouldClose() {
		a.run()
	}
}

func (a *App) run() {
	a.UpdateMusic()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(a.cam)
	a.drawBackgroud()
	a.drawControlPanel()
	a.drawDialogs()
	a.drawActors()
	a.drawMouseCursor()
	rl.EndMode2D()
	rl.EndDrawing()
}
