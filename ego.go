package pctk

import "fmt"

// Ego represents the main character controlled by the player holding references to the current action being performed.
type Ego struct {
	actor  *Actor
	verb   *Verb
	source *Object
}

func (e *Ego) Description(target string) string {
	description := DefaultVerb.Description
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

	if target != "" {
		description = fmt.Sprintf("%s the %s", description, target)
	}
	return description
}

func (e *Ego) setActor(actor *Actor) {
	e.actor = actor
}

func (e *Ego) setVerb(verb *Verb) {
	e.verb = verb
}

func (e *Ego) clearVerb() {
	e.verb = nil
}

func (e *Ego) clear() {
	e.clearVerb()
	e.actor, e.verb = nil, nil
}
