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

func (a *App) drawVerb(v *Verb, color Color) {
	rectangle := v.Rectangle()
	if a.MouseIsInto(rectangle) {
		color = ControlVerbHoverOrSuggestedColor
	}

	a.drawDefaultText(v.Description, NewPos(rectangle.Pos.X, rectangle.Pos.Y), AlignLeft, color)
}

func (a *App) drawEgoVerb() {
	ego := a.ego
	targetDescription := ""
	// check if mouse is hovering an object
	for _, o := range a.objects {
		if a.MouseIsInto(o.Rectangle()) {
			targetDescription = o.name
			a.drawVerb(VerbLookAt, ControlVerbHoverOrSuggestedColor)
			break
		}
	}
	// TODO  hovering actors (discarding ego)
	description := ego.String()
	if targetDescription != "" {
		description = fmt.Sprintf("%s the %s", description, targetDescription)
	}

	pos := NewPos(ScreenWidth/2, ViewportHeight)
	a.drawDefaultText(description, pos, AlignCenter, ControlEgoVerbColor)
}

func (a *App) processControlInputs() {
	ego := a.ego
	if ego != nil && ego.actor != nil && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mouseClick := a.MousePosition()
		if RoomViewport.Contains(mouseClick) {
			// TODO missing check ego verb / object source
			var target *Object
			for _, o := range a.objects {
				if a.MouseIsInto(o.Rectangle()) {
					target = o
					break
				}
			}
			if target != nil && ego.verb != nil {
				a.Do(ObjectOnAction{
					ObjectName: target.name,
					Verb:       ego.verb,
				})
			} else {
				a.Do(ActorWalkToPosition{
					ActorName: ego.actor.name,
					Position:  NewPos(mouseClick.X, a.ego.actor.pos.Y),
				})
			}
		} else {
			// check verbs
			for _, verb := range Verbs {
				if a.MouseIsInto(verb.Rectangle()) {
					ego.setVerb(verb)
					return
				}
			}

			// TODO check inventory (setObject)
		}

		// clean ego status
		ego.verb = nil
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
