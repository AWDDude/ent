package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionCmd_Output(t *testing.T) {
	buf := &bytes.Buffer{}
	versionCmd.SetOut(buf)

	versionCmd.Run(versionCmd, nil)

	out := buf.String()
	require.NotEmpty(t, out)
	assert.True(t, strings.HasPrefix(out, "ent "))
	assert.Contains(t, out, "commit:")
	assert.Contains(t, out, "built:")
}

func TestVersionCmd_InjectsVariables(t *testing.T) {
	origVersion := Version
	origCommit := Commit
	origDate := BuildDate
	defer func() {
		Version = origVersion
		Commit = origCommit
		BuildDate = origDate
	}()

	Version = "v1.2.3"
	Commit = "abc123"
	BuildDate = "2026-01-01"

	buf := &bytes.Buffer{}
	versionCmd.SetOut(buf)
	versionCmd.Run(versionCmd, nil)

	out := buf.String()
	assert.Contains(t, out, "v1.2.3")
	assert.Contains(t, out, "abc123")
	assert.Contains(t, out, "2026-01-01")
}
