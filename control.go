package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlVerbColor      = Green
	ControlVerbHoverColor = BrigthGreen
	ControlActionColor    = Cyan
)

// ControlPane is the screen control pane that shows the action, verbs and inventory.
type ControlPane struct {
	Enabled bool

	actionVerb Verb
	actionArg1 RoomItem
}

// Verb is a type that represents the action verb.
type Verb string

const (
	VerbOpen    Verb = "Open"
	VerbClose   Verb = "Close"
	VerbPush    Verb = "Push"
	VerbPull    Verb = "Pull"
	VerbWalkTo  Verb = "Walk to"
	VerbPickUp  Verb = "Pick up"
	VerbTalkTo  Verb = "Talk to"
	VerbGive    Verb = "Give"
	VerbUse     Verb = "Use"
	VerbLookAt  Verb = "Look at"
	VerbTurnOn  Verb = "Turn on"
	VerbTurnOff Verb = "Turn off"
)

// Draw renders the control panel in the viewport.
func (p *ControlPane) Draw(a *App) {
	if p.Enabled {
		p.drawActionVerb(a, VerbOpen, 0, 0)
		p.drawActionVerb(a, VerbClose, 0, 1)
		p.drawActionVerb(a, VerbPush, 0, 2)
		p.drawActionVerb(a, VerbPull, 0, 3)

		p.drawActionVerb(a, VerbWalkTo, 1, 0)
		p.drawActionVerb(a, VerbPickUp, 1, 1)
		p.drawActionVerb(a, VerbTalkTo, 1, 2)
		p.drawActionVerb(a, VerbGive, 1, 3)

		p.drawActionVerb(a, VerbUse, 2, 0)
		p.drawActionVerb(a, VerbLookAt, 2, 1)
		p.drawActionVerb(a, VerbTurnOn, 2, 2)
		p.drawActionVerb(a, VerbTurnOff, 2, 3)

		p.drawActionLine(a)
	}
}

func (p *ControlPane) drawActionVerb(a *App, verb Verb, col, row int) {
	x := 2 + col*ScreenWidth/6
	y := ViewportHeight + (row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	color := ControlVerbColor
	if a.MouseIsInto(NewRect(x, y, w, h)) {
		color = ControlVerbHoverColor
	}

	DrawDefaultText(string(verb), NewPos(x, y), AlignLeft, color)
}

func (p *ControlPane) drawActionLine(a *App) {
	if p.actionVerb == "" {
		p.actionVerb = VerbWalkTo
	}
	pos := NewPos(ScreenWidth/2, ViewportHeight)
	action := string(p.actionVerb)
	if room := a.room; room != nil {
		item := room.ItemAt(a.MousePosition())
		if item != nil {
			action = action + " " + item.Name()
		}
	}

	DrawDefaultText(action, pos, AlignCenter, ControlActionColor)
}

func (p *ControlPane) processControlInputs(a *App) {
	if a.ego != nil && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mouseClick := a.MousePosition()
		if ViewportRect.Contains(mouseClick) {
			p.processViewportClick(a, mouseClick)
		}
	}
}

func (p *ControlPane) processViewportClick(a *App, click Position) {
	if p.actionVerb == "" {
		p.actionVerb = VerbWalkTo
	}
	if room := a.room; room != nil {
		p.actionArg1 = room.ItemAt(click)
	}
	switch p.actionVerb {
	case VerbWalkTo:
		if p.actionArg1 != nil {
			// TODO: the item might be an actor
			a.Do(ActorWalkToObject{
				ActorID:  a.ego.name,
				ObjectID: p.actionArg1.Name(),
			})
		} else {
			a.Do(ActorWalkToPosition{
				ActorID:  a.ego.name,
				Position: NewPos(click.X, click.Y),
			})
		}
	}
	p.actionArg1 = nil
	p.actionVerb = VerbWalkTo
}

// EnableControlPanel is a command that will enable or disable the control panel.
type EnableControlPanel struct {
	Enable bool
}

func (cmd EnableControlPanel) Execute(app *App, done *Promise) {
	app.control.Enabled = cmd.Enable
	done.Complete()
}
