package pctk

import (
	"io"
	"log"

	"github.com/Shopify/go-lua"
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

	ref       ResourceRef
	l         *lua.State
	including bool
}

// NewScript creates a new script.
func NewScript(lang ScriptLanguage, code []byte) *Script {
	return &Script{
		Language: lang,
		Code:     code,
	}
}

// BinaryDecode decodes the script from a binary stream. The format is:
//   - byte: the script language.
//   - uint32: the length of the script code.
//   - []byte: the script code.
func (s *Script) BinaryEncode(w io.Writer) (int, error) {
	return BinaryEncode(w, s.Language, uint32(len(s.Code)), s.Code)
}

// BinaryDecode decodes the script from a binary stream. See Script.BinaryEncode for the format.
func (s *Script) BinaryDecode(r io.Reader) error {
	var lang ScriptLanguage
	var length uint32
	if err := BinaryDecode(r, &lang, &length); err != nil {
		return err
	}

	code := make([]byte, length)
	if err := BinaryDecode(r, &code); err != nil {
		return err
	}

	s.Language = ScriptLanguage(lang)
	s.Code = code
	return nil
}

func (s *Script) init(app *App, ref ResourceRef) {
	s.ref = ref
	switch s.Language {
	case ScriptLua:
		s.luaInit(app)
	default:
		log.Panicf("Unknown script language: %0x", s.Language)
	}
}

func (s *Script) run(app *App, prom *Promise) {
	switch s.Language {
	case ScriptLua:
		s.luaRun(app, prom)
	default:
		log.Panicf("Unknown script language: %0x", s.Language)
	}
}

func (s *Script) call(chain ...string) Future {
	switch s.Language {
	case ScriptLua:
		return s.luaCall(chain...)
	default:
		log.Panicf("Unknown script language: %0x", s.Language)
		return nil
	}
}

// ScriptRun is a command to run a script.
type ScriptRun struct {
	ScriptRef ResourceRef
}

func (c ScriptRun) Execute(app *App, prom *Promise) {
	script, ok := app.scripts[c.ScriptRef]
	if ok {
		prom.CompleteWithValue(script)
		return
	}

	// The script was not loaded yet, so we load and execute it now.
	script = app.res.LoadScript(c.ScriptRef)
	if script == nil {
		log.Panicf("Script not found: %s", c.ScriptRef)
	}
	script.init(app, c.ScriptRef)
	app.scripts[c.ScriptRef] = script
	script.run(app, prom)
}

// ScriptCall is a command to call a script function.
type ScriptCall struct {
	ScriptRef ResourceRef
	Method    []string
}

func (c ScriptCall) Execute(app *App, prom *Promise) {
	script, ok := app.scripts[c.ScriptRef]
	if !ok {
		log.Panicf("Script not found: %s", c.ScriptRef)
	}

	prom.CompleteWhen(script.call(c.Method...))
}
