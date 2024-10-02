package pctk

// ObjectDeclare is a command that will declare a new object with the given properties.
type ObjectDeclare struct {
	Classes  ObjectClass
	Hotspot  Rectangle
	Name     string
	ObjectID string
	Pos      Position
	RoomID   string
	Sprites  ResourceRef
	States   []*ObjectState
	UseDir   Direction
	UsePos   Position
}

func (cmd ObjectDeclare) Execute(app *App, done *Promise) {
	room := app.RoomByID(cmd.RoomID)
	var sprites *SpriteSheet
	if cmd.Sprites != ResourceRefNull {
		sprites = app.res.LoadSpriteSheet(cmd.Sprites)
	}

	obj := &Object{
		classes: cmd.Classes,
		hotspot: cmd.Hotspot,
		id:      cmd.ObjectID,
		name:    cmd.Name,
		pos:     cmd.Pos,
		room:    room,
		sprites: sprites,
		states:  cmd.States,
		useDir:  cmd.UseDir,
		usePos:  cmd.UsePos,
	}
	room.DeclareObject(obj)
	done.Complete()
}

// ObjectCall is a command that will execute a script function of an object.
type ObjectCall struct {
	Object   *Object
	Function string
}

func (cmd ObjectCall) Execute(app *App, done *Promise) {
	obj := cmd.Object
	call := obj.room.script.Call(
		WithField(obj.room.id, "objects", obj.id, cmd.Function),
		nil,
		true,
	)
	call = Recover(call, func(err error) Future {
		return obj.room.script.Call(
			WithField("default", cmd.Function),
			[]any{WithField(obj.room.id, "objects", obj.id)},
			false,
		)
	})
	done.Bind(call)
}

// FindObject returns the object with the given ID in the room, or nil if not found.
func (a *App) FindObject(roomID, objectID string) *Object {
	room := a.RoomByID(roomID)
	if room == nil {
		return nil
	}
	return room.ObjectByID(objectID)
}
