package main

import (
	"time"

	"github.com/apoloval/pctk"
)

func main() {
	bundle := pctk.NewResourceBundle()
	app := pctk.New(bundle)

	makeScene(bundle)
	app.Do(pctk.PlayScene{SceneResource: "/main"})
	app.Do(pctk.ActorShow{
		ActorResource: "/guybrush",
		ActorName:     "guybrush",
		Position:      pctk.NewPos(340, 90),
		LookAt:        pctk.DirLeft,
	})
	app.Do(pctk.MusicPlay{MusicResource: "/music/on-the-hill"})
	go func() {
		app.Do(pctk.ActorWalkToPosition{
			ActorName: "guybrush",
			Position:  pctk.NewPos(290, 90),
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "Hello, I'm Guybrush Threepwood,\nmighty pirate!",
		}).Wait()
		app.Do(pctk.ShowDialog{
			Text:     "**Oh no! This guy again!**",
			Position: pctk.NewPos(60, 20),
			Color:    pctk.Magenta,
		})
		app.Do(pctk.ActorWalkToPosition{
			ActorName: "guybrush",
			Position:  pctk.NewPos(120, 90),
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "I think I've lost my boat keys.",
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "Have you seen any keys?",
			Delay:     2 * time.Second,
		}).Wait()
		pctk.WithDelay(
			app.Do(pctk.ShowDialog{
				Text:     "Eeerrrr... Nope!",
				Position: pctk.NewPos(60, 20),
				Color:    pctk.Magenta,
			}),
			2*time.Second,
		).Wait()
		app.Do(pctk.MusicPlay{MusicResource: "/music/guitar_noodling"})
		pctk.WithDelay(
			app.Do(pctk.ActorLookAtDirection{
				ActorName: "guybrush",
				Direction: pctk.DirRight,
			}),
			2*time.Second,
		).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "Ok, I will try the Scumm bar.",
		}).Wait()
		app.Do(pctk.ActorLookAtDirection{
			ActorName: "guybrush",
			Direction: pctk.DirLeft,
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "Thank you guys!",
		}).Wait()
		app.Do(pctk.ActorWalkToPosition{
			ActorName: "guybrush",
			Position:  pctk.NewPos(360, 90),
		}).Wait()
		pctk.WithDelay(
			app.Do(pctk.ShowDialog{
				Text:     "Oh, Jesus! I though we would\ntell again that stupid\ntale about LeChuck!",
				Position: pctk.NewPos(60, 20),
				Color:    pctk.Magenta,
			}),
			5*time.Second,
		).Wait()
		pctk.WithDelay(
			app.Do(pctk.ShowDialog{
				Text:     "Who has the keys?",
				Position: pctk.NewPos(60, 20),
				Color:    pctk.BrigthYellow,
			}),
			1*time.Second,
		).Wait()
		app.Do(pctk.ShowDialog{
			Text:     "Me!",
			Position: pctk.NewPos(60, 20),
			Color:    pctk.Magenta,
		}).Wait()
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
		).
		WithSpeakH(pctk.NewAnimation("/guybrush/sprites").
			WithFrame(0, 1, 100*time.Millisecond).
			WithFrame(1, 1, 100*time.Millisecond).
			WithFrame(2, 1, 100*time.Millisecond).
			WithFrame(3, 1, 100*time.Millisecond).
			WithFrame(4, 1, 100*time.Millisecond).
			WithFrame(5, 1, 100*time.Millisecond),
		)
	bundle.PutActor("/guybrush", actor)

	bundle.PutMusic("/music/on-the-hill", pctk.LoadMusicFromFile("On_the_Hill.ogg"))
	bundle.PutMusic("/music/guitar_noodling", pctk.LoadMusicFromFile("guitar_noodling.ogg"))
}
