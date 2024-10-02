package pctk

// MusicPlay is a command that will play the music with the given resource reference.
type MusicPlay struct {
	MusicResource ResourceRef
}

func (cmd MusicPlay) Execute(app *App, done *Promise) {
	app.PlayMusic(cmd.MusicResource)
	// TODO: determine if future is bounded to the music stream end or just the play begin.
	done.Complete()
}

// MusicStop is a command that will stop the music.
type MusicStop struct{}

func (cmd MusicStop) Execute(app *App, done *Promise) {
	app.StopMusic()
	done.Complete()
}

// MusicPause is a command that will pause the music.
type MusicPause struct{}

func (cmd MusicPause) Execute(app *App, done *Promise) {
	app.PauseMusic()
	done.Complete()
}

// MusicResume is a command that will resume the music.
type MusicResume struct{}

func (cmd MusicResume) Execute(app *App, done *Promise) {
	app.ResumeMusic()
	done.Complete()
}
