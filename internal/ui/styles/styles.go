package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/endorse/internal/config"
)

// Styles holds all shared lipgloss styles derived from the active theme.
type Styles struct {
	Theme config.Theme

	// Panel borders
	BorderNormal  lipgloss.Style
	BorderFocused lipgloss.Style

	// Text
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Muted      lipgloss.Style
	Bold       lipgloss.Style
	AccentText lipgloss.Style

	// Status indicators
	Connected    lipgloss.Style
	Disconnected lipgloss.Style
	Unread       lipgloss.Style

	// Header / status bar
	Header    lipgloss.Style
	StatusBar lipgloss.Style
	StatusKey lipgloss.Style

	// Conversation list
	ConvSelected lipgloss.Style
	ConvNormal   lipgloss.Style
	ConvUnread   lipgloss.Style

	// Messages
	OwnBubble     lipgloss.Style
	OtherBubble   lipgloss.Style
	Timestamp     lipgloss.Style
	SenderName    lipgloss.Style
	OwnSenderName lipgloss.Style

	// Compose
	ComposeCursor lipgloss.Style
}

// New creates a Styles set from a theme.
func New(theme config.Theme) Styles {
	s := Styles{Theme: theme}

	s.BorderNormal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border)

	s.BorderFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderFocused)

	s.Title = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Bold(true)

	s.Subtitle = lipgloss.NewStyle().
		Foreground(theme.Comment)

	s.Muted = lipgloss.NewStyle().
		Foreground(theme.Subtle)

	s.Bold = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Bold(true)

	s.AccentText = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	s.Connected = lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	s.Disconnected = lipgloss.NewStyle().
		Foreground(theme.Error).
		Bold(true)

	s.Unread = lipgloss.NewStyle().
		Foreground(theme.Unread).
		Bold(true)

	s.Header = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Bold(true).
		Padding(0, 1)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(theme.Comment).
		Padding(0, 1)

	s.StatusKey = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	s.ConvSelected = lipgloss.NewStyle().
		Background(theme.Selection).
		Foreground(theme.Foreground)

	s.ConvNormal = lipgloss.NewStyle().
		Foreground(theme.Foreground)

	s.ConvUnread = lipgloss.NewStyle().
		Foreground(theme.Foreground).
		Bold(true)

	s.OwnBubble = lipgloss.NewStyle().
		Background(theme.OwnMessage).
		Foreground(theme.Foreground).
		Padding(0, 1)

	s.OtherBubble = lipgloss.NewStyle().
		Background(theme.OtherMessage).
		Foreground(theme.Foreground).
		Padding(0, 1)

	s.Timestamp = lipgloss.NewStyle().
		Foreground(theme.Subtle)

	s.SenderName = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	s.OwnSenderName = lipgloss.NewStyle().
		Foreground(theme.OwnSender).
		Bold(true)

	s.ComposeCursor = lipgloss.NewStyle().
		Foreground(theme.OwnSender)

	return s
}
