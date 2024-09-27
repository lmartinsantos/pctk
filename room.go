package pctk

import (
	"io"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Room represents a room in the game.
type Room struct {
	background *Image    // The background image of the room
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

// BinaryEncode encodes the room data to a binary stream. The format is:
//   - the background image.
func (r *Room) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, r.background)
}

// BinaryDecode decodes the room data from a binary stream. See Room.BinaryEncode for the format.
func (r *Room) BinaryDecode(rd io.Reader) error {
	r.background = new(Image)
	return BinaryDecode(rd, r.background)
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
	app.room.script.call(cmd.RoomID, "enter", done)
}

func (a *App) drawBackgroud() {
	if a.room != nil {
		rl.DrawTexture(a.room.background.Texture(), 0, 0, rl.White)
	}
}
