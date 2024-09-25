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

	costume *Costume

	lookAt    Direction
	pos       Position
	act       action
	inventory *Inventory
}

func NewActor(name string) *Actor {
	return &Actor{
		name:      name,
		inventory: NewInventory(),
	}
}

// Bounds is implemented to satisfy the Interactable interface.
func (a *Actor) Bounds() Rectangle {
	size := a.costume.sprites.frameSize
	return NewRect(a.pos.X, a.pos.Y, size.W, size.H)
}

// Description is implemented to satisfy the Interactable interface.
func (a *Actor) Description() string {
	return a.name
}

// SetCostume sets the costume for the actor.
func (a *Actor) SetCostume(costume *Costume) *Actor {
	a.costume = costume
	return a
}

func (a *Actor) stand(dir Direction) *Actor {
	a.lookAt = dir
	a.act = func(app *App, c *Actor) (completed bool) {
		if cos := a.costume; cos != nil {
			cos.draw(CostumeIdle(dir), c.pos)
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

// ActorShow is a command that will show an actor in the room at the given position.
type ActorShow struct {
	CostumeResource ResourceRef
	ActorName       string
	Position        Position
	LookAt          Direction
}

func (cmd ActorShow) Execute(app *App, done Promise) {
	actor := app.ensureActor(cmd.ActorName)
	actor.pos = cmd.Position
	actor.stand(cmd.LookAt)
	if cmd.CostumeResource != ResourceRefNull {
		actor.SetCostume(app.res.LoadCostume(cmd.CostumeResource))
	}
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
			if cos := a.costume; cos != nil {
				cos.draw(CostumeWalk(a.lookAt), c.pos)
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
			if cos := a.costume; cos != nil {
				cos.draw(CostumeSpeak(a.lookAt), c.pos)
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

func (a *App) ensureActor(name string) *Actor {
	actor, ok := a.actors[name]
	if !ok {
		actor = NewActor(name)
		a.actors[name] = actor
	}
	return actor
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
	ego := app.ego
	if ok {
		ego.actor = actor
	} else {
		ego.Clear()
	}

	done.Complete()
}
