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
	app.Do(pctk.ActorShow{
		ActorResource: "/pirate1",
		ActorName:     "pirate1",
		Position:      pctk.NewPos(30, 90),
		LookAt:        pctk.DirRight,
	})
	app.Do(pctk.ActorShow{
		ActorResource: "/pirate2",
		ActorName:     "pirate2",
		Position:      pctk.NewPos(45, 75),
		LookAt:        pctk.DirRight,
	})
	app.Do(pctk.ActorShow{
		ActorResource: "/pirate3",
		ActorName:     "pirate3",
		Position:      pctk.NewPos(5, 100), // TODO  playing with position can reproduce blinking issue (no z-index solution yet)
		LookAt:        pctk.DirRight,
	})
	app.Do(pctk.MusicPlay{MusicResource: "/music/on-the-hill"})
	app.Do(pctk.SoundPlay{SoundResource: "/sound/cricket"})
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
		app.Do(pctk.ActorSpeak{
			ActorName: "pirate1",
			Text:      "Eeerrrr... Nope!",
			Delay:     2 * time.Second,
			Color:     pctk.BrigthYellow,
		}).Wait()
		app.Do(pctk.MusicPlay{MusicResource: "/music/guitar_noodling"})
		app.Do(pctk.ActorWalkToPosition{
			ActorName: "guybrush",
			Position:  pctk.NewPos(120, 80),
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "Where can I find the keys?",
			Delay:     1 * time.Second,
		}).Wait()
		app.Do(pctk.ActorWalkToPosition{
			ActorName: "guybrush",
			Position:  pctk.NewPos(120, 90),
		}).Wait()
		pctk.WithDelay(
			app.Do(pctk.ActorSpeak{
				ActorName: "guybrush",
				Text:      "Ooooook...",
			}),
			2*time.Second,
		).Wait()
		pctk.WithDelay(
			app.Do(pctk.ActorStand{
				ActorName: "guybrush",
				Direction: pctk.DirRight,
			}),
			2*time.Second,
		).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "Ok, I will try the Scumm bar.",
		}).Wait()
		app.Do(pctk.ActorStand{
			ActorName: "guybrush",
			Direction: pctk.DirLeft,
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "guybrush",
			Text:      "Thank you guys!",
		}).Wait()
		app.Do(pctk.SoundPlay{SoundResource: "/sound/cricket"})
		app.Do(pctk.ActorWalkToPosition{
			ActorName: "guybrush",
			Position:  pctk.NewPos(360, 90),
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "pirate2",
			Text:      "Oh, Jesus! I though we would\ntell again that stupid\ntale about LeChuck!",
			Delay:     5 * time.Second,
			Color:     pctk.Magenta,
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "pirate1",
			Text:      "Who has the keys?",
			Delay:     2 * time.Second,
			Color:     pctk.BrigthYellow,
		}).Wait()
		app.Do(pctk.ActorSpeak{
			ActorName: "pirate2",
			Text:      "Me!",
			Delay:     1 * time.Second,
			Color:     pctk.Magenta,
		}).Wait()
	}()
	app.Run()
}

func makeScene(bundle *pctk.ResourceBundle) {

	bg := pctk.LoadImageFromFile("background-2.png")
	scene := pctk.NewScene(bg)
	bundle.PutScene("/main", scene)

	sprites := pctk.LoadSpriteSheetFromFile("guybrush.png", pctk.Size{W: 32, H: 48})
	bundle.PutSpriteSheet("/guybrush/sprites", sprites)

	actor := pctk.NewActor("Guybrush").
		WithAnimationStand(pctk.DirRight, pctk.NewAnimation("/guybrush/sprites").
			WithFrame(0, 1, time.Second),
		).
		WithAnimationStand(pctk.DirLeft, pctk.NewAnimation("/guybrush/sprites").
			WithFrame(0, 1, time.Second).
			Flip(true),
		).
		WithAnimationStand(pctk.DirUp, pctk.NewAnimation("/guybrush/sprites").
			WithFrame(0, 5, time.Second),
		).
		WithAnimationStand(pctk.DirDown, pctk.NewAnimation("/guybrush/sprites").
			WithFrame(0, 4, time.Second),
		).
		WithAnimationSpeak(pctk.DirRight, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(1, 100*time.Millisecond, 0, 1, 2, 3, 4, 5),
		).
		WithAnimationSpeak(pctk.DirLeft, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(1, 100*time.Millisecond, 0, 1, 2, 3, 4, 5).Flip(true),
		).
		WithAnimationSpeak(pctk.DirUp, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(5, 100*time.Millisecond, 0, 1, 2),
		).
		WithAnimationSpeak(pctk.DirDown, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(4, 100*time.Millisecond, 0, 1, 2, 3, 4, 5),
		).
		WithAnimationWalk(pctk.DirRight, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(0, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3),
		).
		WithAnimationWalk(pctk.DirLeft, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(0, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3).
			Flip(true),
		).
		WithAnimationWalk(pctk.DirUp, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(3, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3),
		).
		WithAnimationWalk(pctk.DirDown, pctk.NewAnimation("/guybrush/sprites").
			WithFramesInRow(2, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3),
		)
	bundle.PutActor("/guybrush", actor)

	sprites = pctk.LoadSpriteSheetFromFile("pirate-001.png", pctk.Size{W: 26, H: 40})
	bundle.PutSpriteSheet("/pirate1/sprites", sprites)

	actor = pctk.NewActor("Pirate1").
		WithAnimationStand(pctk.DirRight, pctk.NewAnimation("/pirate1/sprites").
			WithFrame(0, 0, time.Second),
		).
		WithAnimationSpeak(pctk.DirRight, pctk.NewAnimation("/pirate1/sprites").
			WithFramesInRow(0, 100*time.Millisecond, 0, 1, 2, 3),
		)
	bundle.PutActor("/pirate1", actor)

	sprites = pctk.LoadSpriteSheetFromFile("pirate-002.png", pctk.Size{W: 25, H: 50})
	bundle.PutSpriteSheet("/pirate2/sprites", sprites)

	actor = pctk.NewActor("Pirate2").
		WithAnimationStand(pctk.DirRight, pctk.NewAnimation("/pirate2/sprites").
			WithFrame(0, 0, time.Second),
		).
		WithAnimationSpeak(pctk.DirRight, pctk.NewAnimation("/pirate2/sprites").
			WithFramesInRow(0, 100*time.Millisecond, 0, 1, 2, 3),
		)
	bundle.PutActor("/pirate2", actor)

	sprites = pctk.LoadSpriteSheetFromFile("pirate-003.png", pctk.Size{W: 35, H: 40})
	bundle.PutSpriteSheet("/pirate3/sprites", sprites)

	actor = pctk.NewActor("Pirate3").
		WithAnimationStand(pctk.DirRight, pctk.NewAnimation("/pirate3/sprites").
			WithFramesInRow(0, 5*time.Second, 0, 1),
		)
	bundle.PutActor("/pirate3", actor)

	bundle.PutMusic("/music/on-the-hill", pctk.LoadMusicFromFile("On_the_Hill.ogg"))
	bundle.PutMusic("/music/guitar_noodling", pctk.LoadMusicFromFile("guitar_noodling.ogg"))
	bundle.PutSound("/sound/cricket", pctk.LoadSoundFromFile("cricket.wav"))
}
