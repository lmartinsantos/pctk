package pctk_test

import (
	"testing"

	"github.com/apoloval/pctk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceRefParse(t *testing.T) {
	ref, err := pctk.ParseResourceRef("pkg:foo/bar")
	require.NoError(t, err)

	assert.Equal(t, pctk.ResourcePackage("pkg"), ref.Package())
	assert.Equal(t, pctk.ResourceID("foo/bar"), ref.ID())
}
