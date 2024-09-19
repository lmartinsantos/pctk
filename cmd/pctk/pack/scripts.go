package pack

import (
	"fmt"

	"github.com/apoloval/pctk"
	"gopkg.in/yaml.v3"
)

type ScriptLanguage string

const (
	ScriptLanguageLua ScriptLanguage = "lua"
)

// ScriptData is the data for a script resource.
type ScriptData struct {
	Language pctk.ScriptLanguage
	Code     []byte
}

// AsResource converts the script data to a script resource.
func (d *ScriptData) AsResource() *pctk.Script {
	return &pctk.Script{
		Language: d.Language,
		Code:     d.Code,
	}
}

func (d *ScriptData) UnmarshalYAML(n *yaml.Node) error {
	var data struct {
		Language string
		Code     string
	}
	if err := n.Decode(&data); err != nil {
		return err
	}

	switch data.Language {
	case "", "lua", "Lua", "LUA":
		d.Language = pctk.ScriptLua
	default:
		return fmt.Errorf("unknown script language: %s", data.Language)
	}
	d.Code = []byte(data.Code)

	return nil
}
