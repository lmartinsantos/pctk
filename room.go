package pctk

import (
	"log"
	"slices"
)

// Room represents a room in the game.
type Room struct {
	actors     []*Actor  // The actors in the room
	background *Image    // The background image of the room
	id         string    // The ID of the room
	objects    []*Object // The objects declared in the room
	script     *Script   // The script where this room is defined. Used to call the room functions.
}

// NewRoom creates a new room with the given background image.
func NewRoom(bg *Image) *Room {
	if bg.Width() < ScreenWidth || bg.Height() < ViewportHeight {
		log.Fatal("Background image is too small")
	}
	return &Room{
		background: bg,
	}
}

// RoomByID returns the room with the given ID, panicking if not found.
func (a *App) RoomByID(id string) *Room {
	room, ok := a.rooms[id]
	if !ok {
		log.Fatalf("Room %s not found", id)
	}
	return room
}

// DeclareObject declares an object in the room.
func (r *Room) DeclareObject(obj *Object) {
	obj.room = r
	r.objects = append(r.objects, obj)
}

// Draw renders the room in the viewport.
func (r *Room) Draw() {
	r.background.Draw(NewPos(0, 0), White)
	items := make([]RoomItem, 0, len(r.actors)+len(r.objects))
	for _, actor := range r.actors {
		items = append(items, actor)
	}
	for _, obj := range r.objects {
		items = append(items, obj)
	}
	slices.SortFunc(items, func(a, b RoomItem) int {
		return a.Position().Y - b.Position().Y
	})
	for _, item := range items {
		item.Draw()
	}
}

// ItemAt returns the item at the given position in the room.
func (r *Room) ItemAt(pos Position) RoomItem {
	for _, actor := range r.actors {
		if !actor.IsEgo() && actor.Hotspot().Contains(pos) {
			return actor
		}
	}
	for _, obj := range r.objects {
		if obj.hotspot.Contains(pos) {
			return obj
		}
	}
	return nil
}

// ObjectByID returns the object with the given ID, or nil if not found.
func (r *Room) ObjectByID(id string) *Object {
	for _, obj := range r.objects {
		if obj.name == id {
			return obj
		}
	}
	return nil
}

// PutActor puts an actor in the room.
func (r *Room) PutActor(actor *Actor) {
	actor.room = r
	for _, act := range r.actors {
		if act == actor {
			return
		}
	}
	r.actors = append(r.actors, actor)
}

// RoomItem is an item from a room that can be represented in the viewport.
type RoomItem interface {
	Class() ObjectClass
	Name() string
	Position() Position
	Draw()
}

// RoomDeclare is a command that will declare a new room with the given properties.
type RoomDeclare struct {
	BackgroundRef ResourceRef
	RoomID        string
	Script        *Script
}

func (cmd RoomDeclare) Execute(app *App, done *Promise) {
	if _, ok := app.rooms[cmd.RoomID]; ok {
		log.Fatalf("Room %s already exists", cmd.RoomID)
	}
	room := Room{
		id:         cmd.RoomID,
		background: app.res.LoadImage(cmd.BackgroundRef),
		script:     cmd.Script,
	}
	app.rooms[cmd.RoomID] = &room
	done.CompleteWithValue(room)
}

// RoomShow is a command that will show the room with the given resource.
type RoomShow struct {
	RoomID string
}

func (cmd RoomShow) Execute(app *App, done *Promise) {
	// TODO: execute exit function and dispose the previous room if any.
	var ok bool
	app.room, ok = app.rooms[cmd.RoomID]
	if !ok {
		log.Fatalf("Room %s not found", cmd.RoomID)
	}

	// Call the enter function of the room script.
	done.CompleteWhen(app.room.script.call(cmd.RoomID, "enter"))
}

func (a *App) drawViewport() {
	if a.room != nil {
		a.room.Draw()
	}
}
