package pctk

import (
	"time"
)

const (
	// LettersPerSecond is the number of letters that an adult human could read easily per second.
	LettersPerSecond = 10
)

type Dialog struct {
	text  string
	x, y  int32
	color Color
	speed float32

	expiresAt time.Time
	done      chan<- struct{}
}

func (a *App) ShowDialog(text string, x, y int32, col Color, speed float32) (done <-chan struct{}) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	shownDuring := time.Duration(len(text)/LettersPerSecond) * time.Second
	if shownDuring < 2*time.Second {
		shownDuring = 2 * time.Second
	}
	shownDuring /= time.Duration(speed)
	expiresAt := time.Now().Add(shownDuring)
	dchan := make(chan struct{})

	dialog := Dialog{
		text:      text,
		x:         x,
		y:         y,
		color:     col,
		speed:     speed,
		expiresAt: expiresAt,
		done:      dchan,
	}
	a.dialogs = append(a.dialogs, dialog)
	return dchan
}

func (a *App) drawDialogs() {
	dialogs := make([]Dialog, 0, len(a.dialogs))
	for _, d := range a.dialogs {
		if a.drawDialog(&d) {
			time.AfterFunc(time.Duration(d.speed)*time.Second, func() {
				close(d.done)
			})
		} else {
			dialogs = append(dialogs, d)
		}
	}
	a.dialogs = dialogs
}

func (a *App) drawDialog(d *Dialog) (expired bool) {
	if time.Now().After(d.expiresAt) {
		return true
	}

	a.drawDialogText(d.text, d.x, d.y, d.color)
	return false
}
