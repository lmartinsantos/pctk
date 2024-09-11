package pctk

type Actor struct {
	name string

	standH *Animation
	walkH  *Animation

	lookAt Direction
	pos    Position
	act    action
}

func NewActor(name string) *Actor {
	return &Actor{name: name}
}

func (a *Actor) WithStandH(anim *Animation) *Actor {
	a.standH = anim
	return a
}

func (a *Actor) WithWalkH(anim *Animation) *Actor {
	a.walkH = anim
	return a
}

func (a *Actor) At(pos Position) *Actor {
	a.pos = pos
	return a
}

func (a *Actor) Looking(pos Position) *Actor {
	// TODO: adjust respect vector size
	if a.pos.X < pos.X {
		a.lookAt = DirRight
	} else if a.pos.X > pos.X {
		a.lookAt = DirLeft
	}
	return a
}

func (a *Actor) Stand() *Actor {
	a.act = func(app *App, c *Actor) (completed bool) {
		switch c.lookAt {
		case DirRight:
			if a.standH != nil {
				a.standH.draw(app, c.pos, false)
			}
		case DirLeft:
			if a.standH != nil {
				a.standH.draw(app, c.pos, true)
			}
		}
		return false
	}
	return a
}

func (a *Actor) WalkTo(pos Position) *Actor {
	a.act = func(app *App, c *Actor) (completed bool) {
		flip := false
		if c.pos.X < pos.X {
			c.pos.X++
			a.lookAt = DirRight
		} else if c.pos.X > pos.X {
			flip = true
			c.pos.X--
			a.lookAt = DirLeft
		}
		if c.pos.Y < pos.Y {
			c.pos.Y++
		} else if c.pos.Y > pos.Y {
			c.pos.Y--
		}

		if a.walkH != nil {
			a.walkH.draw(app, c.pos, flip)
		}

		return c.pos == pos
	}
	return a
}

func (a *Actor) draw(app *App) {
	if a.act == nil {
		if a.standH != nil {
			a.standH.draw(app, a.pos, false)
		}
		return
	}
	if a.act(app, a) {
		a.Stand()
	}
}

type action func(*App, *Actor) (completed bool)

func (a *App) ShowActor(loc ResourceLocator, pos Position) *Actor {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	actor := a.res.LoadActor(loc).At(pos)
	a.actors = append(a.actors, actor)
	return actor
}

func (a *App) drawActors() {
	for _, actor := range a.actors {
		actor.draw(a)
	}
}
