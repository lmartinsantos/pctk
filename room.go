package pctk

import (
	"io"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Room represents a room in the game.
type Room struct {
	bg      *Image
	objects map[string]*Object
}

// NewRoom creates a new room with the given background image.
func NewRoom(bg *Image) *Room {
	if bg.Width() < ScreenWidth || bg.Height() < ViewportHeight {
		log.Fatal("Background image is too small")
	}
	return &Room{
		bg:      bg,
		objects: make(map[string]*Object),
	}
}

// Objects returns a map of all objects present in the room.
func (r *Room) Objects() map[string]*Object {
	return r.objects
}

// ObjectByName retrieves an object from the room by its name.
func (r *Room) ObjectByName(name string) *Object {
	return r.objects[name]
}

// SaveObject adds an Object to the room's collection of objects.
func (r *Room) SaveObject(o *Object) {
	r.objects[o.name] = o
}

// BinaryEncode encodes the room data to a binary stream. The format is:
//   - the background image.
func (r *Room) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, r.bg)
}

// BinaryDecode decodes the room data from a binary stream. See Room.BinaryEncode for the format.
func (r *Room) BinaryDecode(rd io.Reader) error {
	r.bg = new(Image)
	// TODO make sense to define objects from the begining without ObjectShow commands
	r.objects = make(map[string]*Object)
	return BinaryDecode(rd, r.bg)
}

// RoomShow is a command that will show the room with the given resource.
type RoomShow struct {
	RoomRef ResourceRef
}

func (cmd RoomShow) Execute(app *App, done Promise) {
	// TODO: dispose the previous room if any
	app.room = app.res.LoadRoom(cmd.RoomRef)
	done.Complete()
}

func (a *App) drawBackgroud() {
	if a.room != nil {
		rl.DrawTexture(a.room.bg.Texture(), 0, 0, rl.White)
	}
}
