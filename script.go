package pctk

import (
	"log"
)

// ScriptLanguage represents the language of a script.
type ScriptLanguage string

const (
	// ScriptLua is the Lua script language.
	ScriptLua ScriptLanguage = "lua"
)

// Script represents a script.
type Script struct {
	Language ScriptLanguage
	Code     []byte
}

func (s *Script) run(app *App, prom Promise) {
	switch s.Language {
	case ScriptLua:
		s.runLua(app, prom)
	default:
		log.Panicf("Unknown script language: %s", s.Language)
	}
}

// ScriptRun is a command to run a script.
type ScriptRun struct {
	ScriptResource ResourceLocator
}

func (c ScriptRun) Execute(app *App, prom Promise) {
	script := app.res.LoadScript(c.ScriptResource)
	if script == nil {
		log.Panicf("Script not found: %s", c.ScriptResource)
	}
	script.run(app, prom)
}
