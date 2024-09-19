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
	a.updateMusic()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.BeginMode2D(a.cam)
	a.drawBackgroud()
	a.drawControlPanel()
	a.drawDialogs()
	a.drawObjects()
	a.drawActors()
	a.drawMouseCursor()
	a.drawEgoAction()
	rl.EndMode2D()
	rl.EndDrawing()
	a.commands.execute(a)
}
