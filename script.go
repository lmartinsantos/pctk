package pctk

import (
	"io"
	"log"
)

// ScriptLanguage represents the language of a script.
type ScriptLanguage byte

const (
	// ScriptUndefined is an undefined script language.
	ScriptUndefined ScriptLanguage = iota

	// ScriptLua is the Lua script language.
	ScriptLua
)

// Script represents a script.
type Script struct {
	Language ScriptLanguage
	Code     []byte
}

func (s *Script) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, uint16(s.Language), uint32(len(s.Code)), s.Code)
}

func (s *Script) run(app *App, prom Promise) {
	switch s.Language {
	case ScriptLua:
		s.runLua(app, prom)
	default:
		log.Panicf("Unknown script language: %0x", s.Language)
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
