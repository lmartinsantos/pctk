package pctk

type VerbType int

// VerbType represents a verb or command that the player can perform in the game.
const (
	Open VerbType = iota
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

// Verb represents an interactive Verb in the game including where is rendered.
type Verb struct {
	Type        VerbType
	Description string
	Col         int
	Row         int
}

var (
	VerbOpen  = &Verb{Type: Open, Description: "Open", Col: 0, Row: 0}
	VerbClose = &Verb{Type: Close, Description: "Close", Col: 0, Row: 1}
	VerbPush  = &Verb{Type: Push, Description: "Push", Col: 0, Row: 2}
	VerbPull  = &Verb{Type: Pull, Description: "Pull", Col: 0, Row: 3}

	VerbWalkTo = &Verb{Type: WalkTo, Description: "Walk to", Col: 1, Row: 0}
	VerbPickUp = &Verb{Type: PickUp, Description: "Pick up", Col: 1, Row: 1}
	VerbTalkTo = &Verb{Type: TalkTo, Description: "Talk to", Col: 1, Row: 2}
	VerbGive   = &Verb{Type: Give, Description: "Give", Col: 1, Row: 3}

	VerbUse     = &Verb{Type: Use, Description: "Use", Col: 2, Row: 0}
	VerbLookAt  = &Verb{Type: LookAt, Description: "Look at", Col: 2, Row: 1}
	VerbTurnOn  = &Verb{Type: TurnOn, Description: "Turn on", Col: 2, Row: 2}
	VerbTurnOff = &Verb{Type: TurnOff, Description: "Turn off", Col: 2, Row: 3}

	DefaultVerb = VerbWalkTo
	Verbs       = []*Verb{
		VerbOpen, VerbClose, VerbPush, VerbPull,
		VerbWalkTo, VerbPickUp, VerbTalkTo, VerbGive,
		VerbUse, VerbLookAt, VerbTurnOn, VerbTurnOff,
	}
)
