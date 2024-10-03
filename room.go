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
	if r == nil {
		return nil
	}
	for _, actor := range r.actors {
		if !actor.IsEgo() && actor.Hotspot().Contains(pos) {
			return actor
		}
	}
	for _, obj := range r.objects {
		if obj.IsVisible() && obj.hotspot.Contains(pos) {
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
	Draw()
	Name() string
	Position() Position
	UsePosition() (Position, Direction)
}

// FindRoom returns the room with the given ID, or nil if not found.
func (a *App) FindRoom(id string) *Room {
	room, ok := a.rooms[id]
	if !ok {
		return nil
	}
	return room
}

// StartRoom starts the given room in the application.
func (a *App) StartRoom(room *Room) {
	if room == nil {
		log.Fatal("Room is nil")
	}
	for _, r := range a.rooms {
		if r == room {
			a.room = room
			return
		}
	}
	log.Fatalf("Room %s not declared", room.id)
}

func (a *App) drawSceneViewport() {
	if a.room != nil {
		a.room.Draw()
	}
}
