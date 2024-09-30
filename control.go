package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlVerbColor      = Green
	ControlVerbHoverColor = BrigthGreen
	ControlActionColor    = Cyan
)

type ControlPane struct {
	Enabled bool
	action  string
}

// Draw renders the control panel in the viewport.
func (p *ControlPane) Draw(a *App) {
	if p.Enabled {
		p.drawActionVerb(a, "Open", 0, 0)
		p.drawActionVerb(a, "Close", 0, 1)
		p.drawActionVerb(a, "Push", 0, 2)
		p.drawActionVerb(a, "Pull", 0, 3)

		p.drawActionVerb(a, "Walk to", 1, 0)
		p.drawActionVerb(a, "Pick up", 1, 1)
		p.drawActionVerb(a, "Talk to", 1, 2)
		p.drawActionVerb(a, "Give", 1, 3)

		p.drawActionVerb(a, "Use", 2, 0)
		p.drawActionVerb(a, "Look at", 2, 1)
		p.drawActionVerb(a, "Turn on", 2, 2)
		p.drawActionVerb(a, "Turn off", 2, 3)

		p.drawActionLine(a)
	}
}

func (p *ControlPane) drawActionVerb(a *App, verb string, col, row int) {
	x := 2 + col*ScreenWidth/6
	y := ViewportHeight + (row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	color := ControlVerbColor
	if a.MouseIsInto(NewRect(x, y, w, h)) {
		color = ControlVerbHoverColor
	}

	DrawDefaultText(verb, NewPos(x, y), AlignLeft, color)
}

func (p *ControlPane) drawActionLine(a *App) {
	if p.action == "" {
		p.action = "Walk to"
	}
	pos := NewPos(ScreenWidth/2, ViewportHeight)
	action := p.action
	if room := a.room; room != nil {
		item := room.ItemAt(a.MousePosition())
		if item != nil {
			action = p.action + " " + item.Name()
		}
	}

	DrawDefaultText(action, pos, AlignCenter, ControlActionColor)
}

func (p *ControlPane) processControlInputs(a *App) {
	if a.ego != nil && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mouseClick := a.MousePosition()
		if ViewportRect.Contains(mouseClick) {
			// TODO missing check action / control selected
			a.Do(ActorWalkToPosition{
				ActorID:  a.ego.name,
				Position: NewPos(mouseClick.X, mouseClick.Y),
			})
		}
	}
}

// EnableControlPanel is a command that will enable or disable the control panel.
type EnableControlPanel struct {
	Enable bool
}

func (cmd EnableControlPanel) Execute(app *App, done *Promise) {
	app.control.Enabled = cmd.Enable
	done.Complete()
}
