package pctk

import (
	"log"
	"time"

	"github.com/Shopify/go-lua"
)

func (s *Script) runLua(app *App, prom Promise) {
	go func() {
		l := lua.NewState()
		lua.BaseOpen(l)

		luaDeclareConstants(l)
		funcs := luaApiFunctions(app)
		for _, f := range funcs {
			l.PushGoFunction(f.Function)
			l.SetGlobal(f.Name)
		}

		if err := lua.DoString(l, string(s.Code)); err != nil {
			log.Panicf("Error running script: %s", err)
		}
		prom.Complete()
	}()
}

func luaApiFunctions(app *App) []lua.RegistryFunction {
	return []lua.RegistryFunction{
		{Name: "ActorShow", Function: func(l *lua.State) int {
			cmd := ActorShow{
				ActorResource:   ResourceLocator(lua.CheckString(l, 1)),
				ActorName:       lua.CheckString(l, 2),
				Position:        luaCheckOption(l, 3, "pos", DefaultActorPosition, luaCheckPosition),
				LookAt:          luaCheckOption(l, 3, "dir", DefaultActorDirection, luaCheckDirection),
				CostumeResource: luaCheckOption(l, 3, "costume", "", luaCheckResourceLocator),
			}
			done := app.Do(cmd)
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "ActorSpeak", Function: func(l *lua.State) int {
			cmd := ActorSpeak{
				ActorName: lua.CheckString(l, 1),
				Text:      lua.CheckString(l, 2),
				Delay:     luaCheckOption(l, 3, "delay", DefaultActorSpeakDelay, luaCheckDurationMillis),
			}
			done := app.Do(cmd)
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "ActorStand", Function: func(l *lua.State) int {
			cmd := ActorStand{
				ActorName: lua.CheckString(l, 1),
				Direction: luaCheckOption(l, 2, "dir", DefaultActorDirection, luaCheckDirection),
			}
			done := app.Do(cmd)
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "ActorWalkToPosition", Function: func(l *lua.State) int {
			cmd := ActorWalkToPosition{
				ActorName: lua.CheckString(l, 1),
				Position:  luaCheckPosition(l, 2),
			}
			done := app.Do(cmd)
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "DialogShow", Function: func(l *lua.State) int {
			cmd := ShowDialog{
				Text:     lua.CheckString(l, 1),
				Position: luaCheckOption(l, 2, "pos", DefaultDialogPosition, luaCheckPosition),
				Color:    luaCheckOption(l, 2, "color", White, luaCheckColor),
			}
			done := app.Do(cmd)
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "MusicPlay", Function: func(l *lua.State) int {
			done := app.Do(MusicPlay{MusicResource: ResourceLocator(lua.CheckString(l, 1))})
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "ScenePlay", Function: func(l *lua.State) int {
			done := app.Do(ScenePlay{SceneResource: ResourceLocator(lua.CheckString(l, 1))})
			luaPushFuture(l, done)
			return 1
		}},
		{Name: "SleepMillis", Function: func(l *lua.State) int {
			time.Sleep(luaCheckDurationMillis(l, 1))
			return 0
		}},
		{Name: "SoundPlay", Function: func(l *lua.State) int {
			done := app.Do(SoundPlay{SoundResource: ResourceLocator(lua.CheckString(l, 1))})
			luaPushFuture(l, done)
			return 1
		}},
	}
}

func luaDeclareConstants(l *lua.State) {
	for k, pushFunc := range map[string]func(){
		"ColorBlack":         func() { luaPushColor(l, Black) },
		"ColorBlue":          func() { luaPushColor(l, Blue) },
		"ColorGreen":         func() { luaPushColor(l, Green) },
		"ColorCyan":          func() { luaPushColor(l, Cyan) },
		"ColorRed":           func() { luaPushColor(l, Red) },
		"ColorMagenta":       func() { luaPushColor(l, Magenta) },
		"ColorBrown":         func() { luaPushColor(l, Brown) },
		"ColorLightGray":     func() { luaPushColor(l, LightGray) },
		"ColorDarkGray":      func() { luaPushColor(l, DarkGray) },
		"ColorBrigthBlue":    func() { luaPushColor(l, BrigthBlue) },
		"ColorBrigthGreen":   func() { luaPushColor(l, BrigthGreen) },
		"ColorBrigthCyan":    func() { luaPushColor(l, BrigthCyan) },
		"ColorBrigthRed":     func() { luaPushColor(l, BrigthRed) },
		"ColorBrigthMagenta": func() { luaPushColor(l, BrigthMagenta) },
		"ColorYellow":        func() { luaPushColor(l, Yellow) },
		"ColorWhite":         func() { luaPushColor(l, White) },
		"DirUp":              func() { l.PushInteger(int(DirUp)) },
		"DirRight":           func() { l.PushInteger(int(DirRight)) },
		"DirDown":            func() { l.PushInteger(int(DirDown)) },
		"DirLeft":            func() { l.PushInteger(int(DirLeft)) },
	} {
		pushFunc()
		l.SetGlobal(k)
	}
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

func luaCheckColor(l *lua.State, index int) Color {
	return Color{
		R: byte(luaCheckFieldInt(l, index, "r")),
		G: byte(luaCheckFieldInt(l, index, "g")),
		B: byte(luaCheckFieldInt(l, index, "b")),
		A: byte(luaCheckOption(l, index, "a", 255, lua.CheckInteger)),
	}
}

func luaCheckDirection(l *lua.State, index int) Direction {
	return Direction(lua.CheckInteger(l, index))
}

func luaCheckDurationMillis(l *lua.State, index int) time.Duration {
	return time.Duration(lua.CheckInteger(l, index)) * time.Millisecond
}

func luaCheckFieldInt(l *lua.State, index int, field string) int {
	l.Field(index, field)
	defer l.Pop(1)
	return lua.CheckInteger(l, -1)
}

func luaCheckOption[T any](l *lua.State, index int, field string, def T, f func(*lua.State, int) T) T {
	luaIfFieldExists(l, index, field, func() {
		def = f(l, -1)
	})
	return def
}

func luaCheckPosition(l *lua.State, index int) Position {
	return Position{
		X: luaCheckFieldInt(l, index, "x"),
		Y: luaCheckFieldInt(l, index, "y"),
	}
}

func luaCheckResourceLocator(l *lua.State, index int) ResourceLocator {
	return ResourceLocator(lua.CheckString(l, index))
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
		{Name: "IsCompleted", Function: func(l *lua.State) int {
			l.PushBoolean(f.IsCompleted())
			return 1
		}},
		{Name: "Wait", Function: func(l *lua.State) int {
			f.Wait()
			return 0
		}},
	})
	return
}
