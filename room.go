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

func (r *Room) BinaryEncode(w io.Writer) (int, error) {
	panic("not implemented")
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
		rl.DrawTexture(a.room.bg.tex, 0, 0, rl.White)
	}
}
