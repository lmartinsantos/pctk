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
	ego     bool
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

// Hotspot returns the hotspot of the actor.
func (a *Actor) Hotspot() Rectangle {
	return Rectangle{Pos: a.costumePos(), Size: a.size}
}

// IsEgo returns true if the actor is the actor under player's control, false otherwise.
func (a *Actor) IsEgo() bool {
	return a.ego
}

// Name returns the name of the actor.
func (a *Actor) Name() string {
	return a.name
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

// ActorWalkToObject is a command that will make an actor walk to an object.
type ActorWalkToObject struct {
	ActorID  string
	ObjectID string
}

func (cmd ActorWalkToObject) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorID, func(a *Actor) {
		obj := app.room.ObjectByID(cmd.ObjectID)
		if obj == nil {
			done.Complete()
			return
		}
		pos, dir := obj.UsePos()
		done.CompleteWhen(Sequence(
			func() Future {
				return app.Do(ActorWalkToPosition{
					ActorID:  a.name,
					Position: pos,
				})
			},
			func() Future {
				return app.Do(ActorStand{
					ActorID:   a.name,
					Direction: dir,
				})
			},
		))
	})
}

// ActorLookAtObject is a command that will make an actor look at an object.
type ActorLookAtObject struct {
	ActorID  string
	ObjectID string
}

func (cmd ActorLookAtObject) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorID, func(a *Actor) {
		obj := app.room.ObjectByID(cmd.ObjectID)
		if obj == nil {
			done.Complete()
			return
		}
		done.CompleteWhen(Sequence(
			func() Future {
				return app.Do(ActorWalkToObject{
					ActorID:  a.name,
					ObjectID: obj.name,
				})
			},
			func() Future {
				return a.room.script.call(app.room.id, "objects", obj.name, "lookat")
			},
		))
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
	if app.ego != nil {
		app.ego.ego = false
	}
	app.ego = app.ensureActor(cmd.ActorID)
	app.ego.ego = true

	done.Complete()
}
