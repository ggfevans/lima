package config

import "github.com/charmbracelet/lipgloss"

// Theme defines a color palette for the application.
type Theme struct {
	Name string

	// Base colors
	Background  lipgloss.Color
	Foreground  lipgloss.Color
	CurrentLine lipgloss.Color
	Selection   lipgloss.Color
	Comment     lipgloss.Color
	Subtle      lipgloss.Color

	// Accent colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Info      lipgloss.Color

	// Semantic
	Accent        lipgloss.Color
	AccentDim     lipgloss.Color
	Border        lipgloss.Color
	BorderFocused lipgloss.Color
	Unread        lipgloss.Color
	OwnMessage    lipgloss.Color
	OtherMessage  lipgloss.Color
	OwnSender     lipgloss.Color
}

// Dracula is the default dark theme.
var Dracula = Theme{
	Name:          "dracula",
	Background:    lipgloss.Color("#282a36"),
	Foreground:    lipgloss.Color("#f8f8f2"),
	CurrentLine:   lipgloss.Color("#44475a"),
	Selection:     lipgloss.Color("#44475a"),
	Comment:       lipgloss.Color("#7e8eb8"),
	Subtle:        lipgloss.Color("#8a9bc4"),
	Primary:       lipgloss.Color("#bd93f9"),
	Secondary:     lipgloss.Color("#ff79c6"),
	Success:       lipgloss.Color("#50fa7b"),
	Warning:       lipgloss.Color("#f1fa8c"),
	Error:         lipgloss.Color("#ff5555"),
	Info:          lipgloss.Color("#8be9fd"),
	Accent:        lipgloss.Color("#bd93f9"),
	AccentDim:     lipgloss.Color("#6272a4"),
	Border:        lipgloss.Color("#44475a"),
	BorderFocused: lipgloss.Color("#bd93f9"),
	Unread:        lipgloss.Color("#50fa7b"),
	OwnMessage:    lipgloss.Color("#44475a"),
	OtherMessage:  lipgloss.Color("#313345"),
	OwnSender:     lipgloss.Color("#8be9fd"),
}

// DefaultTheme returns the default application theme.
func DefaultTheme() Theme {
	return Dracula
}

// ThemeByName returns a theme by name, falling back to Dracula.
func ThemeByName(name string) Theme {
	switch name {
	default:
		return Dracula
	}
}
