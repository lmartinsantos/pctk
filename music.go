package pctk

import (
	"log"
	"os"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Music is music and we love it!
type Music = rl.Music

func (a *App) LoadMusic(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		log.Fatalf("An error occurred while attempting to load the music file %s: %v", fileName, err)
	}

	// dispose the previous music if any
	if unsafe.Pointer(&a.music) != nil {
		a.StopMusic()
		a.UnloadMusic()
	}

	a.music = rl.LoadMusicStream(fileName)
	rl.PlayMusicStream(a.music)
}

func (a *App) SetMasterVolume(volume float32) {
	rl.SetMasterVolume(volume)
}

func (a *App) GetMasterVolume() float32 {
	return rl.GetMasterVolume()
}

func (a *App) UpdateMusic() {
	if unsafe.Pointer(&a.music) != nil {
		rl.UpdateMusicStream(a.music)
	}
}

func (a *App) StopMusic() {
	if unsafe.Pointer(&a.music) != nil {
		rl.StopMusicStream(a.music)
	}
}

func (a *App) UnloadMusic() {
	if unsafe.Pointer(&a.music) != nil {
		rl.UnloadMusicStream(a.music)
	}
}
