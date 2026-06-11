# ent

A fast, opinionated project manager for your terminal. `ent` keeps your projects organised under a single root directory and lets you navigate, create, and clone them without leaving the command line.

## Features

- **Navigate projects** with a fuzzy picker (no external `fzf` dependency)
- **Create local projects** with `git init` in one command
- **Clone remote repos** with automatic `org/project` directory layout
- **Delete projects** safely with a confirmation prompt
- **Shell integration** for zsh, bash, fish, and PowerShell — `ent cd` and `ent new` actually change your directory
- **Configurable** via `~/.config/ent/settings.json`, env var, or `--projects-root` flag
- **Single static binary**, no runtime dependencies

## Installation

### Homebrew (macOS / Linux)

```sh
brew install AWDDude/tap/ent
```

### Download

Grab the latest binary from the [GitHub Releases](https://github.com/AWDDude/ent/releases) page.

### Build from source

```sh
git clone https://github.com/AWDDude/ent
cd ent
make build       # current platform
make install     # install to $GOPATH/bin
make build-all   # all platforms → dist/
```

## Quick Start

### 1. Add shell integration

One line sets up the shell wrapper (enables `ent cd` and `ent new` to change your directory) **and** tab completion:

```sh
# zsh — add to ~/.zshrc
eval "$(ent completion zsh)"

# bash — add to ~/.bashrc
eval "$(ent completion bash)"

# fish — add to ~/.config/fish/config.fish
ent completion fish | source

# PowerShell — add to $PROFILE
ent completion powershell | Out-String | Invoke-Expression
```

### 2. Configure your projects root (optional)

By default `ent` uses `~/projects`. On first run, `ent` automatically creates `~/.config/ent/settings.json` with that default. Edit it to point to your preferred location:

```json
{
  "projects_root": "/Users/you/code"
}
```

Or use the environment variable:

```sh
export PROJECTS_ROOT=/Users/you/code
```

## Commands

| Command | Description |
|---------|-------------|
| `ent list` / `ent ls` | List all projects as `org/project` |
| `ent cd [name]` | Navigate to a project (fuzzy picker if no name given) |
| `ent new <name>` | Create `$projects_root/local/<name>` + `git init`, then cd into it |
| `ent clone <url>` | Clone a repo into `$projects_root/<org>/<project>` |
| `ent rm [name]` | Permanently delete a project (fuzzy picker if no name given) |
| `ent completion <shell>` | Shell cd wrapper + tab completion |
| `ent version` | Print version, commit, and build date |

### `ent cd`

Without a name, opens the fuzzy picker over all projects:

```sh
ent cd
```

With a partial name, navigates directly (or picks if multiple match):

```sh
ent cd myapp       # matches project named *myapp*
ent cd myorg/app   # matches org/project substring
```

### `ent new`

Creates a local project at `$projects_root/local/<name>` and navigates into it (requires shell integration):

```sh
ent new my-new-service
# Created project: /home/you/source/local/my-new-service
# (shell cd's into the new directory automatically)
```

### `ent clone`

Clones a remote repository into the right org subdirectory:

```sh
ent clone git@github.com:myorg/myrepo.git
# → /home/you/source/myorg/myrepo

ent clone https://github.com/myorg/myrepo.git
# → /home/you/source/myorg/myrepo
```

### `ent rm`

Permanently deletes a project directory. Opens the fuzzy picker if no name is given or multiple projects match. Requires typing the full `org/project` name to confirm before anything is deleted. The parent org directory is removed automatically if it becomes empty.

```sh
ent rm myapp         # delete project named *myapp*
ent rm               # fuzzy pick from all projects
```

```
WARNING: This action is irreversible.
Any changes that have not been pushed to a git server will be lost.

Project to delete: myorg/myapp

Type the project name to confirm: myorg/myapp
Deleted: /home/you/source/myorg/myapp
```

## Project Layout

```
$projects_root/
├── local/           # projects created with 'ent new'
│   └── my-app/
├── myorg/           # cloned from github.com/myorg/...
│   ├── backend/
│   └── frontend/
└── other-org/
    └── tool/
```

## Configuration

| Priority | Source | Example |
|----------|--------|---------|
| 1 (highest) | `--projects-root` flag | `ent --projects-root /tmp/work list` |
| 2 | `PROJECTS_ROOT` env var | `export PROJECTS_ROOT=/code` |
| 3 | `~/.config/ent/settings.json` | `{ "projects_root": "/code" }` |
| 4 (default) | `~/projects` | — |

### Settings file location

`~/.config/ent/settings.json` on all platforms, or `$XDG_CONFIG_HOME/ent/settings.json` if that env var is set.

Set `ENT_CONFIG_PATH` to point `ent` at a specific file, overriding all of the above:

```sh
export ENT_CONFIG_PATH=/work/shared/ent-config.json
```

## Contributing

See [CLAUDE.md](CLAUDE.md) for development guidance.

## License

[MIT](LICENSE)
