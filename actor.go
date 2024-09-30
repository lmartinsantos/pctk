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
	DefaultActorSpeed     = NewPosf(80, 20)
	DefaultActorSize      = NewSize(32, 48)
	DefaultActorDirection = DirRight
)

type Actor struct {
	act     action
	costume *Costume
	elev    int
	lookAt  Direction
	name    string
	pos     Positionf
	room    *Room
	size    Size
	speed   Positionf
}

func NewActor(name string) *Actor {
	return &Actor{
		name:  name,
		pos:   DefaultActorPosition.ToPosf(),
		size:  DefaultActorSize,
		speed: DefaultActorSpeed,
	}
}

// Draw renders the actor in the viewport.
func (a *Actor) Draw() {
	if a.act == nil {
		a.stand(a.lookAt)
	}
	if a.act() {
		a.stand(a.lookAt)
	}
}

// Position returns the position of the actor.
func (a *Actor) Position() Position {
	return a.pos.ToPos()
}

// SetCostume sets the costume for the actor.
func (a *Actor) SetCostume(costume *Costume) *Actor {
	a.costume = costume
	return a
}

func (a *Actor) costumePos() Position {
	return a.pos.ToPos().Sub(NewPos(a.size.W/2, a.size.H-a.elev))
}

func (a *Actor) dialogPos() Position {
	return a.pos.ToPos().Above(a.size.H + 40)
}

func (a *Actor) stand(dir Direction) *Actor {
	a.lookAt = dir
	a.act = func() (completed bool) {
		if cos := a.costume; cos != nil {
			cos.draw(CostumeIdle(dir), a.costumePos())
		}
		return false
	}
	return a
}

type action func() (completed bool)

// ActorShow is a command that will show an actor in the room at the given position.
type ActorShow struct {
	CostumeResource ResourceRef
	ActorID         string
	Position        Position
	LookAt          Direction
}

func (cmd ActorShow) Execute(app *App, done *Promise) {
	actor := app.ensureActor(cmd.ActorID)
	app.room.PutActor(actor)
	actor.pos = cmd.Position.ToPosf()
	actor.stand(cmd.LookAt)
	if cmd.CostumeResource != ResourceRefNull {
		actor.SetCostume(app.res.LoadCostume(cmd.CostumeResource))
	}
	app.actors[cmd.ActorID] = actor
	done.Complete()
}

// ActorLookAtPos is a command that will make an actor look at a given position.
type ActorLookAtPos struct {
	ActorName string
	Position  Position
}

func (cmd ActorLookAtPos) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorName, func(a *Actor) {
		a.stand(a.pos.ToPos().DirectionTo(cmd.Position))
	})
	done.Complete()
}

// ActorStand is a command that will make an actor stand in the given direction.
type ActorStand struct {
	ActorID   string
	Direction Direction
}

func (cmd ActorStand) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorID, func(a *Actor) {
		a.stand(cmd.Direction)
	})
	done.Complete()
}

// ActorWalkToPosition is a command that will make an actor walk to a given position.
type ActorWalkToPosition struct {
	ActorID  string
	Position Position
}

func (cmd ActorWalkToPosition) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorID, func(a *Actor) {
		a.act = func() (completed bool) {
			if cos := a.costume; cos != nil {
				cos.draw(CostumeWalk(a.lookAt), a.costumePos())
			}

			if a.pos.ToPos() == cmd.Position {
				done.Complete()
				return true
			}

			a.lookAt = a.pos.ToPos().DirectionTo(cmd.Position)
			a.pos = a.pos.Move(cmd.Position.ToPosf(), a.speed.Scale(rl.GetFrameTime()))
			return false
		}
	})
}

// ActorSpeak is a command that will make an actor speak the given text.
type ActorSpeak struct {
	ActorID string
	Text    string
	Delay   time.Duration
	Color   Color
}

func (cmd ActorSpeak) Execute(app *App, done *Promise) {
	if cmd.Delay == 0 {
		cmd.Delay = DefaultActorSpeakDelay
	}

	if cmd.Color == rl.Blank {
		cmd.Color = rl.White
	}

	app.withActor(cmd.ActorID, func(a *Actor) {
		dialogDone := app.doNow(ShowDialog{
			Text:     cmd.Text,
			Position: a.dialogPos(),
			Color:    cmd.Color,
			Speed:    1.0,
		})
		a.act = func() (completed bool) {
			if cos := a.costume; cos != nil {
				cos.draw(CostumeSpeak(a.lookAt), a.costumePos())
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

// ActorSelectEgo is a command that will make an actor be the actor under player's control.
type ActorSelectEgo struct {
	// Using an empty ActorID allows deselecting the previous ego
	ActorID string
}

func (cmd ActorSelectEgo) Execute(app *App, done *Promise) {
	actor, ok := app.actors[cmd.ActorID]
	if ok {
		app.ego = actor
	} else {
		app.ego = nil
	}

	done.Complete()
}
