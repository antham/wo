package cmd

import "github.com/charmbracelet/lipgloss"

var regularStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#3C3744"))

var highlightedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#3D52D5"))

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#090C9B"))

var separator = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#B4C5E4")).
	Render("---")

func applyDarkTheme() {
	regularStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8ECDDD"))
	highlightedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFADD"))
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFCC70"))
	separator = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#22668D")).
		Render("---")
}
