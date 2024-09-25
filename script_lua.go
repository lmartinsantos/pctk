package pctk

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/go-lua"
)

func (s *Script) luaInit(app AppContext) {
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

func (s *Script) luaRun(app AppContext, prom *Promise) {
	go func() {
		if s.l == nil {
			log.Panic("Script not initialized")
		}
		s.luaEval(app, s.Code, false)
		prom.CompleteWithValue(s)
	}()
}

func (s *Script) luaCall(object, method string, prom *Promise) {
	go func() {
		if s.l == nil {
			log.Panic("Script not initialized")
		}

		s.l.Global(object)
		if !s.l.IsTable(-1) {
			log.Panicf("Object %s not found", object)
		}
		s.l.Field(-1, method)
		if !s.l.IsFunction(-1) {
			log.Panicf("Method %s not found in object %s", method, object)
		}
		s.l.PushValue(-2)
		s.l.Call(1, 0)
		prom.Complete()
	}()
}

func (s *Script) luaEval(app AppContext, code []byte, include bool) {
	prev := s.including
	s.including = include

	input := bytes.NewReader(code)
	if err := s.l.Load(input, "="+s.ref.String(), ""); err != nil {
		log.Panicf("Error loading script: %s", err)
	}
	if err := s.l.ProtectedCall(0, lua.MultipleReturns, 0); err != nil {
		log.Panicf("Error running script: %s", err)
	}

	s.forEachDeclaredObject(func(typ, key string, included bool) {
		// Define the ID from the object key
		s.l.PushString(key)
		s.l.SetField(-2, "id")

		// For rooms that are not included, declare them in the app
		if typ == "room" && !included {
			app.Do(RoomDeclare{
				RoomID:        key,
				Script:        s,
				BackgroundRef: luaCheckField(s.l, -1, "background", luaToResourceRef),
			}).Wait()
		}
	})

	s.including = prev
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
		if typ, ok := luaObjectType(s.l, -1); ok {
			s.l.Field(-1, "included")
			included := s.l.ToBoolean(-1)
			s.l.Pop(1)

			f(typ, key, included)
		}
		s.l.Pop(1)
	}
}

func (s *Script) luaResourceApi(app AppContext) []lua.RegistryFunction {
	return []lua.RegistryFunction{
		//
		// Resource construction functions
		//
		{Name: "actor", Function: func(l *lua.State) int {
			luaPushObject(l, "actor", map[string]any{
				"say": lua.Function(func(l *lua.State) int {
					cmd := ActorSpeak{
						ActorID: luaCheckObjectField(l, 1, "actor", "id", (*lua.State).ToString),
						Text:    lua.CheckString(l, 2),
						Delay:   luaCheckOption(l, 3, "delay", DefaultActorSpeakDelay, luaToDurationMillis),
					}
					done := app.Do(cmd)
					luaPushFuture(l, done)
					return 1
				}),
				"show": lua.Function(func(l *lua.State) int {
					luaCheckObjectType(l, 1, "actor")
					cmd := ActorShow{
						ActorID:         luaCheckObjectField(l, 1, "actor", "id", (*lua.State).ToString),
						Position:        luaCheckOption(l, 2, "pos", DefaultActorPosition, luaToPosition),
						LookAt:          luaCheckOption(l, 2, "dir", DefaultActorDirection, luaToDirection),
						CostumeResource: luaCheckOption(l, 2, "costume", ResourceRefNull, luaToResourceRef),
					}
					done := app.Do(cmd)
					luaPushFuture(l, done)
					return 1
				}),
				"select": lua.Function(func(l *lua.State) int {
					cmd := ActorSelectEgo{
						ActorID: luaCheckObjectField(l, 1, "actor", "id", (*lua.State).ToString),
					}
					done := app.Do(cmd)
					luaPushFuture(l, done)
					return 1
				}),
				"stand": lua.Function(func(l *lua.State) int {
					cmd := ActorStand{
						ActorID:   luaCheckObjectField(l, 1, "actor", "id", (*lua.State).ToString),
						Direction: luaCheckOption(l, 2, "dir", DefaultActorDirection, luaToDirection),
					}
					done := app.Do(cmd)
					luaPushFuture(l, done)
					return 1
				}),
				"walkto": lua.Function(func(l *lua.State) int {
					cmd := ActorWalkToPosition{
						ActorID:  luaCheckObjectField(l, 1, "actor", "id", (*lua.State).ToString),
						Position: luaCheckPosition(l, 2),
					}
					done := app.Do(cmd)
					luaPushFuture(l, done)
					return 1
				}),
			})
			return 1
		}},
		{Name: "class", Function: func(l *lua.State) int {
			luaPushObject(l, "class", map[string]any{
				"mask": luaCheckField(l, 1, "mask", (*lua.State).ToInteger),
			})
			return 1
		}},
		{Name: "costume", Function: func(l *lua.State) int {
			luaPushObject(l, "costume", map[string]any{
				"ref": luaCheckField(l, 1, "ref", luaToResourceRef),
			})
			return 1
		}},
		{Name: "include", Function: func(l *lua.State) int {
			ref := luaCheckResourceRef(l, 1)
			if luaIsIncluded(l, ref) {
				return 0
			}
			script := app.Do(ScriptRun{
				ScriptRef: ref,
			}).Wait().(*Script)
			s.luaEval(app, script.Code, true)

			luaSetIncluded(l, ref)
			return 0
		}},
		{Name: "music", Function: func(l *lua.State) int {
			luaPushObject(l, "music", map[string]any{
				"ref": luaCheckField(l, 1, "ref", luaToResourceRef),
				"play": lua.Function(func(l *lua.State) int {
					done := app.Do(MusicPlay{
						MusicResource: luaCheckObjectField(l, 1, "music", "ref", luaToResourceRef),
					})
					luaPushFuture(l, done)
					return 1
				}),
			})
			return 1
		}},
		{Name: "room", Function: func(l *lua.State) int {
			luaWrapTableAsObject(l, 1, "room")
			luaSetField(l, -1, "show", lua.Function(func(l *lua.State) int {
				done := app.Do(RoomShow{
					RoomID: luaCheckObjectField(l, 1, "room", "id", (*lua.State).ToString),
				})
				luaPushFuture(l, done)
				return 1
			}))
			luaSetField(l, -1, "included", s.including)

			if !s.including {
				// TODO: process the room definitions (objects, etc).
			}
			return 1
		}},
		{Name: "sound", Function: func(l *lua.State) int {
			luaPushObject(l, "sound", map[string]any{
				"ref": luaCheckField(l, 1, "ref", luaToResourceRef),
				"play": lua.Function(func(l *lua.State) int {
					done := app.Do(SoundPlay{
						SoundResource: luaCheckObjectField(l, 1, "sound", "ref", luaToResourceRef),
					})
					luaPushFuture(l, done)
					return 1
				}),
			})
			return 1
		}},
		{Name: "var", Function: func(l *lua.State) int {
			// TODO: declare the variable in the app and bind the getter and setter
			luaPushObject(l, "var", map[string]any{
				"get": lua.Function(func(l *lua.State) int {
					lua.Errorf(l, "not implemented")
					return 0
				}),
				"set": lua.Function(func(l *lua.State) int {
					lua.Errorf(l, "not implemented")
					return 0
				}),
			})
			return 1
		}},

		//
		// API functions
		//

		// TODO: this function uses the ShowDialog. It must be replaced by a function that
		// prints with some selected font.
		{Name: "sayline", Function: func(l *lua.State) int {
			cmd := ShowDialog{
				Text:     lua.CheckString(l, 1),
				Position: luaCheckOption(l, 2, "pos", DefaultDialogPosition, luaToPosition),
				Color:    luaCheckOption(l, 2, "color", White, luaToColor),
			}
			done := app.Do(cmd)
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "sleep", Function: func(l *lua.State) int {
			time.Sleep(luaCheckDurationMillis(l, 1))
			return 0
		}},
		{Name: "usercontrol", Function: func(l *lua.State) int {
			done := app.Do(EnableControlPanel{Enable: l.ToBoolean(1)})
			luaPushFuture(l, done)
			return 1
		}},
	}
}

func luaDeclareConstants(l *lua.State) {
	for k, pushFunc := range map[string]func(){
		"up":    func() { l.PushInteger(int(DirUp)) },
		"right": func() { l.PushInteger(int(DirRight)) },
		"down":  func() { l.PushInteger(int(DirDown)) },
		"left":  func() { l.PushInteger(int(DirLeft)) },
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

func luaIfFieldExists(l *lua.State, index int, field string, f func()) {
	if l.IsTable(index) {
		l.Field(index, field)
		if l.TypeOf(-1) != lua.TypeNil {
			f()
		}
		l.Pop(1)
	}
}

func luaToByte(l *lua.State, index int) (byte, bool) {
	val, ok := l.ToInteger(index)
	return byte(val), ok
}

func luaToColor(l *lua.State, index int) (col Color, ok bool) {
	col.R, ok = luaFieldTo(l, index, "r", luaToByte)
	if !ok {
		return
	}
	col.G, ok = luaFieldTo(l, index, "g", luaToByte)
	if !ok {
		return
	}
	col.B, ok = luaFieldTo(l, index, "b", luaToByte)
	if !ok {
		return
	}
	col.A, ok = luaFieldTo(l, index, "a", luaToByte)
	if !ok {
		col.A = 255
		ok = true
	}
	return
}

func luaToDirection(l *lua.State, index int) (Direction, bool) {
	val, ok := l.ToInteger(index)
	if !ok {
		return 0, false
	}
	return Direction(val), true
}

func luaCheckDurationMillis(l *lua.State, index int) time.Duration {
	val, ok := luaToDurationMillis(l, index)
	if !ok {
		lua.ArgumentError(l, index, "duration expected")
	}
	return val
}

func luaToDurationMillis(l *lua.State, index int) (time.Duration, bool) {
	val, ok := l.ToInteger(index)
	if !ok {
		return 0, false
	}
	return time.Duration(val) * time.Millisecond, true
}

func luaCheckField[T any](l *lua.State, index int, field string, f func(*lua.State, int) (T, bool)) T {
	if !l.IsTable(index) {
		lua.ArgumentError(l, index, "table expected")
	}
	l.Field(index, field)
	defer l.Pop(1)
	val, ok := f(l, -1)
	if !ok {
		lua.ArgumentError(l, index, fmt.Sprintf("required field %s not found", field))
	}
	return val
}

func luaFieldTo[T any](l *lua.State, index int, field string, f func(*lua.State, int) (T, bool)) (val T, ok bool) {
	if !l.IsTable(index) {
		return
	}
	l.Field(index, field)
	defer l.Pop(1)
	val, ok = f(l, -1)
	return
}

func luaCheckOption[T any](l *lua.State, index int, field string, def T, f func(*lua.State, int) (T, bool)) T {
	luaIfFieldExists(l, index, field, func() {
		var ok bool
		def, ok = f(l, -1)
		if !ok {
			lua.ArgumentError(l, index, fmt.Sprintf("invalid value for field %s", field))
		}
	})
	return def
}

func luaCheckPosition(l *lua.State, index int) Position {
	pos, ok := luaToPosition(l, index)
	if !ok {
		lua.ArgumentError(l, index, "position expected")
	}
	return pos
}

func luaToPosition(l *lua.State, index int) (pos Position, ok bool) {
	pos.X, ok = luaFieldTo(l, index, "x", (*lua.State).ToInteger)
	if !ok {
		return
	}
	pos.Y, ok = luaFieldTo(l, index, "y", (*lua.State).ToInteger)
	return
}

func luaCheckResourceRef(l *lua.State, index int) ResourceRef {
	ref, ok := luaToResourceRef(l, index)
	if !ok {
		lua.ArgumentError(l, index, "resource reference expected")
	}
	return ref
}

func luaToResourceRef(l *lua.State, index int) (ResourceRef, bool) {
	val, ok := l.ToString(index)
	if !ok {
		return ResourceRefNull, false
	}
	ref, err := ParseResourceRef(val)
	if err != nil {
		return ResourceRefNull, false
	}
	return ref, true
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

func luaCheckObjectField[T any](
	l *lua.State,
	index int,
	typ string,
	field string,
	f func(*lua.State, int) (T, bool),
) T {
	luaCheckObjectType(l, index, typ)
	l.Field(index, field)
	defer l.Pop(1)
	val, ok := f(l, -1)
	if !ok {
		lua.ArgumentError(l, index, fmt.Sprintf("required field %s not found", field))
	}
	return val
}

func luaCheckObjectType(l *lua.State, index int, expected string) {
	if !l.IsTable(index) {
		lua.ArgumentError(l, index, fmt.Sprintf("object of type %s expected", expected))
	}
	l.Field(index, "__type")
	defer l.Pop(1)

	actual, ok := l.ToString(-1)
	if !ok {
		lua.ArgumentError(l, index, fmt.Sprintf("object of type %s expected", expected))
	}
	if actual != expected {
		lua.ArgumentError(l, index, fmt.Sprintf(
			"object of type %s expected, got %s", expected, actual,
		))
	}
}

func luaObjectType(l *lua.State, index int) (string, bool) {
	if !l.IsTable(index) {
		return "", false
	}
	l.Field(index, "__type")
	defer l.Pop(1)
	return l.ToString(-1)
}

func luaIsObjectType(l *lua.State, index int, expected string) bool {
	if !l.IsTable(index) {
		return false
	}
	l.Field(index, "__type")
	defer l.Pop(1)

	actual, ok := l.ToString(-1)
	if !ok {
		return false
	}
	return actual == expected
}

func luaPushObject(l *lua.State, typ string, fields map[string]any) {
	fields["__type"] = typ
	luaPushTable(l, fields)
}

func luaPushTable(l *lua.State, fields map[string]any) {
	l.NewTable()
	for k, v := range fields {
		luaSetField(l, -1, k, v)
	}
}

func luaWrapTableAsObject(l *lua.State, index int, typ string) {
	l.NewTable()
	l.PushString(typ)
	l.SetField(-2, "__type")
	l.NewTable()
	l.PushValue(index)
	l.SetField(-2, "__index")
	l.SetMetaTable(-2)
}

func luaSetField(l *lua.State, index int, key string, value any) {
	if !l.IsTable(index) {
		lua.Errorf(l, "table expected")
	}
	switch v := value.(type) {
	case int:
		l.PushInteger(v)
	case string:
		l.PushString(v)
	case bool:
		l.PushBoolean(v)
	case lua.Function:
		l.PushGoFunction(v)
	case Color:
		luaPushColor(l, v)
	case Direction:
		l.PushInteger(int(v))
	case ResourceRef:
		l.PushString(v.String())
	default:
		log.Panicf("unsupported type: %T", v)
	}
	l.SetField(index-1, key)
}
