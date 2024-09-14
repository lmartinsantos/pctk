package pctk

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	DefaultActorSpeakDelay = 500 * time.Millisecond
)

type Actor struct {
	name string

	standH *Animation
	speakH *Animation
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

func (a *Actor) WithSpeakH(anim *Animation) *Actor {
	a.speakH = anim
	return a
}

func (a *Actor) WithWalkH(anim *Animation) *Actor {
	a.walkH = anim
	return a
}

func (a *Actor) stand() *Actor {
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

func (a *Actor) draw(app *App) {
	if a.act == nil {
		if a.standH != nil {
			a.standH.draw(app, a.pos, false)
		}
		return
	}
	if a.act(app, a) {
		a.stand()
	}
}

type action func(*App, *Actor) (completed bool)

// ActorShow is a command that will show an actor in the scene at the given position.
type ActorShow struct {
	ActorResource ResourceLocator
	ActorName     string
	Position      Position
	LookAt        Direction
}

func (cmd ActorShow) Execute(app *App, done Promise) {
	actor := app.res.LoadActor(cmd.ActorResource)
	actor.pos = cmd.Position
	actor.lookAt = cmd.LookAt
	actor.stand()
	app.actors[cmd.ActorName] = actor
	done.Complete()
}

// ActorLookAtPos is a command that will make an actor look at a given position.
type ActorLookAtPos struct {
	ActorName string
	Position  Position
}

func (cmd ActorLookAtPos) Execute(app *App, done Promise) {
	app.withActor(cmd.ActorName, func(a *Actor) {
		if a.pos.X < cmd.Position.X {
			a.lookAt = DirRight
		} else if a.pos.X > cmd.Position.X {
			a.lookAt = DirLeft
		}
	})
	done.Complete()
}

// ActorLookAtDirection is a command that will make an actor look at a given direction.
type ActorLookAtDirection struct {
	ActorName string
	Direction Direction
}

func (cmd ActorLookAtDirection) Execute(app *App, done Promise) {
	app.withActor(cmd.ActorName, func(a *Actor) {
		a.lookAt = cmd.Direction
	})
	done.Complete()
}

// ActorWalkToPosition is a command that will make an actor walk to a given position.
type ActorWalkToPosition struct {
	ActorName string
	Position  Position
}

func (cmd ActorWalkToPosition) Execute(app *App, done Promise) {
	app.withActor(cmd.ActorName, func(a *Actor) {
		a.act = func(app *App, c *Actor) (completed bool) {
			flip := false
			if c.pos.X < cmd.Position.X {
				c.pos.X++
				a.lookAt = DirRight
			} else if c.pos.X > cmd.Position.X {
				flip = true
				c.pos.X--
				a.lookAt = DirLeft
			}
			if c.pos.Y < cmd.Position.Y {
				c.pos.Y++
			} else if c.pos.Y > cmd.Position.Y {
				c.pos.Y--
			}

			if a.walkH != nil {
				a.walkH.draw(app, c.pos, flip)
			}
			if c.pos == cmd.Position {
				done.Complete()
			}
			return c.pos == cmd.Position
		}
	})
}

// ActorSpeak is a command that will make an actor speak the given text.
type ActorSpeak struct {
	ActorName string
	Text      string
	Delay     time.Duration
}

func (cmd ActorSpeak) Execute(app *App, done Promise) {
	if cmd.Delay == 0 {
		cmd.Delay = DefaultActorSpeakDelay
	}

	app.withActor(cmd.ActorName, func(a *Actor) {
		dialogDone := app.doNow(ShowDialog{
			Text:     cmd.Text,
			Position: a.pos.Above(50),
			Color:    rl.White,
			Speed:    1.0,
		})
		a.act = func(app *App, c *Actor) (completed bool) {
			switch a.lookAt {
			case DirRight:
				if a.speakH != nil {
					a.speakH.draw(app, c.pos, false)
				}
			case DirLeft:
				if a.speakH != nil {
					a.speakH.draw(app, c.pos, true)
				}
			}
			if dialogDone.IsCompleted() {
				done.CompleteAfter(nil, cmd.Delay)
				return true
			}
			return false
		}
	})
}

func (a *App) withActor(name string, f func(*Actor)) {
	actor, ok := a.actors[name]
	if !ok {
		return
	}
	f(actor)
}

func (a *App) drawActors() {
	for _, actor := range a.actors {
		actor.draw(a)
	}
}
