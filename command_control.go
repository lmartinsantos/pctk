package pctk

// EnableControlPanel is a command that will enable or disable the control panel.
type EnableControlPanel struct {
	Enable bool
}

func (cmd EnableControlPanel) Execute(app *App, done *Promise) {
	app.control.Enabled = cmd.Enable
	done.Complete()
}

// EnableMouseCursor is a command that will enable or disable the mouse control.
type EnableMouseCursor struct {
	Enable bool
}

func (cmd EnableMouseCursor) Execute(app *App, done *Promise) {
	app.control.cursor.Enabled = cmd.Enable
	done.Complete()
}
