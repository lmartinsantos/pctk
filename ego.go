package pctk

import "fmt"

// Ego represents the main character controlled by the player holding references to the current action being performed.
type Ego struct {
	actor  *Actor
	verb   *Verb
	source *Object // TODO review how to represent complex actions (Give X to Y, Use X with Y) Y may be an Actor or an Object
}

func (e *Ego) String(fromInventory bool) string {
	description := DefaultVerb.Description
	if fromInventory {
		description = DefaultInventoryVerb.Description
	}
	source := ""
	if e != nil && e.source != nil {
		source = e.source.name
	}

	if e != nil && e.verb != nil {
		description = e.verb.Description
		switch e.verb.Type {
		case Give:
			if source != "" {
				description = fmt.Sprintf("%s %s to", e.verb.Description, source)
			}
		case Use:
			if source != "" {
				description = fmt.Sprintf("%s %s with", e.verb.Description, source)
			}
		}
	}

	return description
}

func (e *Ego) Clear() {
	e.actor, e.verb = nil, nil
}
