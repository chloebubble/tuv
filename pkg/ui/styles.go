package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#AE88FD")
	accentColor    = lipgloss.Color("#43BF6D")
	warningColor   = lipgloss.Color("#F2B705")
	errorColor     = lipgloss.Color("#F25757")
	textColor      = lipgloss.Color("#FFFFFF")
	dimTextColor   = lipgloss.Color("#CCCCCC")
	bgColor        = lipgloss.Color("#1A1B26")
	highlightColor = lipgloss.Color("#FF79C6") // New color for highlights

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Background(bgColor).
			Foreground(textColor)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			MarginBottom(1).
			PaddingLeft(2).
			PaddingRight(2)

	// Subtitle style
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			MarginLeft(1).
			AlignVertical(lipgloss.Center)

	// Menu styles
	MenuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1).
			MarginRight(2).
			Width(40)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1).
				Width(38)

	ItemStyle = lipgloss.NewStyle().
			Foreground(dimTextColor).
			Width(38)

	// Project list styles
	ProjectListStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(secondaryColor).
				Padding(1).
				Width(60)

	SelectedProjectStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1).
				Width(58)

	ProjectStyle = lipgloss.NewStyle().
			Foreground(dimTextColor).
			Width(58)

	// Info styles
	InfoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1).
			Width(60)

	InfoTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Width(16)

	InfoValueStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Width(42)

	// Status styles
	StatusStyle = lipgloss.NewStyle().
			Foreground(dimTextColor).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Input styles
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(dimTextColor).
			Italic(true)

	// Fancy box styles with double borders
	FancyBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(highlightColor).
			Padding(1).
			Width(60)

	// Highlight text style
	HighlightStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
			Bold(true)

	// Version badge style
	VersionBadgeStyle = lipgloss.NewStyle().
				Foreground(textColor).
				Background(secondaryColor).
				Padding(0, 1).
				AlignVertical(lipgloss.Center)
)
