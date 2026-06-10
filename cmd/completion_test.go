package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShellSnippet_AllShells(t *testing.T) {
	shells := []string{"zsh", "bash", "fish", "powershell"}
	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			snippet, err := shellSnippet(shell)
			require.NoError(t, err)
			assert.NotEmpty(t, snippet)
			assert.Contains(t, snippet, "cd")
			assert.Contains(t, snippet, "new", "shell wrapper should intercept 'ent new' for auto-cd")
		})
	}
}

func TestShellSnippet_Powershell_Alias(t *testing.T) {
	snippet, err := shellSnippet("powershell")
	require.NoError(t, err)
	assert.Contains(t, snippet, "Set-Alias")
	assert.Contains(t, snippet, "Set-Location")
}

func TestShellSnippet_Zsh_BuiltinCd(t *testing.T) {
	snippet, err := shellSnippet("zsh")
	require.NoError(t, err)
	assert.Contains(t, snippet, "builtin cd")
}

func TestShellSnippet_Fish_BuiltinCd(t *testing.T) {
	snippet, err := shellSnippet("fish")
	require.NoError(t, err)
	assert.Contains(t, snippet, "builtin cd")
}

func TestShellSnippet_CaseInsensitive(t *testing.T) {
	_, err := shellSnippet("PowerShell")
	require.NoError(t, err)

	_, err = shellSnippet("ZSH")
	require.NoError(t, err)
}

func TestShellSnippet_UnknownShell(t *testing.T) {
	_, err := shellSnippet("nushell")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell")
}

func TestRunCompletion_IncludesWrapperAndCompletion(t *testing.T) {
	shells := []struct {
		name            string
		wrapperContains string
		completionMark  string
	}{
		{"zsh", "builtin cd", "compdef"},
		{"bash", "builtin cd", "__start_ent"},
		{"fish", "builtin cd", "complete"},
		{"powershell", "Set-Location", "Register-ArgumentCompleter"},
	}

	for _, tt := range shells {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			completionCmd.SetOut(buf)

			err := runCompletion(completionCmd, tt.name)
			require.NoError(t, err)

			out := buf.String()
			assert.Contains(t, out, tt.wrapperContains, "expected cd wrapper")
			assert.Contains(t, out, tt.completionMark, "expected completion script")
		})
	}
}

func TestRunCompletion_UnknownShell(t *testing.T) {
	err := runCompletion(completionCmd, "nushell")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported shell")
}
