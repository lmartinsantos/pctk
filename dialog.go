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

// Dialog is a dialog that will be shown in the screen.
type Dialog struct {
	actor *Actor
	text  string
	pos   Position
	color Color
	speed float32

	completedAt time.Time
	done        *Promise
}

// Actor returns the actor that is speaking the dialog, or nil if it comes from a external voice.
func (d *Dialog) Actor() *Actor {
	return d.actor
}

// Draw will draw the dialog in the screen. It returns true if the dialog is completed.
func (d *Dialog) Draw() (completed bool) {
	if time.Now().After(d.completedAt) {
		return true
	}

	DrawDialogText(d.text, d.pos, d.color)
	return false
}

// ClearDialogsFrom will remove all dialogs from the given actor.
func (a *App) ClearDialogsFrom(actor *Actor) {
	dialogs := make([]Dialog, 0, len(a.dialogs))
	for _, d := range a.dialogs {
		if d.done == nil || d.Actor() != actor {
			dialogs = append(dialogs, d)
		}
	}
	a.dialogs = dialogs
}

func (a *App) drawDialogs() {
	dialogs := make([]Dialog, 0, len(a.dialogs))
	for _, d := range a.dialogs {
		if d.Draw() {
			d.done.Complete()
		} else {
			dialogs = append(dialogs, d)
		}
	}
	a.dialogs = dialogs
}
