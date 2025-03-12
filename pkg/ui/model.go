package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crung/tuv/pkg/config"
	"github.com/crung/tuv/pkg/scanner"
)

// AppState represents the current state of the application
type AppState int

const (
	StateMainMenu AppState = iota
	StateFirstRun
	StateProjectList
	StateProjectDetail
	StateNewProject
	StateLoading
)

// KeyMap defines the keybindings for the application
type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Back   key.Binding
	Quit   key.Binding
	Scan   key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "quit"),
		),
		Scan: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "scan"),
		),
	}
}

// Model represents the main application model
type Model struct {
	config          *config.Config
	scanner         *scanner.Scanner
	keyMap          KeyMap
	state           AppState
	width           int
	height          int
	menuItems       []string
	selectedMenu    int
	projects        []scanner.UVProject
	selectedProject int
	textInput       textinput.Model
	spinner         spinner.Model
	loading         bool
	loadingMsg      string
	statusMsg       string
	error           string
}

// NewModel creates a new application model
func NewModel(cfg *config.Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	ti := textinput.New()
	ti.Placeholder = "Enter path"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	// Set default value to the current parent directory
	if cfg.ParentDirectory != "" {
		ti.SetValue(cfg.ParentDirectory)
	} else {
		// Set to home directory as default
		if homeDir, err := os.UserHomeDir(); err == nil {
			ti.SetValue(filepath.Join(homeDir, "projects"))
		}
	}

	menuItems := []string{
		"List projects",
		"New project",
		"Quit",
	}

	// Check if this is the first run using the IsFirstRun method
	var initialState AppState
	if cfg.IsFirstRun() {
		initialState = StateFirstRun
	} else {
		initialState = StateMainMenu
	}

	scn := scanner.NewScanner(cfg.ParentDirectory)

	m := Model{
		config:          cfg,
		scanner:         scn,
		keyMap:          DefaultKeyMap(),
		state:           initialState,
		menuItems:       menuItems,
		selectedMenu:    0,
		selectedProject: 0,
		textInput:       ti,
		spinner:         s,
		loading:         false,
	}

	// Always scan on startup if not in first run state
	if initialState != StateFirstRun {
		m.loading = true
		m.loadingMsg = "Scanning for uv projects..."
	}

	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
	}

	if m.loading {
		cmds = append(cmds, m.scanProjects)
	}

	return tea.Batch(cmds...)
}

// scanProjects scans for uv projects
func (m Model) scanProjects() tea.Msg {
	projects, err := m.scanner.ScanProjects()
	if err != nil {
		return errMsg{err}
	}
	return projectsFoundMsg{projects}
}

// Update handles updates to the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
		}

		switch m.state {
		case StateMainMenu:
			return m.updateMainMenu(msg)
		case StateFirstRun:
			return m.updateFirstRun(msg)
		case StateProjectList:
			return m.updateProjectList(msg)
		case StateProjectDetail:
			return m.updateProjectDetail(msg)
		case StateNewProject:
			return m.updateNewProject(msg)
		case StateLoading:
			// If we're in the loading state, just return
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case projectsFoundMsg:
		m.projects = msg.projects
		m.loading = false
		m.statusMsg = fmt.Sprintf("Found %d uv projects", len(m.projects))

		// If no projects were found and we're not in the main menu, go to main menu
		if len(m.projects) == 0 && m.state != StateMainMenu {
			m.state = StateMainMenu
		}

	case projectCreatedMsg:
		m.projects = msg.projects
		m.loading = false
		m.state = StateMainMenu
		m.statusMsg = fmt.Sprintf("Project '%s' created successfully!", msg.projectName)
		// Reset text input for next use
		m.textInput.SetValue("")

	case statusMsg:
		m.statusMsg = msg.msg
		m.loading = false

	case errMsg:
		m.error = msg.err.Error()
		m.loading = false
	}

	return m, tea.Batch(cmds...)
}

// updateMainMenu handles updates in the main menu state
func (m Model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			if m.selectedMenu > 0 {
				m.selectedMenu--
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Down):
			if m.selectedMenu < len(m.menuItems)-1 {
				m.selectedMenu++
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Select):
			switch m.selectedMenu {
			case 0: // List projects
				if len(m.projects) == 0 {
					m.loading = true
					m.loadingMsg = "Scanning for uv projects..."
					return m, m.scanProjects
				}
				m.state = StateProjectList
				return m, nil

			case 1: // New project
				m.state = StateNewProject
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, nil

			case 2: // Quit
				return m, tea.Quit
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			// Rescan for projects when 'r' is pressed
			m.loading = true
			m.loadingMsg = "Scanning for uv projects..."
			return m, m.scanProjects
		}
	}

	return m, nil
}

// updateFirstRun handles updates in the first run state
func (m Model) updateFirstRun(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Select):
			if m.textInput.Value() != "" {
				// Expand tilde to home directory if present
				inputPath := m.textInput.Value()
				if strings.HasPrefix(inputPath, "~") {
					homeDir, err := os.UserHomeDir()
					if err != nil {
						m.error = "Could not determine home directory: " + err.Error() + "\nUsing default directory instead."
						// Fall back to default directory
						return m.useDefaultDirectory()
					}
					inputPath = filepath.Join(homeDir, inputPath[1:])
				}

				// Check if directory exists
				fileInfo, err := os.Stat(inputPath)
				if os.IsNotExist(err) {
					// Directory doesn't exist, try to create it
					if err := os.MkdirAll(inputPath, 0755); err != nil {
						m.error = "Could not create directory: " + err.Error() + "\nUsing default directory instead."
						// Fall back to default directory
						return m.useDefaultDirectory()
					}

					// Verify the directory was created successfully
					fileInfo, err = os.Stat(inputPath)
					if err != nil {
						m.error = "Failed to verify directory creation: " + err.Error() + "\nUsing default directory instead."
						// Fall back to default directory
						return m.useDefaultDirectory()
					}
				} else if err != nil {
					// Some other error occurred
					m.error = "Error accessing directory: " + err.Error() + "\nUsing default directory instead."
					// Fall back to default directory
					return m.useDefaultDirectory()
				}

				// Check if it's actually a directory
				if !fileInfo.IsDir() {
					m.error = "The specified path is not a directory.\nUsing default directory instead."
					// Fall back to default directory
					return m.useDefaultDirectory()
				}

				// Check if the directory is writable by trying to create a temporary file
				testFile := filepath.Join(inputPath, ".tuv_write_test")
				if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
					m.error = "Directory is not writable: " + err.Error() + "\nUsing default directory instead."
					// Fall back to default directory
					return m.useDefaultDirectory()
				}
				// Clean up the test file
				os.Remove(testFile)

				m.config.ParentDirectory = inputPath
				err = m.config.Save()
				if err != nil {
					m.error = "Could not save configuration: " + err.Error() + "\nUsing default directory instead."
					// Fall back to default directory
					return m.useDefaultDirectory()
				}

				m.scanner = scanner.NewScanner(m.config.ParentDirectory)
				m.state = StateMainMenu
				m.loading = true
				m.loadingMsg = "Scanning for uv projects..."
				return m, m.scanProjects
			}
		}
	}

	return m, cmd
}

// useDefaultDirectory sets up the default directory and returns the model and command
func (m Model) useDefaultDirectory() (tea.Model, tea.Cmd) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		m.error = "Could not determine home directory: " + err.Error()
		return m, nil
	}

	defaultDir := filepath.Join(homeDir, "projects")

	// Try to create the default directory
	if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
		if err := os.MkdirAll(defaultDir, 0755); err != nil {
			m.error = "Could not create default directory: " + err.Error()
			return m, nil
		}
	}

	m.config.ParentDirectory = defaultDir
	err = m.config.Save()
	if err != nil {
		m.error = "Could not save configuration: " + err.Error()
		return m, nil
	}

	m.scanner = scanner.NewScanner(m.config.ParentDirectory)
	m.state = StateMainMenu
	m.loading = true
	m.loadingMsg = "Scanning for uv projects..."
	return m, m.scanProjects
}

// updateProjectList handles updates in the project list state
func (m Model) updateProjectList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Up):
			if m.selectedProject > 0 {
				m.selectedProject--
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Down):
			if m.selectedProject < len(m.projects)-1 {
				m.selectedProject++
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Select):
			if len(m.projects) > 0 {
				m.state = StateProjectDetail
			}
			return m, nil

		case key.Matches(msg, m.keyMap.Back):
			m.state = StateMainMenu
			return m, nil

		case key.Matches(msg, m.keyMap.Scan):
			m.loading = true
			m.loadingMsg = "Scanning for uv projects..."
			return m, m.scanProjects
		}
	}

	return m, nil
}

// updateProjectDetail handles updates in the project detail state
func (m Model) updateProjectDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Back):
			m.state = StateProjectList
			return m, nil
		}
	}

	return m, nil
}

// updateNewProject handles updates in the new project state
func (m Model) updateNewProject(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Back):
			m.state = StateMainMenu
			return m, nil

		case key.Matches(msg, m.keyMap.Select):
			if m.textInput.Value() != "" {
				projectName := m.textInput.Value()
				projectPath := filepath.Join(m.config.ParentDirectory, projectName)

				// Check if project already exists
				if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
					m.error = fmt.Sprintf("Project %s already exists", projectName)
					return m, nil
				}

				m.loading = true
				m.loadingMsg = fmt.Sprintf("Creating new project: %s", projectName)

				return m, func() tea.Msg {
					// Create project directory
					if err := os.MkdirAll(projectPath, 0755); err != nil {
						return errMsg{err}
					}

					// Create virtual environment
					if _, err := scanner.RunUVCommand(projectPath, "venv"); err != nil {
						return errMsg{err}
					}

					// Create a basic pyproject.toml file
					pyprojectContent := fmt.Sprintf(`[project]
name = "%s"
version = "0.1.0"
description = "A new Python project"
readme = "README.md"
requires-python = ">=3.8"
license = {text = "MIT"}

[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"
`, projectName)

					if err := os.WriteFile(filepath.Join(projectPath, "pyproject.toml"), []byte(pyprojectContent), 0644); err != nil {
						return errMsg{err}
					}

					// Create a .python-version file
					// Get Python version from the virtual environment
					pythonVersion := "3.12" // Default version if we can't detect it
					if output, err := scanner.RunUVCommand(projectPath, "venv", "python", "--version"); err == nil {
						// Parse output like "Python 3.12.0"
						parts := strings.Split(strings.TrimSpace(output), " ")
						if len(parts) >= 2 {
							pythonVersion = parts[1]
						}
					}

					if err := os.WriteFile(filepath.Join(projectPath, ".python-version"), []byte(pythonVersion), 0644); err != nil {
						return errMsg{err}
					}

					// Create a basic README.md file
					readmeContent := fmt.Sprintf("# %s\n\nA new Python project created with TUV.\n", projectName)
					if err := os.WriteFile(filepath.Join(projectPath, "README.md"), []byte(readmeContent), 0644); err != nil {
						return errMsg{err}
					}

					// Create a basic Python file
					helloContent := `def main():
    print("Hello, world!")

if __name__ == "__main__":
    main()
`
					if err := os.WriteFile(filepath.Join(projectPath, "hello.py"), []byte(helloContent), 0644); err != nil {
						return errMsg{err}
					}

					// Scan for projects to update the list
					projects, err := m.scanner.ScanProjects()
					if err != nil {
						return errMsg{err}
					}

					// Return to main menu with success message
					return projectCreatedMsg{
						projects:    projects,
						projectName: projectName,
					}
				}
			}
		}
	}

	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	switch m.state {
	case StateMainMenu:
		return m.viewMainMenu()
	case StateFirstRun:
		return m.viewFirstRun()
	case StateProjectList:
		return m.viewProjectList()
	case StateProjectDetail:
		return m.viewProjectDetail()
	case StateNewProject:
		return m.viewNewProject()
	case StateLoading:
		return m.viewLoading()
	default:
		return "Unknown state"
	}
}

// viewMainMenu renders the main menu
func (m Model) viewMainMenu() string {
	var b strings.Builder

	// Add the ASCII art logo
	b.WriteString(GetLogo() + "\n")

	// Add a version badge and tagline with proper alignment
	version := VersionBadgeStyle.Render("v0.1.0")
	tagline := SubtitleStyle.Render("Terminal UV Environment Manager")
	headerLine := lipgloss.JoinHorizontal(lipgloss.Center, version, tagline)
	b.WriteString(headerLine + "\n\n")

	var menuRows []string
	for i, item := range m.menuItems {
		if i == m.selectedMenu {
			menuRows = append(menuRows, SelectedItemStyle.Render(fmt.Sprintf(" > %s", item)))
		} else {
			menuRows = append(menuRows, ItemStyle.Render(fmt.Sprintf("   %s", item)))
		}
	}

	// Join rows with newlines to ensure vertical layout
	menuContent := strings.Join(menuRows, "\n")
	menu := MenuStyle.Render(menuContent)
	b.WriteString(menu + "\n\n")

	if m.loading {
		b.WriteString(StatusStyle.Render(fmt.Sprintf("%s %s", m.spinner.View(), m.loadingMsg)) + "\n")
	} else if m.statusMsg != "" {
		b.WriteString(StatusStyle.Render(m.statusMsg) + "\n")
	}

	if m.error != "" {
		b.WriteString(ErrorStyle.Render("Error: "+m.error) + "\n")
	}

	help := HelpStyle.Render("↑/↓: Navigate • Enter: Select • s: Rescan • q: Quit")
	b.WriteString("\n" + help)

	return BaseStyle.Render(b.String())
}

// viewFirstRun renders the first run configuration screen
func (m Model) viewFirstRun() string {
	var b strings.Builder

	// Add the ASCII art logo
	b.WriteString(GetLogo() + "\n")

	// Add a version badge and tagline with proper alignment
	version := VersionBadgeStyle.Render("v0.1.0")
	tagline := SubtitleStyle.Render("First Time Setup")
	headerLine := lipgloss.JoinHorizontal(lipgloss.Center, version, tagline)
	b.WriteString(headerLine + "\n\n")

	// Add a fancy divider
	divider := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	b.WriteString(divider + "\n\n")

	// Welcome message in a fancy box with tilde expansion info
	welcomeMsg := "Welcome to TUV!\n\nThis appears to be your first time running TUV.\nPlease configure the parent directory where your Python projects are located."
	welcome := FancyBoxStyle.Render(welcomeMsg)
	b.WriteString(welcome + "\n\n")

	input := InputStyle.Render(
		InputLabelStyle.Render("Parent Directory: ") + "\n" +
			m.textInput.View(),
	)
	b.WriteString(input + "\n\n")

	if m.error != "" {
		b.WriteString(ErrorStyle.Render("Error: "+m.error) + "\n\n")
	}

	// Update help text to remove 'd' key option
	help := HelpStyle.Render("Enter: Save Configuration • Ctrl+C: Quit")
	b.WriteString(help)

	return BaseStyle.Render(b.String())
}

// viewProjectList renders the project list
func (m Model) viewProjectList() string {
	var b strings.Builder

	// Add a compact logo
	b.WriteString(GetCompactLogo() + "\n")

	title := TitleStyle.Render("UV Projects")
	b.WriteString(title + "\n")

	// Add a fancy divider
	divider := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	b.WriteString(divider + "\n\n")

	if m.loading {
		loadingMsg := FancyBoxStyle.Render(fmt.Sprintf("%s %s", m.spinner.View(), m.loadingMsg))
		b.WriteString(loadingMsg + "\n")
	} else if len(m.projects) == 0 {
		emptyMsg := FancyBoxStyle.Render("No uv projects found.\n\nPress 's' to rescan for projects or 'Esc' to go back to the main menu.")
		b.WriteString(emptyMsg + "\n")
	} else {
		// Add project count with highlight
		countMsg := fmt.Sprintf("Found %s uv projects", HighlightStyle.Render(fmt.Sprintf("%d", len(m.projects))))
		b.WriteString(countMsg + "\n\n")

		var rows []string
		for i, project := range m.projects {
			// Only show name and size
			projectInfo := fmt.Sprintf("%s (%s)", project.Name, scanner.FormatSize(project.Size))
			if i == m.selectedProject {
				rows = append(rows, SelectedProjectStyle.Render(fmt.Sprintf(" > %s", projectInfo)))
			} else {
				rows = append(rows, ProjectStyle.Render(fmt.Sprintf("   %s", projectInfo)))
			}
		}

		// Join rows with newlines to ensure vertical layout
		listContent := strings.Join(rows, "\n")
		list := ProjectListStyle.Render(listContent)
		b.WriteString(list + "\n\n")
	}

	if !m.loading && m.statusMsg != "" {
		b.WriteString(StatusStyle.Render(m.statusMsg) + "\n")
	}

	if m.error != "" {
		b.WriteString(ErrorStyle.Render("Error: "+m.error) + "\n")
	}

	help := HelpStyle.Render("↑/↓: Navigate • Enter: Select • s: Rescan • Esc: Back • q: Quit")
	b.WriteString("\n" + help)

	return BaseStyle.Render(b.String())
}

// viewProjectDetail renders the project detail view
func (m Model) viewProjectDetail() string {
	var b strings.Builder
	project := m.projects[m.selectedProject]

	// Add a compact logo
	b.WriteString(GetCompactLogo() + "\n")

	// Project name with fancy styling
	title := TitleStyle.Render(project.Name)
	b.WriteString(title + "\n")

	// Add a fancy divider
	divider := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	b.WriteString(divider + "\n\n")

	// Create a fancy header for the project details
	detailsHeader := lipgloss.NewStyle().
		Foreground(highlightColor).
		Bold(true).
		Render("✧ PROJECT DETAILS ✧")
	b.WriteString(detailsHeader + "\n\n")

	var infoRows []string
	// Only show name, date created, python version, and size with proper spacing
	infoRows = append(infoRows, InfoTitleStyle.Render("Name: ")+InfoValueStyle.Render(project.Name))
	infoRows = append(infoRows, InfoTitleStyle.Render("Date Created: ")+InfoValueStyle.Render(project.LastModified.Format(time.RFC1123)))
	infoRows = append(infoRows, InfoTitleStyle.Render("Python Version: ")+InfoValueStyle.Render(project.PythonVersion))
	infoRows = append(infoRows, InfoTitleStyle.Render("Size: ")+InfoValueStyle.Render(scanner.FormatSize(project.Size)))

	// Join rows with newlines to ensure vertical layout
	infoContent := strings.Join(infoRows, "\n")
	info := FancyBoxStyle.Render(infoContent) // Use fancy box style for details
	b.WriteString(info + "\n\n")

	if m.error != "" {
		b.WriteString(ErrorStyle.Render("Error: "+m.error) + "\n")
	}

	help := HelpStyle.Render("Esc: Back • q: Quit")
	b.WriteString("\n" + help)

	return BaseStyle.Render(b.String())
}

// viewNewProject renders the new project creation screen
func (m Model) viewNewProject() string {
	var b strings.Builder

	title := TitleStyle.Render("Create New UV Project")
	b.WriteString(title + "\n\n")

	if m.loading {
		loadingMsg := FancyBoxStyle.Render(fmt.Sprintf("%s %s", m.spinner.View(), m.loadingMsg))
		b.WriteString(loadingMsg + "\n\n")
	} else {
		b.WriteString("Enter the name for your new UV project:\n\n")

		input := InputStyle.Render(
			InputLabelStyle.Render("Project Name: ") + "\n" +
				m.textInput.View(),
		)
		b.WriteString(input + "\n\n")
	}

	if m.error != "" {
		b.WriteString(ErrorStyle.Render("Error: "+m.error) + "\n\n")
	}

	help := HelpStyle.Render("Enter: Create Project • Esc: Back • Ctrl+C: Quit")
	b.WriteString(help)

	return BaseStyle.Render(b.String())
}

// viewLoading renders the loading screen
func (m Model) viewLoading() string {
	var b strings.Builder

	title := TitleStyle.Render("TUV - UV Environment Manager")
	b.WriteString(title + "\n\n")

	b.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), m.loadingMsg))

	if m.error != "" {
		b.WriteString(ErrorStyle.Render("Error: "+m.error) + "\n")
	}

	return BaseStyle.Render(b.String())
}

// Custom message types
type projectsFoundMsg struct {
	projects []scanner.UVProject
}

type errMsg struct {
	err error
}

type statusMsg struct {
	msg string
}

type projectCreatedMsg struct {
	projects    []scanner.UVProject
	projectName string
}
