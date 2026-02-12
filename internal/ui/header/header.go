package header

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/li-cli/internal/ui/styles"
)

// Model represents the header bar.
type Model struct {
	styles    styles.Styles
	width     int
	username  string
	connected bool
}

// New creates a new header model.
func New(s styles.Styles) Model {
	return Model{
		styles:    s,
		username:  "",
		connected: false,
	}
}

// SetWidth updates the header width.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// SetUsername sets the displayed username.
func (m *Model) SetUsername(u string) {
	m.username = u
}

// SetConnected sets the connection status.
func (m *Model) SetConnected(c bool) {
	m.connected = c
}

// SetStyles updates the styles.
func (m *Model) SetStyles(s styles.Styles) {
	m.styles = s
}

// View renders the header.
func (m Model) View() string {
	logo := m.styles.AccentText.Render("Li-CLI")

	var parts []string
	parts = append(parts, logo)

	if m.username != "" {
		parts = append(parts, m.styles.Muted.Render("@"+m.username))
	}

	if m.connected {
		parts = append(parts, m.styles.Connected.Render("★"))
	} else {
		parts = append(parts, m.styles.Disconnected.Render("★"))
	}

	left := " " + joinParts(parts, "  ")

	return lipgloss.NewStyle().Width(m.width).Render(left)
}

func joinParts(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
