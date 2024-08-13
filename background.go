package pctk

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Background represents a background image.
type Background struct {
	img *rl.Image
	tex *rl.Texture2D
}

// BackgroundFromImage creates a new Background from an image.
func BackgroundFromImage(img *rl.Image) *Background {
	if img.Width < ScreenWidth || img.Height < ScreenHeightScene {
		log.Fatal("Background image is too small")
	}
	return &Background{img: img, tex: nil}
}

// Texture returns the texture of the background, loading it if necessary.
func (b *Background) Texture() rl.Texture2D {
	if b.tex == nil {
		b.tex = new(rl.Texture2D)
		*b.tex = rl.LoadTextureFromImage(b.img)
	}
	return *b.tex
}

// SetBackground sets the background of the application.
func (a *App) SetBackground(loc ResourceLocator) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	bg, err := a.cat.GetBackground(loc)
	if err != nil {
		return err
	}
	a.background = bg
	return nil
}

func (a *App) drawBackgroud() {
	if a.background != nil {
		rl.DrawTexture(a.background.Texture(), 0, 0, rl.White)
	}
}
