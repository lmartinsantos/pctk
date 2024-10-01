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
	act       *Action
	costume   *Costume
	elev      int
	ego       bool
	id        string
	inventory []*Object
	lookAt    Direction
	name      string
	pos       Positionf
	room      *Room
	size      Size
	speed     Positionf
}

// NewActor creates a new actor with the given ID and name.
func NewActor(id, name string) *Actor {
	return &Actor{
		id:    id,
		name:  name,
		pos:   DefaultActorPosition.ToPosf(),
		size:  DefaultActorSize,
		speed: DefaultActorSpeed,
	}
}

// AddToInventory adds an object to the actor's inventory.
func (a *Actor) AddToInventory(obj *Object) {
	a.inventory = append(a.inventory, obj)
	obj.owner = a
}

// CancelAction cancels the current action of the actor.
func (a *Actor) CancelAction() {
	if a.act != nil {
		a.act.Cancel()
	}
	a.act = nil
}

// Class returns the class of the actor.
func (a *Actor) Class() ObjectClass {
	return ObjectClassPerson
}

// Do executes the action in the actor.
func (a *Actor) Do(action *Action) Future {
	if a.act != nil {
		a.act.Cancel()
	}
	a.act = action
	return a.act.Done()
}

// Draw renders the actor in the viewport.
func (a *Actor) Draw() {
	if a.act == nil {
		a.act = Standing(a.lookAt)
	}

	if a.act.RunFrame(a) {
		a.act = nil
	}
}

// Hotspot returns the hotspot of the actor.
func (a *Actor) Hotspot() Rectangle {
	return Rectangle{Pos: a.costumePos(), Size: a.size}
}

// ID returns the ID of the actor.
func (a *Actor) ID() string {
	return a.name
}

// Inventory returns the inventory of the actor.
func (a *Actor) Inventory() []*Object {
	return a.inventory
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

// Action is an action that an actor is performing.
type Action struct {
	prom *Promise
	f    func(*Actor, *Promise)
}

// Standing creates a new action that makes an actor stand in the given direction.
func Standing(dir Direction) *Action {
	return &Action{
		prom: NewPromise(),
		f: func(a *Actor, done *Promise) {
			if cos := a.costume; cos != nil {
				cos.draw(CostumeIdle(dir), a.costumePos())
			}
		},
	}
}

// WalkingTo creates a new action that makes an actor walk to a given position.
func WalkingTo(pos Position) *Action {
	return &Action{
		prom: NewPromise(),
		f: func(a *Actor, done *Promise) {
			if cos := a.costume; cos != nil {
				cos.draw(CostumeWalk(a.lookAt), a.costumePos())
			}

			if a.pos.ToPos() == pos {
				done.Complete()
				return
			}

			a.lookAt = a.pos.ToPos().DirectionTo(pos)
			a.pos = a.pos.Move(pos.ToPosf(), a.speed.Scale(rl.GetFrameTime()))
		},
	}
}

// SpeakingTo creates a new action that makes an actor speak to a dialog.
func SpeakingTo(dialog Future) *Action {
	return &Action{
		prom: NewPromise(),
		f: func(a *Actor, done *Promise) {
			if cos := a.costume; cos != nil {
				cos.draw(CostumeSpeak(a.lookAt), a.costumePos())
			}
			if dialog.IsCompleted() {
				done.Complete()
			}
		},
	}
}

// Cancel cancels the action.
func (a *Action) Cancel() {
	a.prom.Break()
}

// Done returns a future that will be completed when the action is done.
func (a *Action) Done() Future {
	return a.prom
}

// RunFrame runs a frame of the action.
func (a *Action) RunFrame(actor *Actor) (completed bool) {
	a.f(actor, a.prom)
	return a.prom.IsCompleted()
}

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
	actor.Do(Standing(cmd.LookAt))
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
		done.CompleteWhen(a.Do(Standing(a.pos.ToPos().DirectionTo(cmd.Position))))
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
		a.Do(Standing(cmd.Direction))
	})
	done.Complete()
}

// ActorWalkToPosition is a command that will make an actor walk to a given position.
type ActorWalkToPosition struct {
	ActorID  string
	Position Position
}

func (cmd ActorWalkToPosition) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorID, func(actor *Actor) {
		done.CompleteWhen(actor.Do(WalkingTo(cmd.Position)))
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
		done.CompleteWhen(app.Do(ActorWalkToPosition{
			ActorID:  a.name,
			Position: pos,
		}).AndThen(func(_ any) Future {
			return app.Do(ActorStand{
				ActorID:   a.name,
				Direction: dir,
			})
		}))
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
		if obj.Owner() != nil {
			// Object in the inventory. Just call the script.
			done.CompleteWhen(a.room.script.call(app.room.id, "objects", obj.name, "lookat"))
			return
		}
		// Object in the room. First walk to it, then call the script when
		done.CompleteWhen(app.Do(ActorWalkToObject{
			ActorID:  a.name,
			ObjectID: obj.name,
		}).AndThen(func(_ any) Future {
			return a.room.script.call(app.room.id, "objects", obj.name, "lookat")
		}))
	})
}

// ActorPickUpObject is a command that will make an actor pick up an object.
type ActorPickUpObject struct {
	ActorID  string
	ObjectID string
}

func (cmd ActorPickUpObject) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorID, func(a *Actor) {
		obj := app.room.ObjectByID(cmd.ObjectID)
		if obj == nil || obj.owner != nil {
			done.Complete()
			return
		}
		done.CompleteWhen(app.Do(ActorWalkToObject{
			ActorID:  a.name,
			ObjectID: obj.name,
		}).AndThen(func(_ any) Future {
			return a.room.script.call(app.room.id, "objects", obj.name, "pickup")
		}))
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
		done.CompleteWhen(a.Do(SpeakingTo(dialogDone)))
	})
}

func (a *App) withActor(name string, f func(*Actor)) {
	actor, ok := a.actors[name]
	if !ok {
		return
	}
	f(actor)
}

func (a *App) ensureActor(id string) *Actor {
	actor, ok := a.actors[id]
	if !ok {
		actor = NewActor(id, id)
		a.actors[id] = actor
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

// ActorAddToInventory is a command that will add an object to an actor's inventory.
type ActorAddToInventory struct {
	ActorID  string
	ObjectID string
}

func (cmd ActorAddToInventory) Execute(app *App, done *Promise) {
	app.withActor(cmd.ActorID, func(actor *Actor) {
		obj := app.room.ObjectByID(cmd.ObjectID)
		if obj == nil {
			done.Complete()
			return
		}
		actor.AddToInventory(obj)
		done.Complete()
	})
}
