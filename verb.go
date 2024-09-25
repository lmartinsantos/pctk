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
	Default
)

// Verb represents an interactive actor's action in the game.
type Verb struct {
	Type           VerbType
	Description    string
	IsDitransitive bool // This flags means this verb requires both a direct object and an indirect object
}

var (
	VerbOpen  = &Verb{Type: Open, Description: "Open"}
	VerbClose = &Verb{Type: Close, Description: "Close"}
	VerbPush  = &Verb{Type: Push, Description: "Push"}
	VerbPull  = &Verb{Type: Pull, Description: "Pull"}

	VerbWalkTo = &Verb{Type: WalkTo, Description: "Walk to"}
	VerbPickUp = &Verb{Type: PickUp, Description: "Pick up"}
	VerbTalkTo = &Verb{Type: TalkTo, Description: "Talk to"}
	VerbGive   = &Verb{Type: Give, Description: "Give", IsDitransitive: true}

	VerbUse     = &Verb{Type: Use, Description: "Use"}
	VerbLookAt  = &Verb{Type: LookAt, Description: "Look at"}
	VerbTurnOn  = &Verb{Type: TurnOn, Description: "Turn on"}
	VerbTurnOff = &Verb{Type: TurnOff, Description: "Turn off"}

	Verbs = []*Verb{
		VerbOpen, VerbClose, VerbPush, VerbPull,
		VerbWalkTo, VerbPickUp, VerbTalkTo, VerbGive,
		VerbUse, VerbLookAt, VerbTurnOn, VerbTurnOff,
	}
)
