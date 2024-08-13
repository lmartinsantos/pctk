package main

import (
	"log"

	"github.com/apoloval/pctk"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	app := pctk.New()
	bg := pctk.BackgroundFromImage(rl.LoadImage("background.jpg"))
	app.ResourceCatalog().Add("/background", bg)
	if err := app.SetBackground("/background"); err != nil {
		log.Fatalf("Failed setting background: %v", err)
	}
	go func() {
		<-app.ShowDialog("Don't sneak up on me like that!", 160, 20, rl.White, 1.0)
		<-app.ShowDialog("This is an example of a scene\nusing raw functions.", 160, 20, rl.White, 1.0)
		<-app.ShowDialog("Do you remember the years\nof Monkey Island?", 160, 20, rl.Magenta, 1.0)
	}()
	app.Run()
}
