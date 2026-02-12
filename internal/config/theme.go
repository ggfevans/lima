package config

import "github.com/charmbracelet/lipgloss"

// Theme defines a color palette for the application.
type Theme struct {
	Name string

	// Base colors
	Background    lipgloss.Color
	Foreground    lipgloss.Color
	CurrentLine   lipgloss.Color
	Selection     lipgloss.Color
	Comment       lipgloss.Color
	Subtle        lipgloss.Color

	// Accent colors
	Primary       lipgloss.Color
	Secondary     lipgloss.Color
	Success       lipgloss.Color
	Warning       lipgloss.Color
	Error         lipgloss.Color
	Info          lipgloss.Color

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
	Comment:       lipgloss.Color("#6272a4"),
	Subtle:        lipgloss.Color("#6272a4"),
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

// LinkedIn is the light theme using LinkedIn brand colors.
var LinkedIn = Theme{
	Name:          "linkedin",
	Background:    lipgloss.Color("#ffffff"),
	Foreground:    lipgloss.Color("#000000"),
	CurrentLine:   lipgloss.Color("#e8f0fe"),
	Selection:     lipgloss.Color("#d0e4fc"),
	Comment:       lipgloss.Color("#666666"),
	Subtle:        lipgloss.Color("#999999"),
	Primary:       lipgloss.Color("#0a66c2"),
	Secondary:     lipgloss.Color("#004182"),
	Success:       lipgloss.Color("#057642"),
	Warning:       lipgloss.Color("#915907"),
	Error:         lipgloss.Color("#cc1016"),
	Info:          lipgloss.Color("#0a66c2"),
	Accent:        lipgloss.Color("#0a66c2"),
	AccentDim:     lipgloss.Color("#004182"),
	Border:        lipgloss.Color("#e0e0e0"),
	BorderFocused: lipgloss.Color("#0a66c2"),
	Unread:        lipgloss.Color("#057642"),
	OwnMessage:    lipgloss.Color("#e8f0fe"),
	OtherMessage:  lipgloss.Color("#f3f6f8"),
	OwnSender:     lipgloss.Color("#004182"),
}

// DefaultTheme returns the default application theme.
func DefaultTheme() Theme {
	return Dracula
}

// ThemeByName returns a theme by name, falling back to Dracula.
func ThemeByName(name string) Theme {
	switch name {
	case "linkedin", "light":
		return LinkedIn
	default:
		return Dracula
	}
}
