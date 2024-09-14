package pctk

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Music is music and we love it!
type Music = rl.Music

const (
	// StepsPerSecond is the  number of volume adjustment steps per second for a smooth sound transition
	StepsPerSecond = 60
	// MaxMasterVolume is the maximum value for the master volume (1.0 corresponds to 100%)
	MaxMasterVolume = 1.0
)

// LoadMusicFromFile - Load music stream from a file path
func LoadMusicFromFile(path string) *Music {
	music := rl.LoadMusicStream(path)
	if !rl.IsMusicReady(music) {
		log.Fatalf("Failed to load music from file %s", path)
	}
	return &music
}

// MusicPlay is a command that will play the music with the given resource locator.
type MusicPlay struct {
	MusicResource ResourceLocator
}

func (cmd MusicPlay) Execute(app *App, done Promise) {
	app.music = app.res.LoadMusic(cmd.MusicResource)
	if app.isMusicReady() {
		rl.PlayMusicStream(*app.music)
	}
	// TODO: determine if future is bounded to the music stream end or just the play begin.
	done.Complete()
}

// MusicStop is a command that will stop the music.
type MusicStop struct{}

func (cmd MusicStop) Execute(app *App, done Promise) {
	app.stopMusic()
	done.Complete()
}

// SetMasterVolume sets the global volume for the application.
func (a *App) SetMasterVolume(volume float32) {
	rl.SetMasterVolume(volume)
}

// GetMasterVolume returns the global master volume for the application.
func (a *App) GetMasterVolume() float32 {
	return rl.GetMasterVolume()
}
func (a *App) isMusicReady() bool {
	return a.music != nil && rl.IsMusicReady(*a.music)
}

func (a *App) updateMusic() {
	if a.isMusicReady() {
		rl.UpdateMusicStream(*a.music)
	}
}

func (a *App) pauseMusic() {
	if a.isMusicReady() {
		rl.PauseMusicStream(*a.music)
	}
}

func (a *App) resumeMusic() {
	if a.isMusicReady() {
		rl.ResumeMusicStream(*a.music)
	}
}

func (a *App) stopMusic() {
	if a.isMusicReady() {
		rl.StopMusicStream(*a.music)
	}
}

func (a *App) unloadMusic() {
	if a.isMusicReady() {
		rl.UnloadMusicStream(*a.music)
	}
}
