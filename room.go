package pctk

import (
	"io"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Room represents a room in the game.
type Room struct {
	bg *Image
}

// NewRoom creates a new room with the given background image.
func NewRoom(bg *Image) *Room {
	if bg.Width() < ScreenWidth || bg.Height() < ViewportHeight {
		log.Fatal("Background image is too small")
	}
	return &Room{
		bg: bg,
	}
}

// BinaryEncode encodes the room data to a binary stream. The format is:
//   - the background image.
func (r *Room) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, r.bg)
}

// BinaryDecode decodes the room data from a binary stream. See Room.BinaryEncode for the format.
func (r *Room) BinaryDecode(rd io.Reader) error {
	r.bg = new(Image)
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
