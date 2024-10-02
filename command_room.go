package pctk

import "log"

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
	Room *Room
}

func (cmd RoomShow) Execute(app *App, done *Promise) {
	var job Future

	if app.room != nil {
		job = IgnoreError(app.room.script.Call(WithMethod(app.room.id, "exit")), nil)
	}

	// Call the enter function of the room script.
	app.room = cmd.Room
	job = Continue(job, func(a any) Future {
		return IgnoreError(cmd.Room.script.Call(WithMethod(cmd.Room.id, "enter")), nil)
	})

	done.Bind(job)
}
