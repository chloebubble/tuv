package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// ASCII art logo for TUV
const logoArt = `
████████╗██╗   ██╗██╗   ██╗
╚══██╔══╝██║   ██║██║   ██║
   ██║   ██║   ██║██║   ██║
   ██║   ██║   ██║╚██╗ ██╔╝
   ██║   ╚██████╔╝ ╚████╔╝ 
   ╚═╝    ╚═════╝   ╚═══╝  
`

// GetLogo returns a styled logo
func GetLogo() string {
	return lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(logoArt)
}

// GetCompactLogo returns a smaller version of the logo for screens with limited space
const compactLogoArt = `
 ████████╗██╗   ██╗██╗   ██╗
 ╚══██╔══╝██║   ██║██║   ██║
    ██║   ██║   ██║╚██████╔╝
    ╚═╝    ╚═════╝  ╚═════╝ 
`

// GetCompactLogo returns a styled compact logo
func GetCompactLogo() string {
	return lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(compactLogoArt)
}
