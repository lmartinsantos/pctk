package pack

import (
	"testing"

	"github.com/apoloval/pctk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestLoadManifest(t *testing.T) {
	text := `
type: script
compression: none
data:
  language: lua
  code: |
    print("Hello, world!")
`
	var man Manifest
	err := yaml.Unmarshal([]byte(text), &man)

	require.NoError(t, err)
	assert.Equal(t, ManifestTypeScript, man.Type)
	assert.Equal(t, pctk.CompressionNone, man.Compression)

	script, ok := man.Data.(*ScriptData)
	require.True(t, ok)
	assert.Equal(t, "lua", script.Language)
	assert.Equal(t, "print(\"Hello, world!\")\n", script.Code)
}
