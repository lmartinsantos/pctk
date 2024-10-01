package pctk

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlActionColor         = Cyan
	ControlActionOngoingColor  = BrigthCyan
	ControlInventoryColor      = Magenta
	ControlInventoryHoverColor = BrigthMagenta
	ControlVerbColor           = Green
	ControlVerbHoverColor      = BrigthGreen
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
func (s *ActionSentence) Draw(app *App, hover RoomItem) {
	pos := NewPos(ScreenWidth/2, ViewportHeight)
	action := string(s.verb)
	color := ControlActionColor
	if s.fut != nil {
		// Ongoing action.
		if s.args[0] != nil {
			action = action + " " + s.args[0].Name()
		}
		color = ControlActionOngoingColor
	} else if hover != nil {
		action = action + " " + hover.Name()

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
	case VerbLookAt:
		s.lookAtItem(app, item)
	case VerbPickUp:
		s.pickupItem(app, item)
	case VerbWalkTo:
		s.walkToItem(app, item)
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
			s.pickupItem(app, item)
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

func (s *ActionSentence) pickupItem(app *App, item RoomItem) {
	s.verb = VerbPickUp
	s.args[0] = item
	s.fut = app.Do(ActorPickUpObject{
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

// ControlInventory is a screen control that shows the inventory.
type ControlInventory struct {
	slotsRect [6]Rectangle
}

// Draw renders the inventory in the control pane.
func (c *ControlInventory) Draw(app *App) {
	if app.ego == nil {
		return
	}
	mpos := app.MousePosition()
	for i, item := range app.ego.Inventory() {
		rect := c.slotsRect[i]
		color := ControlInventoryColor
		if rect.Contains(mpos) {
			color = ControlInventoryHoverColor
		}
		DrawDefaultText(item.Name(), rect.Pos, AlignLeft, color)
	}
}

// Init initializes the control inventory.
func (c *ControlInventory) Init() {
	arrowsWidth := 32
	for i := range c.slotsRect {
		c.slotsRect[i] = NewRect(
			2+3*ScreenWidth/6+arrowsWidth,
			ViewportHeight+FontDefaultSize*(i+1),
			2*ScreenWidth/6,
			FontDefaultSize,
		)
	}
}

// ObjectAt returns the object at the given position in the inventory box.
func (c *ControlInventory) ObjectAt(app *App, pos Position) *Object {
	if app.ego == nil {
		return nil
	}
	inv := app.ego.Inventory()
	for i, rect := range c.slotsRect {
		if rect.Contains(pos) {
			if i < len(inv) {
				return inv[i]
			}
			return nil
		}
	}
	return nil
}

// ControlPane is the screen control pane that shows the action, verbs and inventory.
type ControlPane struct {
	Enabled bool

	verbs  []VerbSlot
	action ActionSentence
	inv    ControlInventory
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
	p.inv.Init()
}

// Draw renders the control panel in the viewport.
func (p *ControlPane) Draw(app *App) {
	if p.Enabled {
		for _, v := range p.verbs {
			v.Draw(app)
		}
		hover := p.hover(app, app.MousePosition())
		p.action.Draw(app, hover)
		p.inv.Draw(app)
	}
}

func (p *ControlPane) hover(app *App, pos Position) RoomItem {
	var item RoomItem
	if ViewportRect.Contains(pos) && app.room != nil {
		item = app.room.ItemAt(pos)
	} else if ControlPaneRect.Contains(pos) {
		if obj := p.inv.ObjectAt(app, pos); obj != nil {
			item = obj
		}
	}
	return item
}

func (p *ControlPane) processControlInputs(app *App) {
	if app.ego == nil {
		return
	}
	pos := app.MousePosition()
	hover := p.hover(app, pos)
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		if ViewportRect.Contains(pos) {
			p.action.ProcessLeftClick(app, pos, hover)
		}
		if ControlPaneRect.Contains(pos) {
			p.processLeftClick(app, pos)
		}
	} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		if ViewportRect.Contains(pos) {
			p.action.ProcessRightClick(app, pos, hover)
		}
	}
}

func (p *ControlPane) processLeftClick(app *App, click Position) {
	for _, v := range p.verbs {
		if v.Rect().Contains(click) {
			p.action.Reset(v.Verb)
			return
		}
	}
	if obj := p.inv.ObjectAt(app, click); obj != nil {
		p.action.ProcessLeftClick(app, click, obj)
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
