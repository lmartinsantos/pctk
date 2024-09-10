package main

import (
	"log"

	"github.com/apoloval/pctk"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var goodbyes = []string{
	"Venga, que curres.",
	"Deja de perder el tiempo.",
	"No seas perro.",
	"Luego no me eches la culpa\nsi tu manager te pilla.",
	"Por cierto, vendo estas estupendas\nchaquetas de cuero.",
}



func main() {
	app := pctk.New()
	bg := pctk.BackgroundFromImage(rl.LoadImage("background.jpg"))
	app.ResourceCatalog().Add("/background", bg)
	if err := app.SetBackground("/background"); err != nil {
		log.Fatalf("Failed setting background: %v", err)
	}
	go func() {
		<-app.ShowDialog("Hi, my name is Javier Solana\nand I want to be a pirate!", 120, 20, rl.White, 1.0)
		<-app.ShowDialog("Si, correcto. Es un proto-invento\npara hacer aventuras graficas.", 120, 20, rl.Magenta, 1.0)
		<-app.ShowDialog("Aun faltan bastantes cosas.\nPero nos podria servir en Kenia.", 120, 20, rl.Red, 1.0)
		<-app.ShowDialog("Ahora deja de perder el tiempo\ny ponte a currar.", 120, 20, rl.Yellow, 1.0)
		<-app.ShowDialog("Que el Dev/X no sale adelante solo!", 120, 20, rl.Yellow, 1.0)

		for {
			for _, g := range goodbyes {
				<-app.ShowDialog(g, 120, 20, rl.White, 1.0)
			}
		}
	}()
	app.Run()
}
