package urlparse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_Whitespace(t *testing.T) {
	got, err := Parse("  git@github.com:myorg/myrepo.git  ")
	require.NoError(t, err)
	assert.Equal(t, "myorg", got.Org)
	assert.Equal(t, "myrepo", got.Project)
}

func TestParse_HTTPSEmptyOrgOrProject(t *testing.T) {
	// Path with empty segment after trimming
	_, err := Parse("https://github.com//project")
	require.Error(t, err)
}
