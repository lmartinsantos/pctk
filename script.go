package pctk

import (
	"io"
	"log"
	"strings"
	"sync"

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

	mutex     sync.Mutex
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

// Call a  script function. The script must be run before calling this method.
func (s *Script) Call(f FieldAccessor, args []any, method bool) Future {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch s.Language {
	case ScriptLua:
		return s.luaCall(f, args, method)
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

// Run the script. This will evaluate the code in the script, running the declarations (if any) and
// preparing the code to receive calls.
func (s *Script) Run(app *App, prom *Promise) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	switch s.Language {
	case ScriptLua:
		s.luaInit(app)
		s.luaRun(app, prom)
	default:
		log.Panicf("Unknown script language: %0x", s.Language)
	}
}

// FieldAccessor is a sequence of identifiers that references a field in a chain of tables.
type FieldAccessor []string

// WithField creates a new FieldAccessor with the given parts.
func WithField(global string, fields ...string) FieldAccessor {
	return append(FieldAccessor{global}, fields...)
}

// WithActorField creates a new FieldAccessor pointing to an actor value.
func WithActorField(actor *Actor, fields ...string) FieldAccessor {
	return WithField(actor.id, fields...)
}

// WithObjectField creates a new FieldAccessor pointing to an object value.
func WithObjectField(obj *Object, fields ...string) FieldAccessor {
	return WithField(obj.id, append([]string{"objects"}, fields...)...)
}

// WithDefaultsField creates a new FieldAccessor pointing to the defaults object.
func WithDefaultsField(fields ...string) FieldAccessor {
	return WithField("DEFAULT", fields...)
}

// Append appends the given fields to the accessor.
func (m FieldAccessor) Append(fields ...string) FieldAccessor {
	return append(m, fields...)
}

// ForEach calls the given function for each element of the accessor.
func (m FieldAccessor) ForEach(f func(string)) {
	for _, part := range m {
		f(part)
	}
}

// Base returns the base accessor of the accessor. This is the accessor without the last element.
// If the accessor points to a global variable, it returns itself.
func (m FieldAccessor) Base() FieldAccessor {
	if len(m) == 1 {
		return m
	}
	return m[:len(m)-1]
}

// String returns the string representation of the fields accessor.
func (m FieldAccessor) String() string {
	return strings.Join(m, ".")
}
