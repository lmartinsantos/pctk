package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlVerbColor      = Green
	ControlVerbHoverColor = BrigthGreen
	ControlActionColor    = Cyan
)

func (a *App) drawControlPanel() {
	if a.controlPanelEnabled {
		a.drawActionVerb("Open", 0, 0)
		a.drawActionVerb("Close", 0, 1)
		a.drawActionVerb("Push", 0, 2)
		a.drawActionVerb("Pull", 0, 3)

		a.drawActionVerb("Walk to", 1, 0)
		a.drawActionVerb("Pick up", 1, 1)
		a.drawActionVerb("Talk to", 1, 2)
		a.drawActionVerb("Give", 1, 3)

		a.drawActionVerb("Use", 2, 0)
		a.drawActionVerb("Look at", 2, 1)
		a.drawActionVerb("Turn on", 2, 2)
		a.drawActionVerb("Turn off", 2, 3)

		a.drawFullAction("Walk to") // TODO: do not hardcode this
	}
}

func (a *App) drawActionVerb(verb string, col, row int) {
	x := 2 + col*ScreenWidth/6
	y := ScreenHeightScene + (row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	color := ControlVerbColor
	if a.MouseIsInto(NewRect(x, y, w, h)) {
		color = ControlVerbHoverColor
	}

	a.drawDefaultText(verb, NewPos(x, y), AlignLeft, color)
}

func (a *App) drawFullAction(action string) {
	pos := NewPos(ScreenWidth/2, ScreenHeightScene)
	a.drawDefaultText(action, pos, AlignCenter, ControlActionColor)
}

func (a *App) processControlInputs() {
	if a.ego != nil && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mouseClick := a.MousePosition()
		if SceneViewport.Contains(mouseClick) {
			// TODO missing check action / control selected
			a.Do(ActorWalkToPosition{
				ActorName: a.ego.name,
				Position:  NewPos(mouseClick.X, a.ego.pos.Y),
			})
		}
	}
}

// EnableControlPanel is a command that will enable or disable the control panel.
type EnableControlPanel struct {
	Enable bool
}

func (cmd EnableControlPanel) Execute(app *App, done Promise) {
	app.controlPanelEnabled = cmd.Enable
	done.Complete()
}
