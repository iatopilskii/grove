# grove - Git Worktree TUI

A terminal user interface for managing Git worktrees, built with Go and [BubbleTea](https://github.com/charmbracelet/bubbletea).

## Features

- List, create, delete, and prune worktrees
- Two-pane layout with worktree details (path, branch, status)
- Keyboard and mouse navigation
- Adaptive light/dark color scheme
- User-configurable themes via YAML

## Installation

```bash
go install github.com/iatopilskii/grove/cmd/grove@latest
```

Or build from source:

```bash
git clone https://github.com/iatopilskii/grove.git
cd grove
go build ./cmd/grove
```

## Usage

Run `grove` inside any Git repository:

```bash
grove
```

### Shell Wrapper (Recommended)

To automatically cd into newly created worktrees, add this wrapper to your shell rc file:

**Bash/Zsh** (`~/.bashrc` or `~/.zshrc`):

```bash
grove() {
    local target
    target=$(command grove "$@")
    local ec=$?
    [[ $ec -eq 2 && -d "$target" ]] && cd "$target"
}
```

**Fish** (`~/.config/fish/functions/grove.fish`):

```fish
function grove
    set -l target (command grove $argv)
    set -l ec $status
    if test $ec -eq 2 -a -d "$target"
        cd "$target"
    end
end
```

### Keybindings

| Key                   | Action                |
| --------------------- | --------------------- |
| `Tab` / `Shift+Tab`   | Switch tabs           |
| `↑` / `↓` / `j` / `k` | Navigate list         |
| `PgUp` / `PgDn`       | Page navigation       |
| `Enter`               | Open action menu      |
| `n`                   | Create new worktree   |
| `p`                   | Prune stale worktrees |
| `Esc`                 | Close dialog          |
| `q` / `Ctrl+C`        | Quit                  |

Mouse clicks and scroll are also supported.

## Configuration

Config file location: `~/.config/grove/config.yaml`

Example:

```yaml
theme:
  colors:
    primary:
      light: "#7c3aed"
      dark: "#a78bfa"
    text:
      light: "#1f2937"
      dark: "#f9fafb"
```

## Requirements

- Go 1.24+
- Git

## License

MIT
