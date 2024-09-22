package pctk

import (
	"io"
)

// TODO
const (
	NoClass          = 0
	ClassUntouchable = 1 << iota
	ClassPickupable
	ClassTalkable
	ClassGiveable
	ClassOpenable
	ClassDoor
	// ClassActor
)

var (
	DefaultObjectPosition = NewPos(160, 90)
)

// Object object refers to any interactive item or entity within the game's world.
type Object struct {
	name    string
	sprites *SpriteSheet
	pos     Position
	states  []*State
	state   int // By default, objects are in state 0.
	classes uint
}

func NewObject(name string, sprites *SpriteSheet, classes uint) *Object {
	return &Object{
		name:    name,
		sprites: sprites,
		classes: classes,
		states:  []*State{},
	}
}

// WithState sets a new state for the object.
func (o *Object) WithState(newState *State) {
	o.states = append(o.states, newState)
}

// Rectangle returns the screen area (as a Rectangle) associated with the object's position
func (o *Object) Rectangle() Rectangle {
	size := o.sprites.frameSize
	return NewRect(o.pos.X, o.pos.Y, size.W, size.H)
}

func (o *Object) AddClass(newClass uint) {
	o.classes |= newClass
}

func (o *Object) RemoveClass(class uint) {
	o.classes &^= class
}

func (o *Object) HasClass(class uint) bool {
	return o.classes&class != 0
}

// BinaryEncode encodes the object to a binary format. The format is as follows:
// - uint32 name string length (in bytes).
// - name.
// - sprite sheet.
// - classes.
// - uint32: the number of states.
// - states.
func (o *Object) BinaryEncode(w io.Writer) (n int, err error) {
	n, err = BinaryEncode(w, uint32(len(o.name)), []byte(o.name), o.sprites, byte(o.classes), uint32(len(o.states)))
	if err != nil {
		return n, err
	}
	for _, state := range o.states {
		nn, err := BinaryEncode(w, state)
		n += nn
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// BinaryDecode decodes the object from a binary format. See BinaryEncode for the format.
func (o *Object) BinaryDecode(r io.Reader) error {
	o.sprites = new(SpriteSheet)
	o.states = make([]*State, 0)

	var length uint32
	if err := BinaryDecode(r, &length); err != nil {
		return err
	}
	nameBytes := make([]byte, length)
	if err := BinaryDecode(r, nameBytes); err != nil {
		return err
	}
	o.name = string(nameBytes)

	var count uint32
	var classes byte
	if err := BinaryDecode(r, o.sprites, &classes, &count); err != nil {
		return err
	}

	o.classes = uint(classes)

	for i := uint32(0); i < count; i++ {
		state := new(State)
		if err := BinaryDecode(r, state); err != nil {
			return err
		}
		o.states = append(o.states, state)
	}

	return nil
}

// State defines the various states of the object.
type State struct {
	anim    *Animation // be more flexible being an anim instead of an sprite
	scripts map[VerbType]*Script
}

func NewState() *State {
	return &State{
		scripts: make(map[VerbType]*Script),
	}
}

func (s *State) WithAnimation(anim *Animation) *State {
	s.anim = anim
	return s
}

func (s *State) WithScript(v VerbType, script *Script) *State {
	s.scripts[v] = script
	return s
}

// BinaryEncode encodes the object's state to a binary format. The format is as follows:
// - anim.
// - uint32: the number of scripts.
// - for each script:
//   - byte: the verb.
//   - script.
func (s *State) BinaryEncode(w io.Writer) (n int, err error) {
	n, err = BinaryEncode(w, s.anim, uint32(len(s.scripts)))
	if err != nil {
		return n, err
	}

	for verb, script := range s.scripts {
		nn, err := BinaryEncode(w, byte(verb), script)
		n += nn
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// BinaryDecode decodes the object's state from a binary format. See BinaryEncode for the format.
func (s *State) BinaryDecode(r io.Reader) error {
	s.anim = new(Animation)
	s.scripts = make(map[VerbType]*Script)

	if err := BinaryDecode(r, s.anim); err != nil {
		return err
	}
	var count uint32
	if err := BinaryDecode(r, &count); err != nil {
		return err
	}

	for i := uint32(0); i < count; i++ {
		script := new(Script)
		var verb byte
		if err := BinaryDecode(r, &verb, script); err != nil {
			return err
		}

		s.scripts[VerbType(verb)] = script

	}

	return nil
}

// ObjectShow is a command that will show an object in the room at the given position.
type ObjectShow struct {
	ObjectResource ResourceRef
	ObjectName     string
	Position       Position
}

func (cmd ObjectShow) Execute(app *App, done Promise) {
	object := app.res.LoadObject(cmd.ObjectResource)
	object.pos = cmd.Position
	app.objects[cmd.ObjectName] = object
	done.Complete()
}

func (a *App) drawObjects() {
	for _, o := range a.objects {
		state := o.states[o.state]
		// TODO
		state.anim.draw(o.sprites, o.pos)
	}
}

// ObjectRelease is a command that will release an object removing it from the application.
type ObjectRelease struct {
	ObjectName string
}

func (cmd ObjectRelease) Execute(app *App, done Promise) {
	delete(app.objects, cmd.ObjectName)
	done.Complete()
}

// TODO object source & object target
// ObjectOnAction is a command that will run the action script related to an object.
type ObjectOnAction struct {
	ObjectName string
	Verb       *Verb
}

func (cmd ObjectOnAction) Execute(app *App, done Promise) {
	object := app.objects[cmd.ObjectName]
	if object != nil {
		state := object.states[object.state]
		script := state.scripts[cmd.Verb.Type]
		if script != nil {
			script.run(app, done)
		}
	}

	done.Complete()
}
