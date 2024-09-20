package main

import (
	"time"

	"github.com/apoloval/pctk"
)

func main() {
	bundle := pctk.NewResourceBundle()
	app := pctk.New(bundle)

	buildResources(bundle)
	app.Do(pctk.ScriptRun{ScriptResource: pctk.NewResourceRef("resources", "scripts/LostMyKeys")})
	app.Run()
}

func buildResources(bundle *pctk.ResourceBundle) {

	bg := pctk.LoadImageFromFile("background.jpg")
	room := pctk.NewRoom(bg)
	bundle.PutRoom(pctk.NewResourceRef("resources", "rooms/Melee"), room)

	sprites := pctk.LoadSpriteSheetFromFile("guybrush.png", pctk.Size{W: 32, H: 48})
	costume := pctk.NewCostume(sprites).
		WithAnimation(pctk.CostumeIdle(pctk.DirRight), pctk.NewAnimation().
			WithFrame(0, 1, time.Second),
		).
		WithAnimation(pctk.CostumeIdle(pctk.DirLeft), pctk.NewAnimation().
			WithFrame(0, 1, time.Second).
			Flip(true),
		).
		WithAnimation(pctk.CostumeIdle(pctk.DirUp), pctk.NewAnimation().
			WithFrame(0, 5, time.Second),
		).
		WithAnimation(pctk.CostumeIdle(pctk.DirDown), pctk.NewAnimation().
			WithFrame(0, 4, time.Second),
		).
		WithAnimation(pctk.CostumeSpeak(pctk.DirRight), pctk.NewAnimation().
			WithFramesInRow(1, 100*time.Millisecond, 0, 1, 2, 3, 4, 5),
		).
		WithAnimation(pctk.CostumeSpeak(pctk.DirLeft), pctk.NewAnimation().
			WithFramesInRow(1, 100*time.Millisecond, 0, 1, 2, 3, 4, 5).Flip(true),
		).
		WithAnimation(pctk.CostumeSpeak(pctk.DirUp), pctk.NewAnimation().
			WithFramesInRow(5, 100*time.Millisecond, 0, 1, 2),
		).
		WithAnimation(pctk.CostumeSpeak(pctk.DirDown), pctk.NewAnimation().
			WithFramesInRow(4, 100*time.Millisecond, 0, 1, 2, 3, 4, 5),
		).
		WithAnimation(pctk.CostumeWalk(pctk.DirRight), pctk.NewAnimation().
			WithFramesInRow(0, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3),
		).
		WithAnimation(pctk.CostumeWalk(pctk.DirLeft), pctk.NewAnimation().
			WithFramesInRow(0, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3).
			Flip(true),
		).
		WithAnimation(pctk.CostumeWalk(pctk.DirUp), pctk.NewAnimation().
			WithFramesInRow(3, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3),
		).
		WithAnimation(pctk.CostumeWalk(pctk.DirDown), pctk.NewAnimation().
			WithFramesInRow(2, 100*time.Millisecond, 0, 1, 2, 1, 0, 3, 4, 5, 4, 3),
		)
	bundle.PutCostume(pctk.NewResourceRef("resources", "costumes/Guybrush"), costume)

	bundle.PutMusic(
		pctk.NewResourceRef("resources", "audio/OnTheHill"),
		pctk.LoadMusicFromFile("On_the_Hill.ogg"),
	)
	bundle.PutMusic(
		pctk.NewResourceRef("resources", "audio/GuitarNoodling"),
		pctk.LoadMusicFromFile("guitar_noodling.ogg"),
	)
	bundle.PutSound(
		pctk.NewResourceRef("resources", "audio/Cricket"),
		pctk.LoadSoundFromFile("cricket.wav"),
	)

	bundle.PutScript(pctk.NewResourceRef("resources", "scripts/LostMyKeys"), &pctk.Script{
		Language: pctk.ScriptLua,
		Code: []byte(`
			local pirate1_dialog_props = { pos = {x=60, y=20}, color = ColorMagenta }
			local pirate2_dialog_props = { pos = {x=60, y=50}, color = ColorYellow }

			RoomShow("resources:rooms/Melee")
			ActorShow("guybrush", {
				pos={x=340, y=90}, 
				dir=DirLeft,
				costume="resources:costumes/Guybrush"
			})
			MusicPlay("resources:audio/OnTheHill")
			SoundPlay("resources:audio/Cricket")
			ActorWalkToPosition("guybrush", {x=290, y=90}).Wait()
			ActorSpeak("guybrush", "Hello, I'm Guybrush Threepwood,\nmighty pirate!").Wait()
			DialogShow("**Oh no! This guy again!**", pirate1_dialog_props)
			ActorWalkToPosition("guybrush", {x=120, y=90}).Wait()
			ActorSpeak("guybrush", "I think I've lost the keys to my boat.").Wait()
			ActorSpeak("guybrush", "Have you seen any keys?", {delay=2000}).Wait()
			DialogShow("Eeerrrr... Nope!", pirate1_dialog_props)
			SleepMillis(2000)

			MusicPlay("resources:audio/GuitarNoodling")
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
			SoundPlay("resources:audio/Cricket")
			ActorWalkToPosition("guybrush", {x=360, y=90}).Wait()

			DialogShow("Oh, Jesus! I though he would\ntell again that stupid\ntale about LeChuck!", pirate1_dialog_props).Wait()
			SleepMillis(5000)
			DialogShow("Who has the keys?", pirate2_dialog_props).Wait()
			SleepMillis(1000)
			DialogShow("Me!", pirate1_dialog_props)
		`),
	})
}
