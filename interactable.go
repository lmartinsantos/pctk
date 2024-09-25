package pctk

// Interactable represents an element that the player can interact with as Actors, Objects, ControlButtons or Inventory items.
type Interactable interface {
	// Returns the bounding area for interaction
	Bounds() Rectangle
	// Description related to the interatable item
	Description() string
}

func (a *App) Interactables() []Interactable {
	var interactables []Interactable

	// Add room objects
	if a.room != nil {
		for _, obj := range a.room.Objects() {
			if !obj.HasClass(ClassUntouchable) {
				interactables = append(interactables, obj)
			}

		}
	}

	// Add actors (filtering ego)
	ego := &Actor{name: ""}
	if a.ego.actor != nil {
		ego = a.ego.actor
	}
	for _, actor := range a.actors {
		if actor.Description() != ego.Description() {
			interactables = append(interactables, actor)
		}
	}

	// Add inventory items if necessary
	if a.ego != nil && a.ego.actor != nil {
		inventory := a.ego.Inventory()
		for _, item := range inventory.items {
			interactables = append(interactables, item)
		}
	}

	// Add control buttons
	for _, button := range buttons {
		interactables = append(interactables, button)
	}

	return interactables
}
