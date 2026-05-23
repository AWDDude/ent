package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// completionCmd replaces Cobra's default completion command, combining the
// shell cd-wrapper with the tab-completion script in a single output.
var completionCmd = &cobra.Command{
	Use:       "completion <shell>",
	Short:     "Generate shell integration (cd wrapper + tab completion)",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"zsh", "bash", "fish", "powershell"},
	Long: `Generate a shell integration script that provides both:

  1. A wrapper function so 'ent cd' actually changes your directory
  2. Tab completion for all ent commands and flags

Supported shells: zsh, bash, fish, powershell

Usage:
  # zsh — add to ~/.zshrc
  eval "$(ent completion zsh)"

  # bash — add to ~/.bashrc
  eval "$(ent completion bash)"

  # fish — add to ~/.config/fish/config.fish
  ent completion fish | source

  # PowerShell — add to $PROFILE
  ent completion powershell | Out-String | Invoke-Expression`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCompletion(cmd, args[0])
	},
}

func runCompletion(cmd *cobra.Command, shell string) error {
	snippet, err := shellSnippet(shell)
	if err != nil {
		return err
	}
	fmt.Fprint(cmd.OutOrStdout(), snippet)
	return genCompletion(cmd, shell)
}

// genCompletion writes the Cobra-generated tab-completion script for shell.
func genCompletion(cmd *cobra.Command, shell string) error {
	root := cmd.Root()
	w := cmd.OutOrStdout()
	switch strings.ToLower(shell) {
	case "zsh":
		return root.GenZshCompletion(w)
	case "bash":
		return root.GenBashCompletion(w)
	case "fish":
		return root.GenFishCompletion(w, true)
	case "powershell", "pwsh":
		return root.GenPowerShellCompletionWithDesc(w)
	}
	return nil
}

func shellSnippet(shell string) (string, error) {
	switch strings.ToLower(shell) {
	case "zsh":
		return zshSnippet, nil
	case "bash":
		return bashSnippet, nil
	case "fish":
		return fishSnippet, nil
	case "powershell", "pwsh":
		return powershellSnippet, nil
	default:
		return "", fmt.Errorf("unsupported shell %q; supported: zsh, bash, fish, powershell", shell)
	}
}

const zshSnippet = `# ent shell integration — added by 'ent completion zsh'
ent() {
  if [[ "$1" == "cd" ]]; then
    local dir _err
    dir=$(command ent cd "${@:2}" 2>/tmp/.ent-err)
    local rc=$?
    _err=$(cat /tmp/.ent-err 2>/dev/null)
    [[ -n "$_err" ]] && echo "$_err" >&2
    [[ $rc -eq 0 && -n "$dir" ]] && builtin cd "$dir"
  else
    command ent "$@"
  fi
}
`

const bashSnippet = `# ent shell integration — added by 'ent completion bash'
ent() {
  if [[ "$1" == "cd" ]]; then
    local dir _err
    dir=$(command ent cd "${@:2}" 2>/tmp/.ent-err)
    local rc=$?
    _err=$(cat /tmp/.ent-err 2>/dev/null)
    [[ -n "$_err" ]] && echo "$_err" >&2
    [[ $rc -eq 0 && -n "$dir" ]] && builtin cd "$dir"
  else
    command ent "$@"
  fi
}
`

const fishSnippet = `# ent shell integration — added by 'ent completion fish'
function ent
    if test "$argv[1]" = "cd"
        set dir (command ent cd $argv[2..] 2>/tmp/.ent-err)
        set rc $status
        set _err (cat /tmp/.ent-err 2>/dev/null)
        test -n "$_err" && echo $_err >&2
        if test $rc -eq 0 -a -n "$dir"
            builtin cd $dir
        end
    else
        command ent $argv
    end
end
`

const powershellSnippet = `# ent shell integration — added by 'ent completion powershell'
function Invoke-Ent {
    if ($args[0] -eq 'cd') {
        $errFile = [System.IO.Path]::GetTempFileName()
        $dir = & ent @args 2>$errFile
        $rc = $LASTEXITCODE
        $errMsg = Get-Content $errFile -ErrorAction SilentlyContinue
        Remove-Item $errFile -ErrorAction SilentlyContinue
        if ($errMsg) { Write-Error $errMsg }
        if ($rc -eq 0 -and $dir) { Set-Location $dir }
    } else {
        & ent @args
    }
}
Set-Alias -Name ent -Value Invoke-Ent -Option AllScope -Force
`
