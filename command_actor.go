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
	TalkColor Color
	ScriptLoc FieldAccessor
	Size      Size
	UsePos    Position
	UseDir    Direction
}

func (cmd ActorDeclare) Execute(app *App, done *Promise) {
	actor := app.DeclareActor(cmd.ActorID, cmd.ActorName)
	if cmd.Costume != ResourceRefNull {
		actor.SetCostume(app.res.LoadCostume(cmd.Costume))
	}
	actor.Size = cmd.Size
	actor.TalkColor = cmd.TalkColor
	actor.UsePos = cmd.UsePos
	actor.UseDir = cmd.UseDir
	actor.scriptLoc = cmd.ScriptLoc
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
	Actor   *Actor
	Targets [2]RoomItem
	Verb    Verb
}

func (cmd ActorInteractWith) Execute(app *App, done *Promise) {
	if cmd.Verb == VerbWalkTo {
		// This is not an interaction, but a movement command. This should be unreachable.
		done.CompleteWithErrorf("invalid verb %s for actor interaction", cmd.Verb)
		return
	}

	var completed Future
	var args []any
	other := cmd.Targets[1]
	if other != nil {
		args = []any{other.ScriptLocation()}
	}
	switch item := cmd.Targets[0].(type) {
	case *Actor:
		completed = app.RunCommandSequence(
			ActorWalkToItem{
				Actor: cmd.Actor,
				Item:  cmd.Targets[0],
			},
			ActorCall{
				Actor:    item,
				Function: cmd.Verb.Action(),
				Args:     args,
			},
		)
	case *Object:
		if item.Owner() == cmd.Actor {
			switch cmd.Verb {
			case VerbWalkTo, VerbPickUp:
				// Verb not applicable to inventory item
				done.Complete()
				return
			}
			// The first argument is in the inventory. If there is a second argument that is not in
			// the inventory, walk to it and then interact. Othewise, just interact.
			if cmd.Targets[1] != nil && cmd.Targets[1].Owner() != cmd.Actor {
				completed = app.RunCommandSequence(
					ActorWalkToItem{
						Actor: cmd.Actor,
						Item:  cmd.Targets[1],
					},
					ObjectCall{
						Object:   item,
						Function: cmd.Verb.Action(),
						Args:     args,
					},
				)
			} else {
				completed = app.RunCommand(ObjectCall{
					Object:   item,
					Function: cmd.Verb.Action(),
					Args:     args,
				})
			}
		} else {
			// It is in the room.

			// Special case: use verb for a applicable object. Must walk to it, pick it up and then
			// do the rest of the action.
			if cmd.Verb == VerbUse && item.Class().IsOneOf(ObjectClassApplicable) && other != nil {
				completed = app.RunCommandSequence(
					ActorWalkToItem{
						Actor: cmd.Actor,
						Item:  item,
					},
					ObjectCall{
						Object:   item,
						Function: VerbPickUp.Action(),
						Args:     nil,
					},
					ActorWalkToItem{
						Actor: cmd.Actor,
						Item:  other,
					},
					ObjectCall{
						Object:   item,
						Function: cmd.Verb.Action(),
						Args:     args,
					},
				)
			} else {
				// General case. Walk to it and then interact.
				completed = app.RunCommandSequence(
					ActorWalkToItem{
						Actor: cmd.Actor,
						Item:  cmd.Targets[0],
					},
					ObjectCall{
						Object:   item,
						Function: cmd.Verb.Action(),
						Args:     args,
					},
				)
			}
		}
	default:
		log.Fatalf("unknown room item type %T", item)
	}
	completed = RecoverWithValue(completed, func(err error) any {
		if err != PromiseBroken {
			log.Printf("Actor interaction failed: %v", err)
		}
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
		cmd.Color = cmd.Actor.TalkColor
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

// ActorCall is a command that will call a method on an actor.
type ActorCall struct {
	Actor    *Actor
	Function string
	Args     []any
}

func (cmd ActorCall) Execute(app *App, done *Promise) {
	// We call this in the script of the current room. If the actor is imported from a include, it
	// will still have the methods for the actions no matter the room where it is evaluated. In
	// addition, this allows the script of the room to override the default behavior of the actor
	// when it is in the room.
	if app.room == nil {
		done.CompleteWithErrorf("no active room to call actor %s", cmd.Actor.Name())
		return
	}
	call := app.room.script.Call(
		cmd.Actor.ScriptLocation().Append(cmd.Function),
		cmd.Args,
		true,
	)
	call = Recover(call, func(err error) Future {
		return app.room.script.Call(
			WithDefaultsField(cmd.Function),
			append([]any{cmd.Actor.ScriptLocation()}, cmd.Args...),
			false,
		)
	})
	done.Bind(call)
}
