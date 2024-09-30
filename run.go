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
	a.drawViewport()
	a.control.Draw(a)
	a.drawDialogs()
	a.drawMouseCursor()
	rl.EndMode2D()
	rl.EndDrawing()
	a.control.processControlInputs(a)
	a.commands.execute(a)
}
