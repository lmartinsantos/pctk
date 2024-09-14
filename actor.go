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
	standU *Animation
	standD *Animation
	speakH *Animation
	speakU *Animation
	speakD *Animation
	walkH  *Animation
	walkU  *Animation
	walkD  *Animation

	lookAt Direction
	pos    Position
	act    action
}

func NewActor(name string) *Actor {
	return &Actor{name: name}
}

func (a *Actor) WithStandHorizontal(anim *Animation) *Actor {
	a.standH = anim
	return a
}

func (a *Actor) WithStandUp(anim *Animation) *Actor {
	a.standU = anim
	return a
}

func (a *Actor) WithStandDown(anim *Animation) *Actor {
	a.standD = anim
	return a
}

func (a *Actor) WithSpeakHorizontal(anim *Animation) *Actor {
	a.speakH = anim
	return a
}

func (a *Actor) WithSpeakUp(anim *Animation) *Actor {
	a.speakU = anim
	return a
}

func (a *Actor) WithSpeakDown(anim *Animation) *Actor {
	a.speakD = anim
	return a
}

func (a *Actor) WithWalkHorizontal(anim *Animation) *Actor {
	a.walkH = anim
	return a
}

func (a *Actor) WithWalkUp(anim *Animation) *Actor {
	a.walkU = anim
	return a
}

func (a *Actor) WithWalkDown(anim *Animation) *Actor {
	a.walkD = anim
	return a
}

func (a *Actor) stand(dir Direction) *Actor {
	a.lookAt = dir
	a.act = func(app *App, c *Actor) (completed bool) {
		switch dir {
		case DirRight:
			if a.standH != nil {
				a.standH.draw(app, c.pos, false)
			}
		case DirLeft:
			if a.standH != nil {
				a.standH.draw(app, c.pos, true)
			}
		case DirUp:
			if a.standU != nil {
				a.standU.draw(app, c.pos, false)
			}
		case DirDown:
			if a.standD != nil {
				a.standD.draw(app, c.pos, false)
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
				if a.walkH != nil {
					a.walkH.draw(app, c.pos, false)
				}
			case DirLeft:
				c.pos.X--
				if a.walkH != nil {
					a.walkH.draw(app, c.pos, true)
				}
			case DirUp:
				c.pos.Y--
				if a.walkU != nil {
					a.walkU.draw(app, c.pos, false)
				}
			case DirDown:
				c.pos.Y++
				if a.walkD != nil {
					a.walkD.draw(app, c.pos, false)
				}
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
			case DirUp:
				if a.speakU != nil {
					a.speakU.draw(app, c.pos, false)
				}
			case DirDown:
				if a.speakD != nil {
					a.speakD.draw(app, c.pos, false)
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
