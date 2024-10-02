package pctk

import "time"

// ShowDialog is a command that will show a dialog with the given text.
type ShowDialog struct {
	Text     string
	Position Position
	Color    Color
	Speed    float32
}

func (cmd ShowDialog) Execute(app *App, done *Promise) {
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
