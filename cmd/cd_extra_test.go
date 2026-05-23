package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimPrefix_Present(t *testing.T) {
	assert.Equal(t, "myorg/myrepo", trimPrefix("/root/myorg/myrepo", "/root/"))
}

func TestTrimPrefix_Absent(t *testing.T) {
	assert.Equal(t, "/other/path", trimPrefix("/other/path", "/root/"))
}

func TestTrimPrefix_ExactMatch(t *testing.T) {
	assert.Equal(t, "", trimPrefix("/root/", "/root/"))
}
