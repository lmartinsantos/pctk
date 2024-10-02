package pctk

// SoundPlay is a command that will play the sound with the given resource reference.
type SoundPlay struct {
	SoundResource ResourceRef
}

func (cmd SoundPlay) Execute(app *App, done *Promise) {
	app.PlaySound(cmd.SoundResource)
	done.Complete()
}

// SoundStop is a command that will stop the sound with the given resource reference.
type SoundStop struct {
	SoundResource ResourceRef
}

func (cmd SoundStop) Execute(app *App, done *Promise) {
	app.StopSound()
	done.Complete()
}
