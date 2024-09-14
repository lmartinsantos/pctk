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
	pos   Position
	color Color
	speed float32

	expiresAt time.Time
	done      Promise
}

func (a *App) drawDialogs() {
	dialogs := make([]Dialog, 0, len(a.dialogs))
	for _, d := range a.dialogs {
		if a.drawDialog(&d) {
			d.done.Complete()
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

	a.drawDialogText(d.text, d.pos, d.color)
	return false
}

// ShowDialog is a command that will show a dialog with the given text.
type ShowDialog struct {
	Text     string
	Position Position
	Color    Color
	Speed    float32
}

func (cmd ShowDialog) Execute(app *App, done Promise) {
	if cmd.Speed == 0 {
		cmd.Speed = 1
	}

	shownDuring := time.Duration(len(cmd.Text)/LettersPerSecond) * time.Second
	if shownDuring < 2*time.Second {
		shownDuring = 2 * time.Second
	}
	shownDuring /= time.Duration(cmd.Speed)
	expiresAt := time.Now().Add(shownDuring)

	dialog := Dialog{
		text:      cmd.Text,
		pos:       cmd.Position,
		color:     cmd.Color,
		speed:     cmd.Speed,
		expiresAt: expiresAt,
		done:      done,
	}
	app.dialogs = append(app.dialogs, dialog)
}
