package pctk

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	DefaultActorSpeakDelay = 500 * time.Millisecond
)

var (
	DefaultActorPosition  = NewPos(160, 90)
	DefaultActorDirection = DirRight
)

type Actor struct {
	name string

	animStand map[Direction]*Animation
	animSpeak map[Direction]*Animation
	animWalk  map[Direction]*Animation

	lookAt Direction
	pos    Position
	act    action
}

func NewActor(name string) *Actor {
	return &Actor{
		name:      name,
		animStand: make(map[Direction]*Animation),
		animSpeak: make(map[Direction]*Animation),
		animWalk:  make(map[Direction]*Animation),
	}
}

func (a *Actor) WithAnimationStand(dir Direction, anim *Animation) *Actor {
	a.animStand[dir] = anim
	return a
}

func (a *Actor) WithAnimationSpeak(dir Direction, anim *Animation) *Actor {
	a.animSpeak[dir] = anim
	return a
}

func (a *Actor) WithAnimationWalk(dir Direction, anim *Animation) *Actor {
	a.animWalk[dir] = anim
	return a
}

func (a *Actor) stand(dir Direction) *Actor {
	a.lookAt = dir
	a.act = func(app *App, c *Actor) (completed bool) {
		if anim := a.animStand[dir]; anim != nil {
			anim.draw(app, c.pos)
		}
		return false
	}
	return a
}

func (a *Actor) draw(app *App) {
	if a.act == nil {
		a.stand(a.lookAt)
	}
	if a.act(app, a) {
		a.stand(a.lookAt)
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
	actor.stand(cmd.LookAt)
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
		a.stand(a.pos.DirectionTo(cmd.Position))
	})
	done.Complete()
}

// ActorStand is a command that will make an actor stand in the given direction.
type ActorStand struct {
	ActorName string
	Direction Direction
}

func (cmd ActorStand) Execute(app *App, done Promise) {
	app.withActor(cmd.ActorName, func(a *Actor) {
		a.stand(cmd.Direction)
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
			if anim := a.animWalk[a.lookAt]; anim != nil {
				anim.draw(app, c.pos)
			}

			if c.pos == cmd.Position {
				done.Complete()
				return true
			}

			// TODO: This implementation is totally naive. It doesn't take into account the
			// diagonal movement, the obstacles, the speed of the actor, etc.
			a.lookAt = c.pos.DirectionTo(cmd.Position)

			switch a.lookAt {
			case DirRight:
				c.pos.X++
			case DirLeft:
				c.pos.X--
			case DirUp:
				c.pos.Y--
			case DirDown:
				c.pos.Y++
			}
			return false
		}
	})
}

// ActorSpeak is a command that will make an actor speak the given text.
type ActorSpeak struct {
	ActorName string
	Text      string
	Delay     time.Duration
	Color     Color
}

func (cmd ActorSpeak) Execute(app *App, done Promise) {
	if cmd.Delay == 0 {
		cmd.Delay = DefaultActorSpeakDelay
	}

	if cmd.Color == rl.Blank {
		cmd.Color = rl.White
	}

	app.withActor(cmd.ActorName, func(a *Actor) {
		dialogDone := app.doNow(ShowDialog{
			Text:     cmd.Text,
			Position: a.pos.Above(50),
			Color:    cmd.Color,
			Speed:    1.0,
		})
		a.act = func(app *App, c *Actor) (completed bool) {
			if anim := a.animSpeak[a.lookAt]; anim != nil {
				anim.draw(app, c.pos)
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

// ActorSelectEgo is a command that will make an actor be the actor under player's control.
type ActorSelectEgo struct {
	// Using an empty ActorName allows deselecting the previous ego
	ActorName string
}

func (cmd ActorSelectEgo) Execute(app *App, done Promise) {
	actor, ok := app.actors[cmd.ActorName]
	if ok {
		app.ego = actor
	} else {
		app.ego = nil
	}

	done.Complete()
}
