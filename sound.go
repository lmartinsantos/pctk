package pctk

import (
	"io"
	"log"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Sound source type
type Sound struct {
	data   []byte
	raw    rl.Sound
	format [4]byte
}

// LoadSoundFromFile - Load sound stream from a file path
func LoadSoundFromFile(path string) *Sound {
	var err error
	sound := new(Sound)

	sound.data, err = os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read sound file: %v", err)
	}
	wav := rl.LoadWaveFromMemory(filepath.Ext(path), sound.data, int32(len(sound.data)))
	sound.raw = rl.LoadSoundFromWave(wav)
	copy(sound.format[:], filepath.Ext(path))
	return sound
}

// BinaryEncode encodes the sound data to a binary stream. The format is:
//   - [4]byte: data format
//   - uint32: data length
//   - []byte: data
func (s *Sound) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, s.format[:], uint32(len(s.data)), s.data)
}

// BinaryDecode decodes the sound data from a binary stream. See Sound.BinaryEncode for the format.
func (s *Sound) BinaryDecode(r io.Reader) error {
	var format [4]byte
	var length uint32
	if err := BinaryDecode(r, &format, &length); err != nil {
		return err
	}

	data := make([]byte, length)
	if err := BinaryDecode(r, &data); err != nil {
		return err
	}

	s.format = format
	s.data = data
	wav := rl.LoadWaveFromMemory(string(format[:]), data, int32(length))
	s.raw = rl.LoadSoundFromWave(wav)
	return nil
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
