package main

import (
	"time"

	"github.com/apoloval/pctk"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	bundle := pctk.NewResourceBundle()
	app := pctk.New(bundle)

	makeScene(bundle)

	app.LoadMusic("On_the_Hill.ogg")
	app.PlayScene("/main")
	guybrush := app.ShowActor("/guybrush", pctk.NewPos(290, 90)).Looking(pctk.NewPos(0, 0)).Stand()
	go func() {
		<-app.ShowDialog("Don't sneak up on me like that!", 160, 20, rl.White, 1.0)
		guybrush.WalkTo(pctk.NewPos(120, 90))
		<-app.ShowDialog("This is an example of a scene\nusing raw functions.", 160, 20, rl.White, 1.0)
		guybrush.Looking(pctk.NewPos(290, 90)).Stand()
		<-app.ShowDialog("Do you remember the years\nof Monkey Island?", 160, 20, rl.Magenta, 1.0)
		guybrush.WalkTo(pctk.NewPos(360, 90))
	}()
	app.Run()
}

func makeScene(bundle *pctk.ResourceBundle) {
	bg := pctk.LoadImageFromFile("background.jpg")
	scene := pctk.NewScene(bg)
	bundle.PutScene("/main", scene)

	sprites := pctk.LoadSpriteSheetFromFile("guybrush.png", pctk.Size{W: 32, H: 48})
	bundle.PutSpriteSheet("/guybrush/sprites", sprites)

	actor := pctk.NewActor("Guybrush").
		WithStandH(pctk.NewAnimation("/guybrush/sprites").WithFrame(0, 1, time.Second)).
		WithWalkH(pctk.NewAnimation("/guybrush/sprites").
			WithFrame(0, 0, 100*time.Millisecond).
			WithFrame(1, 0, 100*time.Millisecond).
			WithFrame(2, 0, 100*time.Millisecond).
			WithFrame(1, 0, 100*time.Millisecond).
			WithFrame(0, 0, 100*time.Millisecond).
			WithFrame(3, 0, 100*time.Millisecond).
			WithFrame(4, 0, 100*time.Millisecond).
			WithFrame(5, 0, 100*time.Millisecond).
			WithFrame(4, 0, 100*time.Millisecond).
			WithFrame(3, 0, 100*time.Millisecond),
		)
	bundle.PutActor("/guybrush", actor)
}
