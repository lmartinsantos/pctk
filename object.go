package pctk

var (
	DefaultObjectPosition = NewPos(160, 90)
)

// Object object refers to any interactive item or entity within the game's world.
type Object struct {
	name    string
	pos     Position
	anim    *Animation // be more flexible being an anim instead of an sprite
	scripts map[ActionName]*Script
}

func NewObject(name string) *Object {
	return &Object{
		name:    name,
		scripts: make(map[ActionName]*Script),
	}
}

// WithAnimation sets an animation for the object.
func (o *Object) WithAnimation(anim *Animation) *Object {
	o.anim = anim
	return o
}

// WithScript assigns a script to a specific action for the object.
func (o *Object) WithScript(a ActionName, s *Script) *Object {
	o.scripts[a] = s
	return o
}

// ObjectShow is a command that will show an object in the scene at the given position.
type ObjectShow struct {
	ObjectResource ResourceLocator
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
		o.anim.draw(a, o.pos)
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
	Action     *Action
}

func (cmd ObjectOnAction) Execute(app *App, done Promise) {
	object := app.objects[cmd.ObjectName]
	if object != nil {
		script := object.scripts[cmd.Action.ActionName]
		if script != nil {
			script.run(app, done)
		}
	}

	done.Complete()
}
