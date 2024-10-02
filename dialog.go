package pctk

import (
	"time"
)

const (
	// LettersPerSecond is the number of letters that an adult human could read easily per second.
	LettersPerSecond = 10
)

var (
	// DefaultDialogColor is the default color of a dialog.
	DefaultDialogColor = Magenta

	// DefaultDialogPosition is the default position of a dialog.
	DefaultDialogPosition = Position{X: 160, Y: 20}
)

type Dialog struct {
	text  string
	pos   Position
	color Color
	speed float32

	expiresAt time.Time
	done      *Promise
}

func (d *Dialog) draw() (expired bool) {
	if time.Now().After(d.expiresAt) {
		return true
	}

	DrawDialogText(d.text, d.pos, d.color)
	return false
}

func (a *App) drawDialogs() {
	dialogs := make([]Dialog, 0, len(a.dialogs))
	for _, d := range a.dialogs {
		if d.draw() {
			d.done.Complete()
		} else {
			dialogs = append(dialogs, d)
		}
	}
	a.dialogs = dialogs
}
