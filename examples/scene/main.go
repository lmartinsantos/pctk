package main

import (
	"time"

	"github.com/apoloval/pctk"
)

func main() {
	bundle := pctk.NewResourceBundle()
	app := pctk.New(bundle)

	makeScene(bundle)
	app.Do(pctk.ScriptRun{ScriptResource: "/scripts/LostMyKeys"})
	app.Run()
}

func makeScene(bundle *pctk.ResourceBundle) {

	bg := pctk.LoadImageFromFile("background.jpg")
	scene := pctk.NewScene(bg)
	bundle.PutScene("/main", scene)

	sprites := pctk.LoadSpriteSheetFromFile("guybrush.png", pctk.Size{W: 32, H: 48})
	bundle.PutSpriteSheet("/guybrush/sprites", sprites)

	actor := pctk.NewActor("guybrush").
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

	bundle.PutMusic("/music/on-the-hill", pctk.LoadMusicFromFile("On_the_Hill.ogg"))
	bundle.PutMusic("/music/guitar_noodling", pctk.LoadMusicFromFile("guitar_noodling.ogg"))
	bundle.PutSound("/sound/cricket", pctk.LoadSoundFromFile("cricket.wav"))

	bundle.PutScript("/scripts/LostMyKeys", &pctk.Script{
		Language: pctk.ScriptLua,
		Code: []byte(`
			local pirate1_dialog_props = { pos = {x=60, y=20}, color = ColorMagenta }
			local pirate2_dialog_props = { pos = {x=60, y=50}, color = ColorYellow }

			ScenePlay("/main")
			ActorShow("/guybrush", "guybrush", {
				pos={x=340, y=90}, 
				dir=DirLeft
			})
			MusicPlay("/music/on-the-hill")
			SoundPlay("/sound/cricket")
			ActorWalkToPosition("guybrush", {x=290, y=90}).Wait()
			ActorSpeak("guybrush", "Hello, I'm Guybrush Threepwood,\nmighty pirate!").Wait()
			DialogShow("**Oh no! This guy again!**", pirate1_dialog_props)
			ActorWalkToPosition("guybrush", {x=120, y=90}).Wait()
			ActorSpeak("guybrush", "I think I've lost the keys to my boat.").Wait()
			ActorSpeak("guybrush", "Have you seen any keys?", {delay=2000}).Wait()
			DialogShow("Eeerrrr... Nope!", pirate1_dialog_props)
			SleepMillis(2000)

			MusicPlay("/music/guitar_noodling")
			ActorWalkToPosition("guybrush", {x=120, y=80}).Wait()
			ActorSpeak("guybrush", "Where can I find the keys?", {delay=1000}).Wait()
			ActorWalkToPosition("guybrush", {x=120, y=90}).Wait()
			ActorSpeak("guybrush", "Ooooook...").Wait()
			SleepMillis(2000)
			ActorStand("guybrush", {dir = DirRight}).Wait()
			SleepMillis(2000)
			ActorSpeak("guybrush", "Ok, I will try the Scumm bar.").Wait()
			ActorStand("guybrush", {dir = DirLeft}).Wait()
			ActorSpeak("guybrush", "Thank you guys!").Wait()
			SoundPlay("/sound/cricket")
			ActorWalkToPosition("guybrush", {x=360, y=90}).Wait()

			DialogShow("Oh, Jesus! I though he would\ntell again that stupid\ntale about LeChuck!", pirate1_dialog_props).Wait()
			SleepMillis(5000)
			DialogShow("Who has the keys?", pirate2_dialog_props).Wait()
			SleepMillis(1000)
			DialogShow("Me!", pirate1_dialog_props)
			ActorSelectEgo("guybrush").Wait()
			EnableControlPanel(true)
		`),
	})
}
