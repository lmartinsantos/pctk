package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlVerbColor          = Green
	ControlVerbHoverColor     = BrigthGreen
	ControlActionColor        = Cyan
	ControlActionOngoingColor = BrigthCyan
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
	Alt  ObjectClass
}

// Draw renders the verb slot in the control pane.
func (s VerbSlot) Draw(a *App) {
	rect := s.Rect()
	color := ControlVerbColor
	if a.MouseIsInto(rect) {
		color = ControlVerbHoverColor
	}
	if room := a.room; room != nil {
		if item := room.ItemAt(a.MousePosition()); item != nil {
			if item.Class().Is(s.Alt) {
				color = ControlVerbHoverColor
			} else if s.Verb == VerbLookAt {
				color = ControlVerbHoverColor
			}
		}
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

// ActionSentence is a sentence that represents the action the player is doing.
type ActionSentence struct {
	verb Verb
	args [2]RoomItem
	fut  Future
}

// Draw renders the action sentence in the control pane.
func (s *ActionSentence) Draw(app *App) {
	pos := NewPos(ScreenWidth/2, ViewportHeight)
	action := string(s.verb)
	color := ControlActionColor
	if s.fut != nil {
		// Ongoing action.
		if s.args[0] != nil {
			action = action + " " + s.args[0].Name()
		}
		color = ControlActionOngoingColor
	} else {
		if room := app.room; room != nil {
			item := room.ItemAt(app.MousePosition())
			if item != nil {
				action = action + " " + item.Name()
			}
		}
	}
	DrawDefaultText(action, pos, AlignCenter, color)
}

// ProcessLeftClick processes a left click in the control pane.
func (s *ActionSentence) ProcessLeftClick(app *App, click Position, item RoomItem) {
	if item == nil {
		if s.verb == VerbWalkTo || s.fut != nil {
			s.walkToPos(app, click)
		}
		return
	} else if s.args[0] != nil {
		// TODO: handle the second argument
		s.Reset(VerbWalkTo)
		return
	}

	switch s.verb {
	case VerbWalkTo:
		s.walkToItem(app, item)
	case VerbLookAt:
		s.lookAtItem(app, item)
	default:
		s.Reset(VerbWalkTo)
	}
}

// ProcessRightClick processes a right click in the control pane.
func (s *ActionSentence) ProcessRightClick(app *App, click Position, item RoomItem) {
	if item != nil {
		// Execute quick action
		if item.Class().Is(ObjectClassPerson) {
			// TODO: do a talk to action
		} else if item.Class().Is(ObjectClassOpenable) {
			// TODO: do a open action
		} else if item.Class().Is(ObjectClassPickable) {
			// TODO: do a pick up action
		} else {
			s.lookAtItem(app, item)
		}
		return
	}
	// No item there. Only respond if current verb is walk to.
	if s.verb == VerbWalkTo {
		s.walkToPos(app, click)
	}
}

func (s *ActionSentence) lookAtItem(app *App, item RoomItem) {
	s.verb = VerbLookAt
	s.args[0] = item
	s.fut = app.Do(ActorLookAtObject{
		ActorID:  app.ego.name,
		ObjectID: item.Name(),
	}).AndThen(func(_ any) Future {
		return app.Do(SyncCommandFunc(func(app *App) { s.Reset(VerbWalkTo) }))
	})
}

func (s *ActionSentence) walkToItem(app *App, item RoomItem) {
	// TODO: the item might be an actor
	s.verb = VerbWalkTo
	s.args[0] = item
	s.fut = app.Do(ActorWalkToObject{
		ActorID:  app.ego.name,
		ObjectID: s.args[0].Name(),
	}).AndThen(func(_ any) Future {
		return app.Do(SyncCommandFunc(func(app *App) { s.Reset(VerbWalkTo) }))
	})
}

func (s *ActionSentence) walkToPos(app *App, click Position) {
	app.Do(ActorWalkToPosition{
		ActorID:  app.ego.name,
		Position: click,
	})
	s.Reset(VerbWalkTo)
}

// Reset resets the action sentence to the given verb.
func (s *ActionSentence) Reset(verb Verb) {
	s.verb = verb
	s.args[0] = nil
	s.args[1] = nil
	s.fut = nil
}

// ControlPane is the screen control pane that shows the action, verbs and inventory.
type ControlPane struct {
	Enabled bool

	verbs  []VerbSlot
	action ActionSentence
}

// Init initializes the control pane.
func (p *ControlPane) Init() {
	p.Enabled = true
	p.verbs = []VerbSlot{
		{Verb: VerbOpen, Row: 0, Col: 0, Alt: ObjectClassOpenable},
		{Verb: VerbClose, Row: 1, Col: 0},
		{Verb: VerbPush, Row: 2, Col: 0},
		{Verb: VerbPull, Row: 3, Col: 0},

		{Verb: VerbWalkTo, Row: 0, Col: 1},
		{Verb: VerbPickUp, Row: 1, Col: 1},
		{Verb: VerbTalkTo, Row: 2, Col: 1, Alt: ObjectClassPerson},
		{Verb: VerbGive, Row: 3, Col: 1},

		{Verb: VerbUse, Row: 0, Col: 2},
		{Verb: VerbLookAt, Row: 1, Col: 2},
		{Verb: VerbTurnOn, Row: 2, Col: 2},
		{Verb: VerbTurnOff, Row: 3, Col: 2},
	}
	p.action.Reset(VerbWalkTo)
}

// Draw renders the control panel in the viewport.
func (p *ControlPane) Draw(app *App) {
	if p.Enabled {
		for _, v := range p.verbs {
			v.Draw(app)
		}
		p.action.Draw(app)
	}
}

func (p *ControlPane) processControlInputs(app *App) {
	if app.ego == nil {
		return
	}
	pos := app.MousePosition()
	item := app.room.ItemAt(pos)
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		if ViewportRect.Contains(pos) {
			p.action.ProcessLeftClick(app, pos, item)
		}
		if ControlPaneRect.Contains(pos) {
			p.processControlPaneClick(app, pos)
		}
	} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		if ViewportRect.Contains(pos) {
			p.action.ProcessRightClick(app, pos, item)
		}
	}
}

func (p *ControlPane) processControlPaneClick(_ *App, click Position) {
	for _, v := range p.verbs {
		if v.Rect().Contains(click) {
			p.action.Reset(v.Verb)
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
