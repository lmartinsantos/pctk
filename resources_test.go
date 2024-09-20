package pctk_test

import (
	"testing"

	"github.com/apoloval/pctk"
	"github.com/stretchr/testify/assert"
)

func TestResourceRefPackage(t *testing.T) {
	ref := pctk.NewResourceRef("001", "foo/bar")
	assert.Equal(t, "001", ref.Package())
}
