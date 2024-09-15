package pctk

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Music is music and we love it!
type Music = rl.Music

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

// MusicPause is a command that will pause the music.
type MusicPause struct{}

func (cmd MusicPause) Execute(app *App, done Promise) {
	app.pauseMusic()
	done.Complete()
}

// MusicResume is a command that will resume the music.
type MusicResume struct{}

func (cmd MusicResume) Execute(app *App, done Promise) {
	app.resumeMusic()
	done.Complete()
}

// SetMusicVolume - Set volume for music (1.0 is max level)
func (a *App) SetMusicVolume(volume float32) {
	if a.isMusicReady() {
		rl.SetMusicVolume(*a.music, volume)
	}
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

// Sound source type
type Sound = rl.Sound

// LoadSoundFromFile - Load sound stream from a file path
func LoadSoundFromFile(path string) *Sound {
	sound := rl.LoadSound(path)
	if !rl.IsSoundReady(sound) {
		log.Fatalf("Failed to load sound from file %s", path)
	}
	return &sound
}

// SetSoundVolume - Set volume for sound (1.0 is max level)
func (a *App) SetSoundVolume(volume float32) {
	if a.isSoundReady() {
		rl.SetSoundVolume(*a.sound, volume)
	}
}

func (a *App) isSoundReady() bool {
	return a.sound != nil && rl.IsSoundReady(*a.sound)
}

// SoundPlay is a command that will play the sound with the given resource locator.
type SoundPlay struct {
	SoundResource ResourceLocator
}

func (cmd SoundPlay) Execute(app *App, done Promise) {
	app.sound = app.res.LoadSound(cmd.SoundResource)
	if app.isSoundReady() {
		rl.PlaySound(*app.sound)
	}
	done.Complete()
}

// SoundStop is a command that will stop the sound with the given resource locator.
type SoundStop struct {
	SoundResource ResourceLocator
}

func (cmd SoundStop) Execute(app *App, done Promise) {
	app.stopSound()
	done.Complete()
}

// SoundPause is a command that will pause the sound with the given resource locator.
type SoundPause struct {
	SoundResource ResourceLocator
}

func (cmd SoundPause) Execute(app *App, done Promise) {
	app.pauseSound()
	done.Complete()
}

// SoundResume is a command that will resume the sound with the given resource locator.
type SoundResume struct {
	SoundResource ResourceLocator
}

func (cmd SoundResume) Execute(app *App, done Promise) {
	app.resumeSound()
	done.Complete()
}

func (a *App) pauseSound() {
	if a.isSoundReady() {
		rl.PauseSound(*a.sound)
	}
}

func (a *App) resumeSound() {
	if a.isSoundReady() {
		rl.ResumeSound(*a.sound)
	}
}

func (a *App) stopSound() {
	if a.isSoundReady() {
		rl.StopSound(*a.sound)
	}
}
