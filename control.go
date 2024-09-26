package pctk

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	DefaultControlButtonColor             = Green
	SelectedOrSuggestedControlButtonColor = BrigthGreen
	EgoIntentionColor                     = Cyan
	DefaultInventoryItemColor             = Blue

	// ButtonsPerRow specifies the number of buttons to display in each row of the control rendering grid.
	ButtonsPerRow = 3
)

// ControlButton represents a button for interacting with  in the game,
// including its associated Verb, column, and row in a grid layout.
type ControlButton struct {
	Verb *Verb
	Col  int
	Row  int
}

// Bounds is implemented to satisfy the Interactable interface.
func (cb *ControlButton) Bounds() Rectangle {
	x := 2 + cb.Col*ScreenWidth/6
	y := ViewportHeight + (cb.Row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	return NewRect(x, y, w, h)
}

// Description is implemented to satisfy the Interactable interface.
func (cb *ControlButton) Description() string {
	return cb.Verb.Description
}

var (
	// A map to hold ControlButtons indexed by their VerbType
	buttons = make(map[VerbType]*ControlButton)
)

// Initializes the control buttons
func init() {
	for i, verb := range Verbs {
		buttons[verb.Type] = &ControlButton{
			Verb: verb,
			Col:  i % ButtonsPerRow,
			Row:  i / ButtonsPerRow,
		}
	}
}

func (a *App) drawControlPanel() {
	if a.controlPanelEnabled {

		a.drawControlButtons()
		a.drawEgoCurrentIntention()
		a.drawInventory()
	}

}

func (a *App) drawControlButtons() {
	for _, button := range buttons {
		a.drawControlButton(button, DefaultControlButtonColor)
	}
}

func (a *App) drawControlButton(cb *ControlButton, color Color) {
	rectangle := cb.Bounds()
	if a.MouseIsInto(rectangle) {
		color = SelectedOrSuggestedControlButtonColor
	}

	DrawDefaultText(cb.Verb.Description, NewPos(rectangle.Pos.X, rectangle.Pos.Y), AlignLeft, color)
}

func (a *App) drawEgoCurrentIntention() {
	ego := a.ego
	targetDescription := ""
	description := ego.String()
	for _, i := range a.Interactables() {
		if a.MouseIsInto(i.Bounds()) {
			switch i.(type) {
			case *Object:
				a.drawControlButton(buttons[LookAt], SelectedOrSuggestedControlButtonColor)
				targetDescription = i.Description()
			case *ControlButton:
				continue
			case *InventoryItem:
				a.drawControlButton(buttons[Use], SelectedOrSuggestedControlButtonColor)
				if description == "" {
					description = VerbUse.Description
				}
				targetDescription = i.Description()
			case *Actor:
				a.drawControlButton(buttons[TalkTo], SelectedOrSuggestedControlButtonColor)
				if description == "" {
					description = VerbTalkTo.Description
				}
				targetDescription = i.Description()
			}
			break
		}
	}

	if description == "" {
		description = VerbWalkTo.Description
	}

	if targetDescription != "" {
		description = fmt.Sprintf("%s the %s", description, targetDescription)
	}

	pos := NewPos(ScreenWidth/2, ViewportHeight)
	DrawDefaultText(description, pos, AlignCenter, EgoIntentionColor)

}

func (a *App) drawInventory() {
	ego := a.ego
	// TODO missing scroll
	for _, i := range ego.Inventory().items {
		r := i.Bounds()
		if a.MouseIsInto(r) {
			DrawDefaultText(i.Description(), NewPos(r.Pos.X, r.Pos.Y), AlignCenter, SelectedOrSuggestedControlButtonColor)
		} else {
			DrawDefaultText(i.Description(), NewPos(r.Pos.X, r.Pos.Y), AlignCenter, DefaultInventoryItemColor)
		}
	}
}

func (a *App) processControlInputs() {
	ego := a.ego
	if ego != nil && ego.actor != nil && rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		mouseClick := a.MousePosition()
		for _, i := range a.Interactables() {
			if a.MouseIsInto(i.Bounds()) {
				switch target := i.(type) {
				case *Object:
					if ego.verb != nil {
						a.Do(EgoInteraction{
							Object: target,
							Verb:   ego.verb,
						})
					} else {
						a.Do(ActorWalkToPosition{
							ActorID:  ego.actor.name,
							Position: NewPos(mouseClick.X, int(a.ego.actor.pos.Y)),
						})
					}
				case *ControlButton:
					ego.verb = target.Verb
					continue
				case *InventoryItem:
					a.Do(EgoInteraction{
						Object: target.object,
						Verb:   ego.verb,
					})
					continue
				case *Actor:
					// TODO actor dialogtree
					continue

				}

				// default
			} else if RoomViewport.Contains(mouseClick) {
				a.Do(ActorWalkToPosition{
					ActorID:  ego.actor.name,
					Position: NewPos(mouseClick.X, int(a.ego.actor.pos.Y)),
				})
			}
		}

	}
}

// EnableControlPanel is a command that will enable or disable the control panel.
type EnableControlPanel struct {
	Enable bool
}

func (cmd EnableControlPanel) Execute(app *App, done *Promise) {
	app.controlPanelEnabled = cmd.Enable
	done.Complete()
}
