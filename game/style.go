package game

import "github.com/charmbracelet/lipgloss"

var (
	defaultStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("#52bf90")).
			Background(lipgloss.Color("#252627"))
	onStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#BF5281")).
		Background(lipgloss.Color("#252627"))
	offStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#607A9D")).
			Background(lipgloss.Color("#252627"))
	completeStyle = lipgloss.NewStyle().
			Width(32).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#e990c5"))
)
