# tuv

A terminal UI written in Go for managing Python projects created with [uv](https://github.com/astral-sh/uv). 
Create, view, and manage your Python projects with a simple interface.

## Features

- Create new Python projects with uv virtual environments
- Scan and detect existing uv projects
- View project details (size, Python version, creation date)
- Delete projects with confirmation

## Screenshots
<img width="414" alt="Screenshot 2025-06-22 at 19 29 43" src="https://github.com/user-attachments/assets/0b2caf31-0b09-49d8-8155-ac8d8ece315b" />
<img width="461" alt="Screenshot 2025-06-22 at 19 30 15" src="https://github.com/user-attachments/assets/d566b2ab-628c-4f5b-b308-1efa8d3909b5" />
<img width="501" alt="Screenshot 2025-06-22 at 19 30 27" src="https://github.com/user-attachments/assets/68b48b0e-0243-4d0a-a431-ab995c010a88" />
<img width="497" alt="Screenshot 2025-06-22 at 19 30 23" src="https://github.com/user-attachments/assets/e9c2a1f4-abec-4acb-9dc0-2ea48217f229" />

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
