package pctk

// Object represents an object in the game. Objects are defined in the scope of rooms and generated
// by the room scripts.
type Object struct {
	classes ObjectClass    // The classes the object belongs to as OR-ed bit flags
	hotspot Rectangle      // The hotspot of the object (for mouse interaction)
	id      string         // The ID of the object
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

// Class returns the class of the object.
func (o *Object) Class() ObjectClass {
	return o.classes
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
	if !o.IsVisible() {
		return
	}
	if st := o.CurrentState(); st != nil {
		st.Anim.Draw(o.sprites, o.pos.Sub(NewPos(o.sprites.frameSize.W/2, o.sprites.frameSize.H)))
	}
}

// ID returns the ID of the object.
func (o *Object) ID() string {
	return o.id
}

// IsVisible returns true if the object is visible in the room, false otherwise.
func (o *Object) IsVisible() bool {
	return o.owner == nil
}

// Name returns the name of the object.
func (o *Object) Name() string {
	return o.name
}

// Owner returns the actor that owns the object, or nil if not picked up.
func (o *Object) Owner() *Actor {
	if o == nil {
		return nil
	}
	return o.owner
}

// Position returns the position of the object.
func (o *Object) Position() Position {
	return o.pos
}

// UsePosition returns the position where actors interact with the object.
func (o *Object) UsePosition() (Position, Direction) {
	return o.usePos, o.useDir
}

// ObjectState represents a state of an object.
type ObjectState struct {
	Anim *Animation // The animation while in this state.
}

// ObjectClass represents a class of objects. Classes are aimed to be used as bit flags that can be
// OR-ed together. As this type is backed by a uint64, there can be up to 64 different classes.
// There are two kind of classes: the built-in classes and the custom classes.
type ObjectClass uint64

const (
	// ObjectClassPerson is a built-in class that represents objects that are persons.
	ObjectClassPerson ObjectClass = 1 << 0

	// ObjectClassPickable is a built-in class that represents objects that can be picked up by the
	// player.
	ObjectClassPickable ObjectClass = 1 << 1

	// ObjectClassOpenable is a built-in class that represents objects that can be opened by the
	// player.
	ObjectClassOpenable ObjectClass = 1 << 2
)

// Is returns true if the object class is the given class, false otherwise.
func (c ObjectClass) Is(other ObjectClass) bool {
	return c&other != 0
}
