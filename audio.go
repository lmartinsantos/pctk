package pctk

import (
	"io"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Music is music and we love it!
type Music struct {
	raw rl.Music
}

func (m *Music) BinaryEncode(w io.Writer) (int, error) {
	panic("not implemented")
}

// LoadMusicFromFile - Load music stream from a file path
func LoadMusicFromFile(path string) *Music {
	music := rl.LoadMusicStream(path)
	if !rl.IsMusicReady(rl.Music(music)) {
		log.Fatalf("Failed to load music from file %s", path)
	}
	return &Music{raw: music}
}

// MusicPlay is a command that will play the music with the given resource reference.
type MusicPlay struct {
	MusicResource ResourceRef
}

func (cmd MusicPlay) Execute(app *App, done Promise) {
	app.music = app.res.LoadMusic(cmd.MusicResource)
	if app.isMusicReady() {
		rl.PlayMusicStream(app.music.raw)
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

func (a *App) isMusicReady() bool {
	return a.music != nil && rl.IsMusicReady(a.music.raw)
}

func (a *App) updateMusic() {
	if a.isMusicReady() {
		rl.UpdateMusicStream(a.music.raw)
	}
}

func (a *App) pauseMusic() {
	if a.isMusicReady() {
		rl.PauseMusicStream(a.music.raw)
	}
}

func (a *App) resumeMusic() {
	if a.isMusicReady() {
		rl.ResumeMusicStream(a.music.raw)
	}
}

func (a *App) stopMusic() {
	if a.isMusicReady() {
		rl.StopMusicStream(a.music.raw)
	}
}

func (a *App) unloadMusic() {
	if a.isMusicReady() {
		rl.UnloadMusicStream(a.music.raw)
	}
}

// Sound source type
type Sound struct {
	raw rl.Sound
}

// LoadSoundFromFile - Load sound stream from a file path
func LoadSoundFromFile(path string) *Sound {
	sound := rl.LoadSound(path)
	if !rl.IsSoundReady(sound) {
		log.Fatalf("Failed to load sound from file %s", path)
	}
	return &Sound{raw: sound}
}

func (s *Sound) BinaryEncode(w io.Writer) (int, error) {
	panic("not implemented")
}

func (a *App) isSoundReady() bool {
	return a.sound != nil && rl.IsSoundReady(a.sound.raw)
}

// SoundPlay is a command that will play the sound with the given resource reference.
type SoundPlay struct {
	SoundResource ResourceRef
}

func (cmd SoundPlay) Execute(app *App, done Promise) {
	app.sound = app.res.LoadSound(cmd.SoundResource)
	if app.isSoundReady() {
		rl.PlaySound(app.sound.raw)
	}
	done.Complete()
}

// SoundStop is a command that will stop the sound with the given resource reference.
type SoundStop struct {
	SoundResource ResourceRef
}

func (cmd SoundStop) Execute(app *App, done Promise) {
	app.stopSound()
	done.Complete()
}

func (a *App) stopSound() {
	if a.isSoundReady() {
		rl.StopSound(a.sound.raw)
	}
}
