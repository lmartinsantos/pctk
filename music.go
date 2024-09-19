package pctk

import (
	"io"
	"log"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Music is music and we love it!
type Music struct {
	data   []byte
	format [4]byte
	raw    rl.Music
}

// BinaryEncode encodes the music data to a binary stream. The format is:
//   - [4]byte: data format
//   - uint32: data length
//   - []byte: data
func (m *Music) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, m.format[:], uint32(len(m.data)), m.data)
}

// LoadMusicFromFile - Load music stream from a file path
func LoadMusicFromFile(path string) *Music {
	var err error
	music := new(Music)
	music.data, err = os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read music file: %v", err)
	}

	music.raw = rl.LoadMusicStreamFromMemory(filepath.Ext(path), music.data, int32(len(music.data)))
	copy(music.format[:], filepath.Ext(path))
	return music
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
