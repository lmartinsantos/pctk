package pctk

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlVerbColor                 = Green
	ControlVerbHoverOrSuggestedColor = BrigthGreen
	ControlEgoVerbColor              = Cyan
)

func (a *App) drawControlPanel() {
	if a.controlPanelEnabled {
		for _, verb := range Verbs {
			a.drawVerb(verb, ControlVerbColor)
		}
		a.drawEgoVerb()
	}
}

func (a *App) drawVerb(Verb *Verb, color Color) {
	x := 2 + Verb.Col*ScreenWidth/6
	y := ViewportHeight + (Verb.Row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	if a.MouseIsInto(NewRect(x, y, w, h)) {
		color = ControlVerbHoverOrSuggestedColor
	}

	a.drawDefaultText(Verb.Description, NewPos(x, y), AlignLeft, color)
}

func (a *App) drawEgoVerb() {
	description := DefaultVerb.Description
	if a.egoVerbSelected != nil {
		description = a.egoVerbSelected.Description
	}
	// check if mouse is hovering an object
	for _, o := range a.objects {
		size := o.FrameSize()
		if a.MouseIsInto(NewRect(o.pos.X, o.pos.Y, size.W, size.H)) {
			description = fmt.Sprintf("%s %s", description, o.name)
			a.drawVerb(VerbLookAt, ControlVerbHoverOrSuggestedColor)
			break
		}
	}

	// TODO  hovering actors (discarding ego)

	pos := NewPos(ScreenWidth/2, ViewportHeight)
	a.drawDefaultText(description, pos, AlignCenter, ControlEgoVerbColor)
}

func (a *App) processControlInputs() {
	if a.ego != nil && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mouseClick := a.MousePosition()
		if RoomViewport.Contains(mouseClick) {
			// TODO missing check Verb / control selected
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
