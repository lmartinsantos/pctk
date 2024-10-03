package pctk

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/go-lua"
)

func (s *Script) luaInit(app *App) {
	if s.l == nil {
		s.l = lua.NewState()
		lua.BaseOpen(s.l)

		luaDeclareConstants(s.l)
		api := s.luaResourceApi(app)
		for _, f := range api {
			s.l.PushGoFunction(f.Function)
			s.l.SetGlobal(f.Name)
		}
	}
}

func (s *Script) luaRun(app *App, prom *Promise) {
	go func() {
		if s.l == nil {
			log.Panic("Script not initialized")
		}
		s.luaEval(app, s.Code, false)
		prom.CompleteWithValue(s)
	}()
}

func (s *Script) luaCall(f FieldAccessor, args []any, method bool) Future {
	prom := NewPromise()
	go func() {
		if s.l == nil {
			log.Panic("Script not initialized")
		}

		prom.Bind(CallFunction(s.l, f, args, method))
	}()
	return prom
}

func (s *Script) luaEval(app *App, code []byte, include bool) {
	prev := s.including
	s.including = include

	input := bytes.NewReader(code)
	if err := s.l.Load(input, "="+s.ref.String(), ""); err != nil {
		log.Panicf("Error loading script: %s", err)
	}
	if err := s.l.ProtectedCall(0, lua.MultipleReturns, 0); err != nil {
		log.Panicf("Error running script: %s", err)
	}
	s.declareGlobalObjects(app)
	s.including = prev
}

func (s *Script) declareGlobalObjects(app *App) {
	s.forEachDeclaredObject(func(typ, key string, included bool) {
		obj := withLuaTableAtIndex(s.l, -1)
		obj.SetString("id", key)

		if included {
			return
		}

		switch typ {
		case "actor":
			s.declareActor(app, key, obj)
		case "room":
			s.declareRoom(app, key, obj)
		}
	})
}

func (s *Script) declareActor(app *App, actorID string, actor luaTableUtils) {
	app.RunCommand(ActorDeclare{
		ActorID:   actorID,
		ActorName: actor.GetString("name"),
		Costume:   actor.GetRefOpt("costume", ResourceRefNull),
		ScriptLoc: WithField(actorID),
		Size:      actor.GetSizeOpt("size", DefaultActorSize),
		TalkColor: actor.GetColorOpt("talkcolor", DefaultActorTalkColor),
		UsePos:    actor.GetPositionOpt("usepos", DefaultActorUsePos),
		UseDir:    actor.GetDirectionOpt("usedir", DefaultActorDirection),
	}).Wait()
}

func (s *Script) declareRoom(app *App, roomID string, room luaTableUtils) {
	app.RunCommand(RoomDeclare{
		RoomID:          roomID,
		Script:          s,
		BackgroundRef:   room.GetRef("background"),
		CollisionMapRef: room.GetRef("collision_map"),
	}).Wait()

	room.IfTableFieldExists("objects", func(objs luaTableUtils) {
		objs.ForEach(func(key int, value int) {
			objID := lua.CheckString(s.l, key)
			obj := withLuaTableAtIndex(s.l, value).CheckObjectType("object")

			// Set the related IDs
			obj.SetString("room", roomID)
			obj.SetString("id", objID)

			cmd := ObjectDeclare{
				Class:     obj.GetClassOpt("class", 0),
				Hotspot:   obj.GetRectangle("hotspot"),
				Name:      lua.CheckString(s.l, key),
				ObjectID:  objID,
				Pos:       obj.GetPositionOpt("pos", NewPos(0, 0)),
				RoomID:    roomID,
				ScriptLoc: WithField(roomID, "objects", objID),
				Sprites:   obj.GetRefOpt("sprites", ResourceRefNull),
				UseDir:    obj.GetDirection("usedir"),
				UsePos:    obj.GetPosition("usepos"),
			}
			obj.IfTableFieldExists("states", func(states luaTableUtils) {
				states.ForEach(func(_ int, value int) {
					state := withLuaTableAtIndex(s.l, value)
					cmd.States = append(cmd.States, &ObjectState{
						Anim: state.GetAnimationOpt("anim", nil),
					})
				})
			})
			app.RunCommand(cmd).Wait()
		})
	})
}

func (s *Script) forEachDeclaredObject(f func(typ, key string, included bool)) {
	if s.l == nil {
		log.Panic("Script not initialized")
	}
	s.l.PushGlobalTable()
	defer s.l.Pop(1)

	s.l.PushNil()
	for s.l.Next(-2) {
		key := lua.CheckString(s.l, -2)
		tab := withLuaTableAtIndex(s.l, -1)
		if typ, ok := tab.ObjectType(); ok {
			f(typ, key, tab.GetBoolean("included"))
		}
		s.l.Pop(1)
	}
}

func (s *Script) luaResourceApi(app *App) []lua.RegistryFunction {
	return []lua.RegistryFunction{
		//
		// Resource construction functions
		//
		{Name: "actor", Function: func(l *lua.State) int {
			actor := withNewLuaObjectWrapping(l, 1, "actor")
			actor.SetBoolean("included", s.including)
			actor.SetFunction("say", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("actor")
				text := lua.CheckString(l, 2)
				opts := withLuaTableAtIndex(l, 3)
				cmd := ActorSpeak{
					Actor: self.GetActorByID(app, "id"),
					Text:  text,
					Delay: opts.GetDurationOpt("delay", DefaultActorSpeakDelay),
				}
				done := app.RunCommand(cmd)
				luaPushFuture(l, done)
				return 1
			}))
			actor.SetFunction("show", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("actor")
				opts := withLuaTableAtIndex(l, 2)
				cmd := ActorShow{
					Actor:    self.GetActorByID(app, "id"),
					Position: opts.GetPositionOpt("pos", DefaultActorPosition),
					LookAt:   opts.GetDirectionOpt("dir", DefaultActorDirection),
				}
				app.RunCommand(cmd).Wait()
				return 0
			}))
			actor.SetFunction("select", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("actor")
				cmd := ActorSelectEgo{
					Actor: self.GetActorByID(app, "id"),
				}
				done := app.RunCommand(cmd)
				luaPushFuture(l, done)
				return 1
			}))
			actor.SetFunction("stand", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("actor")
				opts := withLuaTableAtIndex(l, 2)
				cmd := ActorStand{
					Actor:     app.ActorByID(self.GetString("id")),
					Direction: opts.GetDirectionOpt("dir", DefaultActorDirection),
				}
				done := app.RunCommand(cmd)
				luaPushFuture(l, done)
				return 1
			}))
			actor.SetFunction("toinventory", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("actor")
				obj := withLuaTableAtIndex(l, 2)
				cmd := ActorAddToInventory{
					Actor:  self.GetActorByID(app, "id"),
					Object: obj.GetObjectByID(app, "room", "id"),
				}
				done := app.RunCommand(cmd)
				luaPushFuture(l, done)
				return 1
			}))
			actor.SetFunction("walkto", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("actor")
				pos := luaCheckPosition(l, 2)
				cmd := ActorWalkToPosition{
					Actor:    self.GetActorByID(app, "id"),
					Position: pos,
				}
				done := app.RunCommand(cmd)
				luaPushFuture(l, done)
				return 1
			}))
			return 1
		}},
		{Name: "class", Function: func(l *lua.State) int {
			opts := withLuaTableAtIndex(l, 1)
			class := withNewLuaObject(l, "class")
			class.SetInteger("mask", opts.GetInteger("mask"))
			return 1
		}},
		{Name: "costume", Function: func(l *lua.State) int {
			opts := withLuaTableAtIndex(l, 1)
			cost := withNewLuaObject(l, "costume")
			cost.SetResourceRef("ref", opts.GetRef("ref"))
			return 1
		}},
		{Name: "music", Function: func(l *lua.State) int {
			opts := withLuaTableAtIndex(l, 1)
			music := withNewLuaObject(l, "music")
			music.SetResourceRef("ref", opts.GetRef("ref"))
			music.SetFunction("play", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("music")
				done := app.RunCommand(MusicPlay{
					MusicResource: self.GetRef("ref"),
				})
				luaPushFuture(l, done)
				return 1
			}))
			return 1
		}},
		{Name: "object", Function: func(l *lua.State) int {
			obj := withNewLuaObjectWrapping(l, 1, "object")
			obj.SetFunction("owner", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("object")
				obj := app.FindObject(self.GetString("room"), self.GetString("id"))
				if owner := obj.Owner(); owner == nil {
					l.PushNil()
				} else {
					l.Global(owner.ID())
				}
				return 1
			}))
			return 1
		}},
		{Name: "room", Function: func(l *lua.State) int {
			room := withNewLuaObjectWrapping(l, 1, "room")
			room.SetBoolean("included", s.including)
			room.SetFunction("show", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("room")
				done := app.RunCommand(RoomShow{
					Room: self.GetRoomByID(app, "id"),
				})
				luaPushFuture(l, done)
				return 1
			}))
			return 1
		}},
		{Name: "sound", Function: func(l *lua.State) int {
			opts := withLuaTableAtIndex(l, 1)
			sound := withNewLuaObject(l, "sound")
			sound.SetResourceRef("ref", opts.GetRef("ref"))
			sound.SetFunction("play", lua.Function(func(l *lua.State) int {
				self := withLuaTableAtIndex(l, 1).CheckObjectType("sound")
				done := app.RunCommand(SoundPlay{
					SoundResource: self.GetRef("ref"),
				})
				luaPushFuture(l, done)
				return 1
			}))
			return 1
		}},
		{Name: "var", Function: func(l *lua.State) int {
			// TODO: declare the variable in the app and bind the getter and setter
			v := withNewLuaObject(l, "var")
			v.SetFunction("get", lua.Function(func(l *lua.State) int {
				lua.Errorf(l, "not implemented")
				return 0
			}))
			v.SetFunction("set", lua.Function(func(l *lua.State) int {
				lua.Errorf(l, "not implemented")
				return 0
			}))
			return 1
		}},

		//
		// API functions
		//
		{Name: "include", Function: func(l *lua.State) int {
			ref := luaCheckResourceRef(l, 1)
			if luaIsIncluded(l, ref) {
				return 0
			}
			script, err := WaitAs[*Script](app.RunCommand(ScriptRun{
				ScriptRef: ref,
			}))
			if err != nil {
				lua.Errorf(l, "Error including script: %s", err)
				return 0
			}
			s.luaEval(app, script.Code, true)

			luaSetIncluded(l, ref)
			return 0
		}},

		// TODO: this function uses the ShowDialog. It must be replaced by a function that
		// prints with some selected font.
		{Name: "sayline", Function: func(l *lua.State) int {
			text := lua.CheckString(l, 1)
			opts := withLuaTableAtIndex(l, 2)
			cmd := ShowDialog{
				Text:     text,
				Position: opts.GetPositionOpt("pos", DefaultDialogPosition),
				Color:    opts.GetColorOpt("color", DefaultDialogColor),
			}
			done := app.RunCommand(cmd)
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "sleep", Function: func(l *lua.State) int {
			time.Sleep(luaCheckDurationMillis(l, 1))
			return 0
		}},
		{Name: "userputon", Function: func(l *lua.State) int {
			done := app.RunCommand(EnableControlPanel{Enable: true})
			done.Wait()
			return 0
		}},
		{Name: "userputoff", Function: func(l *lua.State) int {
			done := app.RunCommand(EnableControlPanel{Enable: false})
			done.Wait()
			return 0
		}},
		{Name: "cursoron", Function: func(l *lua.State) int {
			done := app.RunCommand(EnableMouseCursor{Enable: true})
			done.Wait()
			return 0
		}},
		{Name: "cursoroff", Function: func(l *lua.State) int {
			done := app.RunCommand(EnableMouseCursor{Enable: false})
			done.Wait()
			return 0
		}},
	}
}

// Call the function pointed by f with the given arguments. If method is true, the first argument
// is the object that has the method.
func CallFunction(l *lua.State, f FieldAccessor, args []any, method bool) Future {
	prom := NewPromise()
	err := luaPushField(l, f)
	if err != nil {
		prom.CompleteWithError(err)
		return prom
	}
	if !l.IsFunction(-1) {
		prom.CompleteWithErrorf("Field %s is not a function", f)
		return prom
	}
	nargs := len(args)
	if method {
		luaPushValue(l, f.Base())
		nargs++
	}
	for _, arg := range args {
		luaPushValue(l, arg)
	}
	l.Call(nargs, 0)
	prom.Complete()
	return prom
}

func luaDeclareConstants(l *lua.State) {
	for k, pushFunc := range map[string]func(){
		// Defaults table. Needed to declare default actions.
		"DEFAULT": func() { l.NewTable() },

		// Predefined directions
		"UP":    func() { l.PushInteger(int(DirUp)) },
		"RIGHT": func() { l.PushInteger(int(DirRight)) },
		"DOWN":  func() { l.PushInteger(int(DirDown)) },
		"LEFT":  func() { l.PushInteger(int(DirLeft)) },

		// Predefined classes
		"PERSON":     func() { luaPushClass(l, ObjectClassPerson) },
		"PICKABLE":   func() { luaPushClass(l, ObjectClassPickable) },
		"OPENABLE":   func() { luaPushClass(l, ObjectClassOpenable) },
		"CLOSEABLE":  func() { luaPushClass(l, ObjectClassCloseable) },
		"APPLICABLE": func() { luaPushClass(l, ObjectClassApplicable) },
	} {
		pushFunc()
		l.SetGlobal(k)
	}
}

func luaIsIncluded(l *lua.State, ref ResourceRef) bool {
	l.Field(lua.RegistryIndex, "__pctk_includes")
	defer l.Pop(1)
	if l.IsNil(-1) {
		return false
	}

	l.Field(-1, ref.String())
	defer l.Pop(1)
	if l.IsNil(-1) {
		return false
	}
	return true
}

func luaSetIncluded(l *lua.State, ref ResourceRef) {
	l.Field(lua.RegistryIndex, "__pctk_includes")
	defer l.SetField(lua.RegistryIndex, "__pctk_includes")
	if l.IsNil(-1) {
		l.Pop(1)
		l.NewTable()
	}
	l.PushBoolean(true)
	l.SetField(-2, ref.String())
}

func luaCheckColor(l *lua.State, index int) (col Color) {
	tab := withLuaTableAtIndex(l, index)
	col.R = byte(tab.GetInteger("r"))
	col.G = byte(tab.GetInteger("g"))
	col.B = byte(tab.GetInteger("b"))
	col.A = byte(tab.GetIntegerOpt("a", 255))
	return
}

func luaCheckClass(l *lua.State, index int) ObjectClass {
	tab := withLuaTableAtIndex(l, index).CheckObjectType("class")
	return ObjectClass(tab.GetInteger("mask"))
}

func luaCheckDurationMillis(l *lua.State, index int) time.Duration {
	val := lua.CheckInteger(l, index)
	return time.Duration(val) * time.Millisecond
}

func luaCheckPosition(l *lua.State, index int) (pos Position) {
	tab := withLuaTableAtIndex(l, index)
	pos.X = tab.GetInteger("x")
	pos.Y = tab.GetInteger("y")
	return
}

func luaCheckSize(l *lua.State, index int) (size Size) {
	tab := withLuaTableAtIndex(l, index)
	size.W = tab.GetInteger("w")
	size.H = tab.GetInteger("h")
	return
}

func luaCheckRectangle(l *lua.State, index int) (rect Rectangle) {
	tab := withLuaTableAtIndex(l, index)
	rect.Pos.X = tab.GetInteger("x")
	rect.Pos.Y = tab.GetInteger("y")
	rect.Size.W = tab.GetInteger("w")
	rect.Size.H = tab.GetInteger("h")
	return
}

func luaCheckResourceRef(l *lua.State, index int) ResourceRef {
	val := lua.CheckString(l, index)
	ref, err := ParseResourceRef(val)
	if err != nil {
		lua.ArgumentError(l, index, "invalid resource reference")
	}
	return ref
}

func luaCheckAnimation(l *lua.State, index int) (anim *Animation) {
	tab := withLuaTableAtIndex(l, index)

	anim = NewAnimation()
	tab.ForEach(func(key int, value int) {
		frame := withLuaTableAtIndex(l, value)
		anim.AddFrames(
			frame.GetDuration("delay"),
			frame.GetInteger("row"),
			frame.GetIntegers("seq")...,
		)
	})
	return anim
}

func luaPushColor(l *lua.State, c Color) {
	l.NewTable()
	for k, v := range map[string]int{
		"r": int(c.R),
		"g": int(c.G),
		"b": int(c.B),
		"a": int(c.A),
	} {
		l.PushInteger(v)
		l.SetField(-2, k)
	}
}

func luaPushFuture(l *lua.State, f Future) {
	lua.NewLibrary(l, []lua.RegistryFunction{
		{Name: "iscompleted", Function: func(l *lua.State) int {
			l.PushBoolean(f.IsCompleted())
			return 1
		}},
		{Name: "wait", Function: func(l *lua.State) int {
			f.Wait()
			return 0
		}},
	})
	return
}

func luaPushClass(l *lua.State, c ObjectClass) {
	class := withNewLuaObject(l, "class")
	class.SetInteger("mask", int(c))
}

func luaPushField(l *lua.State, f FieldAccessor) error {
	var n int
	str := FieldAccessor(f).String()
	l.PushGlobalTable()
	n++
	for ; len(f) > 1; f = f[1:] {
		l.Field(-1, f[0])
		n++
		if !l.IsTable(-1) {
			l.Pop(n)
			return fmt.Errorf("Object %s not found", str)
		}
	}
	l.Field(-1, f[0])
	n++
	if l.IsNil(-1) {
		l.Pop(n)
		return fmt.Errorf("Field %s not found in object %s", f[0], str)
	}
	for i := 1; i < n; i++ {
		l.Remove(-2)
	}
	return nil
}

func luaPushValue(l *lua.State, val any) {
	switch v := val.(type) {
	case bool:
		l.PushBoolean(v)
	case int:
		l.PushInteger(v)
	case string:
		l.PushString(v)
	case time.Duration:
		l.PushInteger(int(v / time.Millisecond))
	case FieldAccessor:
		if err := luaPushField(l, v); err != nil {
			lua.ArgumentError(l, 1, err.Error())
		}
	default:
		log.Panicf("Unsupported value type: %T", val)
	}
}

type luaTableUtils struct {
	l     *lua.State
	index int
}

func withLuaTableAtIndex(l *lua.State, index int) luaTableUtils {
	// Convert top-relative index to absolute
	if index < 0 {
		index = l.Top() + index + 1
	}
	return luaTableUtils{l, index}
}

func withNewLuaTable(l *lua.State) luaTableUtils {
	l.NewTable()
	return withLuaTableAtIndex(l, -1)
}

func withNewLuaObject(l *lua.State, typ string) luaTableUtils {
	l.NewTable()
	l.PushString(typ)
	l.SetField(-2, "__type")
	return withLuaTableAtIndex(l, -1)
}

func withNewLuaObjectWrapping(l *lua.State, index int, typ string) luaTableUtils {
	l.NewTable() // object table
	l.PushString(typ)
	l.SetField(-2, "__type")
	l.NewTable() // object metatable
	if index < 0 {
		index -= 2 // adjust if index is top relative
	}
	l.PushValue(index)
	l.SetField(-2, "__index")
	l.SetMetaTable(-2)
	return withLuaTableAtIndex(l, -1)
}

func (t luaTableUtils) IfTableFieldExists(key string, then func(luaTableUtils)) {
	t.l.Field(t.index, key)
	defer t.l.Pop(1)
	if t.l.IsNil(-1) {
		return
	}
	then(withLuaTableAtIndex(t.l, -1))
}

func (t luaTableUtils) ForEach(then func(key int, value int)) {
	t.l.PushNil()
	for t.l.Next(t.index) {
		then(-2, -1)
		t.l.Pop(1)
	}
}

func (t luaTableUtils) GetString(key string) (val string) {
	t.getField(key, lua.TypeString, func() { val = lua.CheckString(t.l, -1) })
	return
}

func (t luaTableUtils) GetStringOpt(key string, def string) (val string) {
	val = def
	t.getFieldOpt(key, lua.TypeString, func() { val = lua.CheckString(t.l, -1) })
	return
}

func (t luaTableUtils) GetInteger(key string) (val int) {
	t.getField(key, lua.TypeNumber, func() { val = lua.CheckInteger(t.l, -1) })
	return
}

func (t luaTableUtils) GetIntegerOpt(key string, def int) (val int) {
	val = def
	t.getFieldOpt(key, lua.TypeNumber, func() { val = lua.CheckInteger(t.l, -1) })
	return
}

func (t luaTableUtils) GetIntegers(key string) (val []int) {
	t.getField(key, lua.TypeTable, func() {
		tab := withLuaTableAtIndex(t.l, -1)
		tab.ForEach(func(_, value int) {
			val = append(val, lua.CheckInteger(t.l, value))
		})
	})
	return
}

func (t luaTableUtils) GetBoolean(key string) (val bool) {
	t.getField(key, lua.TypeNone, func() { val = t.l.ToBoolean(-1) })
	return
}

func (t luaTableUtils) GetDuration(key string) (val time.Duration) {
	t.getField(key, lua.TypeNumber, func() {
		val = time.Duration(lua.CheckInteger(t.l, -1)) * time.Millisecond
	})
	return
}

func (t luaTableUtils) GetDurationOpt(key string, def time.Duration) (val time.Duration) {
	val = def
	t.getFieldOpt(key, lua.TypeNumber, func() {
		val = time.Duration(lua.CheckInteger(t.l, -1)) * time.Millisecond
	})
	return
}

func (t luaTableUtils) GetRef(key string) (val ResourceRef) {
	t.getField(key, lua.TypeString, func() { val = luaCheckResourceRef(t.l, -1) })
	return
}

func (t luaTableUtils) GetRefOpt(key string, def ResourceRef) (val ResourceRef) {
	val = def
	t.getFieldOpt(key, lua.TypeString, func() { val = luaCheckResourceRef(t.l, -1) })
	return
}

func (t luaTableUtils) GetDirection(key string) (val Direction) {
	t.getField(key, lua.TypeNumber, func() { val = Direction(lua.CheckInteger(t.l, -1)) })
	return
}

func (t luaTableUtils) GetDirectionOpt(key string, def Direction) (val Direction) {
	val = def
	t.getFieldOpt(key, lua.TypeNumber, func() { val = Direction(lua.CheckInteger(t.l, -1)) })
	return
}

func (t luaTableUtils) GetClass(key string) (val ObjectClass) {
	t.getField(key, lua.TypeTable, func() { val = luaCheckClass(t.l, -1) })
	return
}

func (t luaTableUtils) GetClassOpt(key string, def ObjectClass) (val ObjectClass) {
	val = def
	t.getFieldOpt(key, lua.TypeTable, func() { val = luaCheckClass(t.l, -1) })
	return
}

func (t luaTableUtils) GetColor(key string) (val Color) {
	t.getField(key, lua.TypeTable, func() { val = luaCheckColor(t.l, -1) })
	return
}

func (t luaTableUtils) GetColorOpt(key string, def Color) (val Color) {
	val = def
	t.getFieldOpt(key, lua.TypeTable, func() { val = luaCheckColor(t.l, -1) })
	return
}

func (t luaTableUtils) GetPosition(key string) (val Position) {
	t.getField(key, lua.TypeTable, func() { val = luaCheckPosition(t.l, -1) })
	return
}

func (t luaTableUtils) GetPositionOpt(key string, def Position) (val Position) {
	val = def
	t.getFieldOpt(key, lua.TypeTable, func() { val = luaCheckPosition(t.l, -1) })
	return
}

func (t luaTableUtils) GetSize(key string) (val Size) {
	t.getField(key, lua.TypeTable, func() { val = luaCheckSize(t.l, -1) })
	return
}

func (t luaTableUtils) GetSizeOpt(key string, def Size) (val Size) {
	val = def
	t.getFieldOpt(key, lua.TypeTable, func() { val = luaCheckSize(t.l, -1) })
	return
}

func (t luaTableUtils) GetRectangle(key string) (val Rectangle) {
	t.getField(key, lua.TypeTable, func() {
		val = luaCheckRectangle(t.l, -1)
	})
	return
}

func (t luaTableUtils) GetAnimation(key string) (val *Animation) {
	t.getField(key, lua.TypeTable, func() {
		val = luaCheckAnimation(t.l, -1)
	})
	return
}

func (t luaTableUtils) GetAnimationOpt(key string, def *Animation) (val *Animation) {
	val = def
	t.getFieldOpt(key, lua.TypeTable, func() {
		val = luaCheckAnimation(t.l, -1)
	})
	return
}
func (t luaTableUtils) GetActorByID(app *App, key string) (val *Actor) {
	actorID := t.GetString(key)
	val = app.ActorByID(actorID)
	if val == nil {
		lua.ArgumentError(t.l, 1, fmt.Sprintf("actor %s not found", actorID))
	}
	return
}

func (t luaTableUtils) GetObjectByID(app *App, roomk, idk string) (val *Object) {
	roomID := t.GetString(roomk)
	objID := t.GetString(idk)
	val = app.FindObject(roomID, objID)
	if val == nil {
		lua.ArgumentError(t.l, 1, fmt.Sprintf("object %s not found in room %s", objID, roomID))
	}
	return
}

func (t luaTableUtils) GetRoomByID(app *App, key string) (val *Room) {
	roomID := t.GetString(key)
	val = app.FindRoom(roomID)
	if val == nil {
		lua.ArgumentError(t.l, 1, fmt.Sprintf("room %s not found", roomID))
	}
	return
}

func (t luaTableUtils) getField(key string, expected lua.Type, pull func()) {
	t.l.Field(t.index, key)
	defer t.l.Pop(1)
	if given := t.l.TypeOf(-1); expected != lua.TypeNone && given != expected {
		lua.ArgumentError(t.l, t.index, fmt.Sprintf(
			"field '%s' has type '%s', '%s' expected", key, given, expected))
	}
	pull()
}

func (t luaTableUtils) getFieldOpt(key string, expected lua.Type, pull func()) {
	if !t.l.IsTable(t.index) {
		return
	}
	t.l.Field(t.index, key)
	defer t.l.Pop(1)
	if t.l.IsNil(-1) {
		return
	}
	if given := t.l.TypeOf(-1); expected != lua.TypeNone && given != expected {
		lua.ArgumentError(t.l, t.index, fmt.Sprintf(
			"field '%s' has type '%s', '%s' expected", key, given, expected))
	}
	pull()
}

func (t luaTableUtils) SetObjectType(typ string) {
	t.l.PushString(typ)
	t.l.SetField(t.index, "__type")
}

func (t luaTableUtils) SetString(key, value string) {
	t.setField(key, func() { t.l.PushString(value) })
}

func (t luaTableUtils) SetInteger(key string, value int) {
	t.setField(key, func() { t.l.PushInteger(value) })
}

func (t luaTableUtils) SetBoolean(key string, value bool) {
	t.setField(key, func() { t.l.PushBoolean(value) })
}

func (t luaTableUtils) SetFunction(key string, value lua.Function) {
	t.setField(key, func() { t.l.PushGoFunction(value) })
}

func (t luaTableUtils) SetColor(key string, value Color) {
	t.setField(key, func() { luaPushColor(t.l, value) })
}

func (t luaTableUtils) SetDirection(key string, value Direction) {
	t.setField(key, func() { t.l.PushInteger(int(value)) })
}

func (t luaTableUtils) SetResourceRef(key string, value ResourceRef) {
	t.setField(key, func() { t.l.PushString(value.String()) })
}

func (t luaTableUtils) setField(key string, push func()) {
	push()
	index := t.index
	if index < 0 {
		index--
	}
	t.l.SetField(index, key)
}

func (t luaTableUtils) CheckObjectType(expected string) luaTableUtils {
	t.l.Field(t.index, "__type")
	defer t.l.Pop(1)

	if t.l.TypeOf(-1) != lua.TypeString {
		lua.ArgumentError(t.l, t.index, fmt.Sprintf("object of type %s expected", expected))
	}
	if actual := lua.CheckString(t.l, -1); actual != expected {
		lua.ArgumentError(t.l, t.index, fmt.Sprintf(
			"object of type %s expected, got %s", expected, actual,
		))
	}
	return t
}

func (t luaTableUtils) ObjectType() (typ string, ok bool) {
	if !t.l.IsTable(t.index) {
		return "", false
	}
	t.l.Field(t.index, "__type")
	defer t.l.Pop(1)
	return t.l.ToString(-1)
}

func (t luaTableUtils) IsTable() bool {
	return t.l.IsTable(t.index)
}

func (t luaTableUtils) IsObject() bool {
	if !t.l.IsTable(t.index) {
		return false
	}
	t.l.Field(t.index, "__type")
	defer t.l.Pop(1)
	return t.l.IsString(-1)
}

func (t luaTableUtils) IsObjectType(typ string) bool {
	if !t.l.IsTable(t.index) {
		return false
	}
	t.l.Field(t.index, "__type")
	defer t.l.Pop(1)
	if !t.l.IsString(-1) {
		return false
	}
	return lua.CheckString(t.l, -1) == typ
}
