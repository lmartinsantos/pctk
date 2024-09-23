package pctk

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	ControlVerbColor                 = Green
	ControlVerbHoverOrSuggestedColor = BrigthGreen
	ControlEgoVerbColor              = Cyan
	ControlInventoryColor            = Blue
)

func (a *App) drawControlPanel() {
	if a.controlPanelEnabled {

		a.drawVerbs()
		a.drawEgoVerb()
		a.drawInventory()
	}
}

func (a *App) drawVerbs() {
	for _, verb := range Verbs {
		a.drawVerb(verb, ControlVerbColor)
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
		if a.MouseIsInto(o.Rectangle()) && !o.HasClass(ClassUntouchable) {
			targetDescription = o.name
			a.drawVerb(VerbLookAt, ControlVerbHoverOrSuggestedColor)
			break
		}
	}

	// check if mouse is hovering an object in the inventory
	var row int
	var fromInventory bool
	for _, o := range ego.actor.inventory {
		r := getInventoryItemRectangle(row)
		if a.MouseIsInto(r) {
			targetDescription = o.name
			a.drawVerb(VerbUse, ControlVerbHoverOrSuggestedColor)
			fromInventory = true
			break
		}
	}

	// TODO  hovering actors (discarding ego)

	description := ego.String(fromInventory)
	if targetDescription != "" {
		description = fmt.Sprintf("%s the %s", description, targetDescription)
	}

	pos := NewPos(ScreenWidth/2, ViewportHeight)
	a.drawDefaultText(description, pos, AlignCenter, ControlEgoVerbColor)
}

func (a *App) drawInventory() {
	ego := a.ego
	// TODO missing scroll
	var row int
	for _, o := range ego.actor.inventory {
		r := getInventoryItemRectangle(row)
		if a.MouseIsInto(r) {
			a.drawDefaultText(o.name, NewPos(r.Pos.X, r.Pos.Y), AlignCenter, ControlVerbHoverOrSuggestedColor)
		} else {
			a.drawDefaultText(o.name, NewPos(r.Pos.X, r.Pos.Y), AlignCenter, ControlInventoryColor)
		}
	}
}

func getInventoryItemRectangle(row int) Rectangle {
	x := 2 + 4*ScreenWidth/6
	y := ViewportHeight + (row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	return NewRect(x, y, w, h)
}

func (a *App) processControlInputs() {
	ego := a.ego
	if ego != nil && ego.actor != nil && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mouseClick := a.MousePosition()
		if RoomViewport.Contains(mouseClick) {
			// TODO missing check ego verb / object source
			var target *Object
			for _, o := range a.objects {
				if a.MouseIsInto(o.Rectangle()) && !o.HasClass(ClassUntouchable) {
					target = o
					break
				}
			}
			if target != nil && ego.verb != nil {
				a.Do(ObjectOnVerb{
					Object: target,
					Verb:   ego.verb,
				})
				ego.verb = nil
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
					ego.verb = verb
					return
				}
			}

			// TODO check inventory (setObject)
			var row int
			for _, o := range ego.actor.inventory {
				r := getInventoryItemRectangle(row)
				if a.MouseIsInto(r) && ego.verb != nil {
					a.Do(ObjectOnVerb{
						Object: o,
						Verb:   ego.verb,
					})
					ego.verb = nil
					return
				}
			}
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
