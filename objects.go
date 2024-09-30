package pctk

// Object represents an object in the game. Objects are defined in the scope of rooms and generated
// by the room scripts.
type Object struct {
	classes ObjectClass    // The classes the object belongs to as OR-ed bit flags
	hotspot Rectangle      // The hotspot of the object (for mouse interaction)
	name    string         // The name of the object as seen by the player
	owner   *Actor         // The actor that owns the object, or nil if not picked up
	pos     Position       // The position of the object in its room (for rendering)
	room    *Room          // The room where the object is declared, and where actions code resides
	sprites *SpriteSheet   // The sprites of the object
	states  []*ObjectState // The states the object can be in
	state   int            // The current state of the object
	useDir  Direction      // The direction the actor when using the object
	usePos  Position       // The position the actor was when using the object
}

// CurrentState returns the current state of the object.
func (o *Object) CurrentState() *ObjectState {
	if o.state < 0 || o.state >= len(o.states) {
		return nil
	}
	return o.states[o.state]
}

// Draw renders the object in the viewport.
func (o *Object) Draw() {
	if st := o.CurrentState(); st != nil {
		st.Anim.Draw(o.sprites, o.pos.Sub(NewPos(o.sprites.frameSize.W/2, o.sprites.frameSize.H)))
	}
}

// IsVisible returns true if the object is visible in the room, false otherwise.
func (o *Object) IsVisible() bool {
	return o.owner == nil
}

// Position returns the position of the object.
func (o *Object) Position() Position {
	return o.pos
}

// ObjectState represents a state of an object.
type ObjectState struct {
	Anim *Animation // The animation while in this state.
}

// ObjectClass represents a class of objects. Classes are aimed to be used as bit flags that can be
// OR-ed together. As this type is backed by a uint64, there can be up to 64 different classes.
// Classes are defined in the game scripts, so their meaning is up to the game designer.
type ObjectClass uint64

// ObjectDeclare is a command that will declare a new object with the given properties.
type ObjectDeclare struct {
	Classes ObjectClass
	Hotspot Rectangle
	Name    string
	Pos     Position
	RoomID  string
	Sprites ResourceRef
	States  []*ObjectState
	UseDir  Direction
	UsePos  Position
}

func (cmd ObjectDeclare) Execute(app *App, done *Promise) {
	room := app.RoomByID(cmd.RoomID)
	sprites := app.res.LoadSpriteSheet(cmd.Sprites)
	obj := &Object{
		classes: cmd.Classes,
		hotspot: cmd.Hotspot,
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
