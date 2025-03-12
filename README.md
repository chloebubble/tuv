# tuv

A terminal UI written in Go for managing Python projects created with [uv](https://github.com/astral-sh/uv). 
Create, view, and manage your Python projects with a simple interface.

## Features

- Create new Python projects with uv virtual environments
- Scan and detect existing uv projects
- View project details (size, Python version, creation date)
- Delete projects with confirmation

## Installation

Requires Go 1.21+ and [uv](https://github.com/astral-sh/uv).

```bash
git clone git@github.com:chloebubble/tuv.git
cd tuv
go mod download
go install github.com/chloebubble/tuv
```

## Usage

Run `tuv` in your terminal. On first run, you'll be prompted to configure your projects directory.

Navigation:
- ↑/↓ or j/k to navigate
- Enter to select
- s to rescan for projects
- d to delete a project (with confirmation)
- Esc to go back
- q or Ctrl+C to quit

## Config

Settings are stored in `~/.config/tuv/config.yaml`.

## Acknowledgments

Built with:
- [uv](https://github.com/astral-sh/uv)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
