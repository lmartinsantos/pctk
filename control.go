package pctk

import "fmt"

type ActionName int

// ActionName represents a verb or command that the player can perform in the game.
const (
	Open ActionName = iota
	Close
	Push
	Pull
	WalkTo
	PickUp
	TalkTo
	Give
	Use
	LookAt
	TurnOn
	TurnOff
)

// Action represents an interactive action in the game including where is rendered.
type Action struct {
	ActionName  ActionName
	Description string
	Col         int
	Row         int
}

var (
	ControlVerbColor                 = Green
	ControlVerbHoverOrSuggestedColor = BrigthGreen
	ControlActionColor               = Cyan

	// Actions
	ActionOpen  = &Action{ActionName: Open, Description: "Open", Col: 0, Row: 0}
	ActionClose = &Action{ActionName: Close, Description: "Close", Col: 0, Row: 1}
	ActionPush  = &Action{ActionName: Push, Description: "Push", Col: 0, Row: 2}
	ActionPull  = &Action{ActionName: Pull, Description: "Pull", Col: 0, Row: 3}

	ActionWalkTo = &Action{ActionName: WalkTo, Description: "Walk to", Col: 1, Row: 0}
	ActionPickUp = &Action{ActionName: PickUp, Description: "Pick up", Col: 1, Row: 1}
	ActionTalkTo = &Action{ActionName: TalkTo, Description: "Talk to", Col: 1, Row: 2}
	ActionGive   = &Action{ActionName: Give, Description: "Give", Col: 1, Row: 3}

	ActionUse     = &Action{ActionName: Use, Description: "Use", Col: 2, Row: 0}
	ActionLookAt  = &Action{ActionName: LookAt, Description: "Look at", Col: 2, Row: 1}
	ActionTurnOn  = &Action{ActionName: TurnOn, Description: "Turn on", Col: 2, Row: 2}
	ActionTurnOff = &Action{ActionName: TurnOff, Description: "Turn off", Col: 2, Row: 3}

	DefaultAction = ActionWalkTo
)

func (a *App) drawControlPanel() {

	actions := []*Action{
		ActionOpen, ActionClose, ActionPush, ActionPull,
		ActionWalkTo, ActionPickUp, ActionTalkTo, ActionGive,
		ActionUse, ActionLookAt, ActionTurnOn, ActionTurnOff,
	}

	for _, action := range actions {
		a.drawAction(action, ControlVerbColor)
	}
}

func (a *App) drawAction(action *Action, color Color) {
	x := 2 + action.Col*ScreenWidth/6
	y := ScreenHeightScene + (action.Row+1)*FontDefaultSize
	w := ScreenWidth / 6
	h := FontDefaultSize

	if a.MouseIsInto(NewRect(x, y, w, h)) {
		color = ControlVerbHoverOrSuggestedColor
	}

	a.drawDefaultText(action.Description, NewPos(x, y), AlignLeft, color)
}

func (a *App) drawEgoAction() {
	action := DefaultAction.Description
	if a.egoActionSelected != nil {
		action = a.egoActionSelected.Description
	}
	// check if mouse is hovering an object
	for _, o := range a.objects {
		size := o.anim.getAnimationSize(a)
		if a.MouseIsInto(NewRect(o.pos.X, o.pos.Y, size.W, size.H)) {
			action = fmt.Sprintf("%s %s", action, o.name)
			a.drawAction(ActionLookAt, ControlVerbHoverOrSuggestedColor)
			break
		}
	}

	// TODO  hovering actors (discarding ego)

	pos := NewPos(ScreenWidth/2, ScreenHeightScene)
	a.drawDefaultText(action, pos, AlignCenter, ControlVerbColor)
}
