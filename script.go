package pctk

import (
	"io"
	"log"
	"strings"

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

// Call a method in the script. The method is a chain of identifiers that references a function in
// the script.
func (s *Script) Call(method Method) Future {
	switch s.Language {
	case ScriptLua:
		return s.luaCall(method)
	default:
		log.Panicf("Unknown script language: %0x", s.Language)
		return nil
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

// Method is a chain of identifiers that references a function in a script.
type Method []string

// WithMethod creates a new Method with the given parts.
func WithMethod(head string, tail ...string) Method {
	return append(Method{head}, tail...)
}

// ForEach calls the given function for each part of the method.
func (m Method) ForEach(f func(string)) {
	for _, part := range m {
		f(part)
	}
}

// String returns the string representation of the method.
func (m Method) String() string {
	return strings.Join(m, ".")
}
