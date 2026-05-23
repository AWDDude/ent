# ent — Claude Code Instructions

## Project Overview

`ent` is a Go CLI project manager. It organises projects under a root directory as `$root/<org>/<project>` and provides commands to list, navigate, create, and clone them.

**Module:** `github.com/AWDDude/ent`  
**Binary:** `ent`  
**Go version:** see `go.mod`

## Commands

```sh
go build ./...                    # compile
go test ./...                     # run all tests
go test -cover ./...              # with coverage
go test -run TestXxx ./pkg/...    # run specific test
goreleaser build --snapshot --clean  # local snapshot build (no publish)
goreleaser release --snapshot --clean  # full local release dry-run
```

## Architecture

```
ent/
├── cmd/               # Cobra commands (root, list, cd, new, clone, init, version)
├── internal/
│   ├── config/        # Settings loading: settings.json → env → flag
│   ├── projects/      # Core logic: List, Resolve, New, Clone
│   ├── urlparse/      # Git URL parser (SSH + HTTPS)
│   └── picker/        # go-fuzzyfinder wrapper + MockPicker
├── main.go            # Bootstraps cobra only; excluded from test coverage
└── .goreleaser.yaml   # Cross-platform release config
```

## Key Design Decisions

### Shell `cd` integration
A Go binary cannot change its parent shell's directory. `ent cd` **prints the resolved path** to stdout. A shell wrapper function (output of `ent init <shell>`) intercepts `ent cd` and calls `builtin cd` on the result. Users install this with `eval "$(ent init zsh)"`.

`ent init` also appends the Cobra-generated tab-completion script in the same output, so one line in `.zshrc` handles both `cd` integration and tab completion. This is implemented in `cmd/init.go` via `genCompletion()` which calls `root.GenZshCompletion()` (or the equivalent for other shells).

### Config precedence
`--projects-root` flag > `PROJECTS_ROOT` env var > `settings.json` > default `~/projects`

### CGo disabled
All builds use `CGO_ENABLED=0` for portable static binaries.

### Fuzzy picker
Uses `github.com/ktr0731/go-fuzzyfinder` (pure Go, no external `fzf` binary needed). The `picker.Picker` interface allows commands to be tested with `picker.MockPicker`.

### Git operations
Git is invoked via `exec.Command` (subprocess). The `GitRunner` function type on `projects.Manager` is injectable for testing (see `noopGit` in `manager_test.go`).

## Testing Guidelines

- **Full coverage is the goal** for all packages except `main.go` and the real I/O paths (`defaultGitRunner`, `FuzzyPicker.Pick`).
- Commands are tested by calling `runXxx()` directly (not through `cobra.Command.Execute()`) to control config via `t.Setenv`.
- Always set `PROJECTS_ROOT` and `XDG_CONFIG_HOME` in tests to avoid touching real user config.
- Use `picker.MockPicker` when testing commands that call the picker.
- Use `noopGit` (defined in `projects/manager_test.go`) or a custom `GitRunner` when testing project creation/cloning.

## GoReleaser

Releases are triggered by pushing a semver tag:

```sh
git tag v1.2.3
git push origin v1.2.3
```

Required secrets in GitHub Actions:
- `GITHUB_TOKEN` — standard, auto-provided
- `HOMEBREW_TAP_TOKEN` — PAT with write access to `AWDDude/homebrew-tap`

The Homebrew formula is published to `AWDDude/homebrew-tap` under `Formula/ent.rb`.

## Adding a New Command

1. Create `cmd/<name>.go` with a `var <name>Cmd = &cobra.Command{...}` and a `run<Name>(cmd, args)` function.
2. Register it in `cmd/root.go` `init()`: `rootCmd.AddCommand(<name>Cmd)`.
3. Create `cmd/<name>_test.go` testing `run<Name>` directly.
4. Update `README.md` commands table.

## Adding a New Settings Field

1. Add the field to `config.Settings` in `internal/config/config.go`.
2. Apply it in `config.Load()`.
3. Add a test case in `internal/config/config_test.go`.
4. Update `README.md` configuration section.
