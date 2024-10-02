package pctk

import (
	"log"
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

// Room returns the room where the actor is.
func (a *Actor) Room() *Room {
	return a.room
}

// SetCostume sets the costume for the actor.
func (a *Actor) SetCostume(costume *Costume) *Actor {
	a.costume = costume
	return a
}

// UsePosition returns the position where actors interact with the actor.
func (a *Actor) UsePosition() (Position, Direction) {
	// TODO: this might be wrong, specially if the actor is looking to the edge of a walking box
	return a.pos.ToPos(), a.lookAt
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
	Actor    *Actor
	Position Position
}

func (cmd ActorLookAtPos) Execute(app *App, done *Promise) {
	done.Bind(cmd.Actor.Do(Standing(cmd.Actor.pos.ToPos().DirectionTo(cmd.Position))))
}

// ActorStand is a command that will make an actor stand in the given direction.
type ActorStand struct {
	Actor     *Actor
	Direction Direction
}

func (cmd ActorStand) Execute(app *App, done *Promise) {
	cmd.Actor.Do(Standing(cmd.Direction))
	done.Complete()
}

// ActorWalkToPosition is a command that will make an actor walk to a given position.
type ActorWalkToPosition struct {
	Actor    *Actor
	Position Position
}

func (cmd ActorWalkToPosition) Execute(app *App, done *Promise) {
	if cmd.Actor.Room() != app.room {
		done.CompleteWithErrorf("actor %s is not in the room", cmd.Actor.Name())
		return
	}
	done.Bind(cmd.Actor.Do(WalkingTo(cmd.Position)))
}

// ActorWalkToItem is a command that will make an actor walk to a room item.
type ActorWalkToItem struct {
	Actor *Actor
	Item  RoomItem
}

func (cmd ActorWalkToItem) Execute(app *App, done *Promise) {
	switch item := cmd.Item.(type) {
	case *Actor:
		if item.Room() != app.room {
			done.CompleteWithErrorf("actor %s is not in the room", item.Name())
			return
		}
	case *Object:
		if item.Owner() != nil {
			done.CompleteWithErrorf("object %s is in the inventory", item.Name())
		}
	}
	pos, dir := cmd.Item.UsePosition()

	done.Bind(app.RunCommandSequence(
		ActorWalkToPosition{
			Actor:    cmd.Actor,
			Position: pos,
		},
		ActorStand{
			Actor:     cmd.Actor,
			Direction: dir,
		},
	))

}

// ActorInteractWith is a command that will make an actor interact with an object.
type ActorInteractWith struct {
	Actor  *Actor
	Target RoomItem
	Verb   Verb
}

func (cmd ActorInteractWith) Execute(app *App, done *Promise) {
	var completed Future
	switch item := cmd.Target.(type) {
	case *Actor:
		completed = app.RunCommandSequence(
			ActorWalkToItem{
				Actor: cmd.Actor,
				Item:  cmd.Target,
			},
			// TODO: implement ActorCall
		)
	case *Object:
		if item.Owner() != nil {
			switch cmd.Verb {
			case VerbWalkTo, VerbPickUp:
				// Verb not applicable to inventory item
				done.Complete()
				return
			default:
				// It is in the inventory. Do not walk to it, just call.
				completed = app.RunCommand(ObjectCall{
					Object: item,
					Action: cmd.Verb.Action(),
				})
			}
		} else {
			// It is in the room. Walk to it and then interact.
			completed = app.RunCommandSequence(
				ActorWalkToItem{
					Actor: cmd.Actor,
					Item:  cmd.Target,
				},
				ObjectCall{
					Object: item,
					Action: cmd.Verb.Action(),
				},
			)
		}
	default:
		log.Fatalf("unknown room item type %T", item)
	}
	completed = RecoverWithValue(completed, func(err error) any {
		log.Printf("Actor interaction failed: %v", err)
		return nil
	})
	done.Bind(completed)
}

// ActorSpeak is a command that will make an actor speak the given text.
type ActorSpeak struct {
	Actor *Actor
	Text  string
	Delay time.Duration
	Color Color
}

func (cmd ActorSpeak) Execute(app *App, done *Promise) {
	if cmd.Delay == 0 {
		cmd.Delay = DefaultActorSpeakDelay
	}

	if cmd.Color == rl.Blank {
		cmd.Color = rl.White
	}

	dialogDone := app.RunCommand(ShowDialog{
		Text:     cmd.Text,
		Position: cmd.Actor.dialogPos(),
		Color:    cmd.Color,
		Speed:    1.0,
	})
	done.Bind(cmd.Actor.Do(SpeakingTo(dialogDone)))
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

// SelectEgo sets actor as the ego.
func (a *App) SelectEgo(actor *Actor) {
	if a.ego != nil {
		a.ego.ego = false
	}
	a.ego = actor
	if a.ego != nil {
		a.ego.ego = true
	}
}

// ActorSelectEgo is a command that will make an actor be the actor under player's control.
type ActorSelectEgo struct {
	Actor *Actor
}

func (cmd ActorSelectEgo) Execute(app *App, done *Promise) {
	app.SelectEgo(cmd.Actor)
	done.CompleteWithValue(cmd.Actor)
}

// ActorAddToInventory is a command that will add an object to an actor's inventory.
type ActorAddToInventory struct {
	Actor  *Actor
	Object *Object
}

func (cmd ActorAddToInventory) Execute(app *App, done *Promise) {
	cmd.Actor.AddToInventory(cmd.Object)
	done.CompleteWithValue(cmd)
}

// ActorByID returns the actor with the given ID, or nil if not found.
func (a *App) ActorByID(id string) *Actor {
	actor, _ := a.actors[id]
	return actor
}
