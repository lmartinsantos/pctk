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

// Locate the actor in the given room, position and direction.
func (a *Actor) Locate(room *Room, pos Position, dir Direction) {
	a.room = room
	a.pos = pos.ToPosf()
	a.Do(Standing(dir))
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
			a.lookAt = dir
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

// DeclareActor declares a new actor with the given ID and name.
func (a *App) DeclareActor(id, name string) *Actor {
	if _, ok := a.actors[id]; ok {
		log.Fatalf("Actor %s already exists", id)
	}
	actor := NewActor(id, name)
	a.actors[id] = actor
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
