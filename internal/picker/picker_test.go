package picker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockPicker_ReturnsSelection(t *testing.T) {
	p := &MockPicker{Selection: "myorg/myrepo"}
	got, err := p.Pick("Select project", []string{"myorg/myrepo", "other/repo"})
	require.NoError(t, err)
	assert.Equal(t, "myorg/myrepo", got)
}

func TestMockPicker_ReturnsError(t *testing.T) {
	p := &MockPicker{Err: ErrAborted}
	_, err := p.Pick("Select project", []string{"a", "b"})
	assert.ErrorIs(t, err, ErrAborted)
}
