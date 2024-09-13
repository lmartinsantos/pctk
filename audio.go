package pctk

import (
	"log"
	"time"

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

// IsMusicPlaying - Check if music is playing
func (a *App) IsMusicPlaying() bool {
	return a.IsMusicReady() && rl.IsMusicStreamPlaying(*a.music)
}

// IsMusicReady - Check if the music stream is ready to play
func (a *App) IsMusicReady() bool {
	return a.music != nil && rl.IsMusicReady(*a.music)
}

// PlayMusic - Load and start playing music from a given resource location
func (a *App) PlayMusic(loc ResourceLocator) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.music = a.res.LoadMusic(loc)
	if a.IsMusicReady() {
		rl.PlayMusicStream(*a.music)
	}
}

// UpdateMusic - Update the currently playing music stream
func (a *App) UpdateMusic() {
	if a.IsMusicReady() {
		rl.UpdateMusicStream(*a.music)
	}
}

// PauseMusic - Pause the currently playing music stream
func (a *App) PauseMusic() {
	if a.IsMusicReady() {
		rl.PauseMusicStream(*a.music)
	}
}

// ResumeMusic - Resume the paused music stream
func (a *App) ResumeMusic() {
	if a.IsMusicReady() {
		rl.ResumeMusicStream(*a.music)
	}
}

// StopMusic - Stop the currently playing music stream
func (a *App) StopMusic() {
	if a.IsMusicReady() {
		rl.StopMusicStream(*a.music)
	}
}

// UnloadMusic - Unload and free memory associated with the music stream
func (a *App) UnloadMusic() {
	if a.IsMusicReady() {
		rl.UnloadMusicStream(*a.music)
	}
}

// SetMasterVolume - Set the global master volume for the application
func (a *App) SetMasterVolume(volume float32) {
	rl.SetMasterVolume(volume)
}

// GetMasterVolume - Get the current global master volume
func (a *App) GetMasterVolume() float32 {
	return rl.GetMasterVolume()
}

// SetMusicPan - Set the stereo panning for the currently playing music
func (a *App) SetMusicPan(pan float32) {
	if a.IsMusicReady() {
		rl.SetMusicPan(*a.music, pan)
	}
}

// SetMusicPitch - Set the pitch (frequency) for the currently playing music
func (a *App) SetMusicPitch(pitch float32) {
	if a.IsMusicReady() {
		rl.SetMusicPitch(*a.music, pitch)
	}
}

// GetMusicTimeLength - Get the total duration of the music stream (in seconds)
func (a *App) GetMusicTimeLength(music Music) float32 {
	if !a.IsMusicReady() {
		return 0
	}
	return rl.GetMusicTimeLength(*a.music)
}

// GetMusicTimePlayed - Get the current time played in the music stream (in seconds)
func (a *App) GetMusicTimePlayed(music Music) float32 {
	if !a.IsMusicReady() {
		return 0
	}
	return rl.GetMusicTimePlayed(*a.music)
}

// SeekMusicStream - Seek to a specific position in the music stream (in seconds)
func (a *App) SeekMusicStream(music Music, position float32) {
	if a.IsMusicReady() {
		rl.SeekMusicStream(*a.music, position)
	}
}

// MusicFadeIn - Gradually fade in the music volume from a starting value over a specified duration
func (a *App) MusicFadeIn(from float32, duration time.Duration) {
	if a.IsMusicReady() {
		a.SetMasterVolume(from)
		go func() {
			steps := int(duration.Seconds() * float64(StepsPerSecond))

			stepDuration := duration / time.Duration(steps)
			incrementValue := 1.0 / float64(steps)

			for i := 0; i <= steps && a.GetMasterVolume() <= MaxMasterVolume-float32(incrementValue); i++ {
				a.SetMasterVolume(a.GetMasterVolume() + float32(incrementValue))
				time.Sleep(stepDuration)
			}
		}()
	}
}

// MusicFadeOut - Gradually fade out the music volume from a starting value over a specified duration
func (a *App) MusicFadeOut(from float32, duration time.Duration) {
	if a.IsMusicReady() {
		a.SetMasterVolume(from)

		go func() {
			steps := int(duration.Seconds() * float64(StepsPerSecond))

			stepDuration := duration / time.Duration(steps)
			incrementValue := 1.0 / float64(steps)

			for i := 0; i <= steps && a.GetMasterVolume() >= float32(incrementValue); i++ {
				a.SetMasterVolume(a.GetMasterVolume() - float32(incrementValue))
				time.Sleep(stepDuration)
			}
		}()
	}
}

// SwitchMusic - Provide a smooth transition between the current and the new music using FadeIn / FadeOut effects
func (a *App) SwitchMusic(loc ResourceLocator, duration time.Duration) {
	if a.IsMusicPlaying() {
		a.MusicFadeOut(a.GetMasterVolume(), duration)
		time.Sleep(duration)
	}
	a.PlayMusic(loc)
	a.MusicFadeIn(a.GetMasterVolume(), 5*time.Second)
}

// TODO spatial sounds (play with pan & pitch)
