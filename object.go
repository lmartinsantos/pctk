package pctk

import (
	"io"
)

var (
	DefaultObjectPosition = NewPos(160, 90)
)

// Object object refers to any interactive item or entity within the game's world.
type Object struct {
	name    string
	sprites *SpriteSheet
	pos     Position
	anim    *Animation // be more flexible being an anim instead of an sprite
	scripts map[VerbType]*Script
}

func NewObject(name string, sprites *SpriteSheet) *Object {
	return &Object{
		name:    name,
		sprites: sprites,
		scripts: make(map[VerbType]*Script),
	}
}

// WithAnimation sets an animation for the object.
func (o *Object) WithAnimation(anim *Animation) *Object {
	o.anim = anim
	return o
}

// WithScript assigns a script to a specific action for the object.
func (o *Object) WithScript(a VerbType, s *Script) *Object {
	o.scripts[a] = s
	return o
}

// FrameSize gets the object frame size
func (o *Object) FrameSize() Size {
	return o.sprites.frameSize
}

// BinaryEncode encodes the object to a binary format. The format is as follows:
// - uint32 name string length (in bytes).
// - name.
// - sprite sheet.
// - the animation.
func (o *Object) BinaryEncode(w io.Writer) (n int, err error) {
	//TODO
	return BinaryEncode(w, uint32(len(o.name)), []byte(o.name), o.sprites, o.anim)
}

// BinaryDecode decodes the object from a binary format. See BinaryEncode for the format.
func (o *Object) BinaryDecode(r io.Reader) error {
	// TODO
	o.sprites = new(SpriteSheet)
	o.anim = new(Animation)

	var length uint32
	if err := BinaryDecode(r, &length); err != nil {
		return err
	}
	nameBytes := make([]byte, length)
	if err := BinaryDecode(r, nameBytes); err != nil {
		return err
	}
	o.name = string(nameBytes)

	if err := BinaryDecode(r, o.sprites, o.anim); err != nil {
		return err
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
		o.anim.draw(o.sprites, o.pos)
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
		script := object.scripts[cmd.Verb.Type]
		if script != nil {
			script.run(app, done)
		}
	}

	done.Complete()
}
