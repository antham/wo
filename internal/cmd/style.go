package cmd

import "github.com/charmbracelet/lipgloss"

var regularStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8ECDDD"))

var highlightedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFADD"))

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFCC70"))

var separator = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#22668D")).
	Render("---")
