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
		<-app.ShowDialog("Hello, world!", 20, 20, rl.White, 1.0)
		<-app.ShowDialog("This is an example of a scene using raw functions.", 20, 20, rl.White, 1.0)
		<-app.ShowDialog("Do you remember the years of Monkey Island?", 20, 20, rl.Magenta, 1.0)
	}()
	app.Run()
}
