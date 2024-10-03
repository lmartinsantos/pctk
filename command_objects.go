package pctk

// ObjectDeclare is a command that will declare a new object with the given properties.
type ObjectDeclare struct {
	Class     ObjectClass
	Hotspot   Rectangle
	Name      string
	ObjectID  string
	Pos       Position
	RoomID    string
	ScriptLoc FieldAccessor // The location of the object in the script
	Sprites   ResourceRef
	States    []*ObjectState
	UseDir    Direction
	UsePos    Position
}

func (cmd ObjectDeclare) Execute(app *App, done *Promise) {
	room := app.RoomByID(cmd.RoomID)
	var sprites *SpriteSheet
	if cmd.Sprites != ResourceRefNull {
		sprites = app.res.LoadSpriteSheet(cmd.Sprites)
	}

	obj := &Object{
		classes:   cmd.Class,
		hotspot:   cmd.Hotspot,
		id:        cmd.ObjectID,
		name:      cmd.Name,
		pos:       cmd.Pos,
		room:      room,
		scriptLoc: cmd.ScriptLoc,
		sprites:   sprites,
		states:    cmd.States,
		useDir:    cmd.UseDir,
		usePos:    cmd.UsePos,
	}
	room.DeclareObject(obj)
	done.Complete()
}

// ObjectCall is a command that will execute a script function of an object.
type ObjectCall struct {
	Object   *Object
	Function string
	Args     []any
}

func (cmd ObjectCall) Execute(app *App, done *Promise) {
	obj := cmd.Object
	call := obj.room.script.Call(
		cmd.Object.ScriptLocation().Append(cmd.Function),
		cmd.Args,
		true,
	)
	call = Recover(call, func(err error) Future {
		return obj.room.script.Call(
			WithDefaultsField(cmd.Function),
			append([]any{cmd.Object.ScriptLocation()}, cmd.Args...),
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
