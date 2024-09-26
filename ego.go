package pctk

import (
	"fmt"
)

// Ego represents the main character controlled by the player holding references to the current action being performed.
type Ego struct {
	actor  *Actor
	verb   *Verb
	source *Object // TODO review how to represent complex actions (Give X to Y, Use X with Y) Y may be an Actor or an Object
}

// Returns the current inventory of the ego
func (e *Ego) Inventory() *Inventory {
	if e.actor != nil {
		return e.actor.inventory
	}
	return NewInventory() // avoid segment faults
}

func (e *Ego) String() string {
	description := ""
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

// EgoAddObjectToInventory is a command that adds an item to the ego actor's inventory.
type EgoAddObjectToInventory struct {
	ObjectName string
}

func (cmd EgoAddObjectToInventory) Execute(app *App, done *Promise) {
	if ego := app.ego; ego.actor != nil {
		ego.Inventory().AddItem(app.room.ObjectByName(cmd.ObjectName))
	}

	done.Complete()
}

// EgoRemoveObjectFromInventory is a command that removes an object from the ego actor's inventory.
type EgoRemoveObjectFromInventory struct {
	ObjectName string
}

func (cmd EgoRemoveObjectFromInventory) Execute(app *App, done *Promise) {
	ego := app.ego
	if ego.actor != nil {
		ego.Inventory().RemoveItemByName(cmd.ObjectName)
	}

	done.Complete()
}

// EgoInteract is a command that runs the action script related to an interaction between Ego and an object or actor.
type EgoInteraction struct {
	Object *Object
	Verb   *Verb
}

func (cmd EgoInteraction) Execute(app *App, done *Promise) {
	state := cmd.Object.State()
	script := state.scripts[cmd.Verb.Type]
	if script == nil {
		script = state.scripts[Default]

	}
	script.luaInit(app)
	script.luaRun(app, done)
	app.ego.verb = nil
}
