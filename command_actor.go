package pctk

import (
	"log"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ActorDeclare is a command that will declare an actor.
type ActorDeclare struct {
	ActorID   string
	ActorName string
	Costume   ResourceRef
}

func (cmd ActorDeclare) Execute(app *App, done *Promise) {
	actor := app.DeclareActor(cmd.ActorID, cmd.ActorName)
	if cmd.Costume != ResourceRefNull {
		actor.SetCostume(app.res.LoadCostume(cmd.Costume))
	}
	done.CompleteWithValue(cmd)
}

// ActorShow is a command that will show an actor in the room at the given position.
type ActorShow struct {
	Actor    *Actor
	Position Position
	LookAt   Direction
}

func (cmd ActorShow) Execute(app *App, done *Promise) {
	if app.room == nil {
		done.CompleteWithErrorf("no active room to show actor %s", cmd.Actor.Name())
		return
	}

	app.room.PutActor(cmd.Actor)
	cmd.Actor.Locate(app.room, cmd.Position, cmd.LookAt)
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
