package main

import (
	"github.com/apoloval/pctk"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	bundle := pctk.NewResourceBundle()
	app := pctk.New(bundle)

	makeScene(bundle)

	app.PlayScene("/main")
	go func() {
		<-app.ShowDialog("Don't sneak up on me like that!", 160, 20, rl.White, 1.0)
		<-app.ShowDialog("This is an example of a scene\nusing raw functions.", 160, 20, rl.White, 1.0)
		<-app.ShowDialog("Do you remember the years\nof Monkey Island?", 160, 20, rl.Magenta, 1.0)
	}()
	app.Run()
}

func makeScene(bundle *pctk.ResourceBundle) {
	bg := pctk.LoadImageFromFile("background.jpg")
	scene := pctk.NewScene(bg)
	bundle.PutScene("/main", scene)
}
