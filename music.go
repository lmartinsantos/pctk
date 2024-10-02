package pctk

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Music is music and we love it!
type Music struct {
	data   []byte
	format [4]byte
	raw    rl.Music
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
	copy(music.format[:], strings.ToUpper(filepath.Ext(path)))
	return music
}

// BinaryEncode encodes the music data to a binary stream. The format is:
//   - [4]byte: data format
//   - uint32: data length
//   - []byte: data
func (m *Music) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, m.format[:], uint32(len(m.data)), m.data)
}

// BinaryDecode decodes the music data from a binary stream. See Music.BinaryEncode for the format.
func (m *Music) BinaryDecode(r io.Reader) error {
	var format [4]byte
	var length uint32
	if err := BinaryDecode(r, &format, &length); err != nil {
		return err
	}

	data := make([]byte, length)
	if err := BinaryDecode(r, &data); err != nil {
		return err
	}

	m.format = format
	m.data = data
	m.raw = rl.LoadMusicStreamFromMemory(strings.ToLower(string(format[:])), data, int32(length))
	return nil
}

// PlayMusic plays the music from the given resource.
func (a *App) PlayMusic(ref ResourceRef) {
	a.music = a.res.LoadMusic(ref)
	if a.isMusicReady() {
		rl.PlayMusicStream(a.music.raw)
	}
}

// PauseMusic pauses the music being played.
func (a *App) PauseMusic() {
	if a.isMusicReady() {
		rl.PauseMusicStream(a.music.raw)
	}
}

// ResumeMusic resumes the music being played. It must be paused first.
func (a *App) ResumeMusic() {
	if a.isMusicReady() {
		rl.ResumeMusicStream(a.music.raw)
	}
}

// StopMusic stops the music being played.
func (a *App) StopMusic() {
	if a.isMusicReady() {
		rl.StopMusicStream(a.music.raw)
	}
}

func (a *App) updateMusic() {
	if a.isMusicReady() {
		rl.UpdateMusicStream(a.music.raw)
	}
}

func (a *App) unloadMusic() {
	if a.isMusicReady() {
		rl.UnloadMusicStream(a.music.raw)
	}
}
func (a *App) isMusicReady() bool {
	return a.music != nil && rl.IsMusicReady(a.music.raw)
}
