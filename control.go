package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlVerbColor      = Green
	ControlVerbHoverColor = BrigthGreen
	ControlActionColor    = Cyan
)

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

// VerbSlot is a slot in the control panel that holds a verb.
type VerbSlot struct {
	Verb Verb
	Row  int
	Col  int
}

// Draw renders the verb slot in the control pane.
func (s VerbSlot) Draw(a *App) {
	rect := s.Rect()
	color := ControlVerbColor
	if a.MouseIsInto(rect) {
		color = ControlVerbHoverColor
	}

	DrawDefaultText(string(s.Verb), rect.Pos, AlignLeft, color)
}

// Rect returns the rectangle of the verb slot in the screen.
func (v VerbSlot) Rect() Rectangle {
	x := 2 + v.Col*ScreenWidth/6
	y := ViewportHeight + (v.Row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize
	return NewRect(x, y, w, h)
}

// ControlPane is the screen control pane that shows the action, verbs and inventory.
type ControlPane struct {
	Enabled bool

	verbs      []VerbSlot
	actionVerb Verb
	actionArg1 RoomItem
}

// Init initializes the control pane.
func (p *ControlPane) Init() {
	p.Enabled = true
	p.verbs = []VerbSlot{
		{Verb: VerbOpen, Row: 0, Col: 0},
		{Verb: VerbClose, Row: 1, Col: 0},
		{Verb: VerbPush, Row: 2, Col: 0},
		{Verb: VerbPull, Row: 3, Col: 0},

		{Verb: VerbWalkTo, Row: 0, Col: 1},
		{Verb: VerbPickUp, Row: 1, Col: 1},
		{Verb: VerbTalkTo, Row: 2, Col: 1},
		{Verb: VerbGive, Row: 3, Col: 1},

		{Verb: VerbUse, Row: 0, Col: 2},
		{Verb: VerbLookAt, Row: 1, Col: 2},
		{Verb: VerbTurnOn, Row: 2, Col: 2},
		{Verb: VerbTurnOff, Row: 3, Col: 2},
	}
}

// Draw renders the control panel in the viewport.
func (p *ControlPane) Draw(a *App) {
	if p.Enabled {
		for _, v := range p.verbs {
			v.Draw(a)
		}
		p.drawActionLine(a)
	}
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
		if ControlPaneRect.Contains(mouseClick) {
			p.processControlPaneClick(a, mouseClick)
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
	case VerbLookAt:
		if p.actionArg1 != nil {
			a.Do(ActorLookAtObject{
				ActorID:  a.ego.name,
				ObjectID: p.actionArg1.Name(),
			})
		}
	}
	p.actionArg1 = nil
	p.actionVerb = VerbWalkTo
}

func (p *ControlPane) processControlPaneClick(_ *App, click Position) {
	for _, v := range p.verbs {
		if v.Rect().Contains(click) {
			p.actionVerb = v.Verb
			return
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
