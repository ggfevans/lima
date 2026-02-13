package statusbar

import (
	"fmt"
	"strings"

	"github.com/ggfevans/endorse/internal/ui/styles"
)

// Hint represents a keybinding hint.
type Hint struct {
	Key  string
	Desc string
}

// Model represents the status bar.
type Model struct {
	styles    styles.Styles
	width     int
	hints     []Hint
	err       string
	username  string
	connected bool
}

// New creates a new status bar model.
func New(s styles.Styles) Model {
	return Model{
		styles: s,
		hints: []Hint{
			{Key: "↑↓/jk", Desc: "Navigate"},
			{Key: "←→/Tab", Desc: "Focus"},
			{Key: "Enter", Desc: "Select"},
			{Key: "f", Desc: "Filter"},
			{Key: "r", Desc: "Reply"},
			{Key: "m", Desc: "Read/Unread"},
			{Key: "d", Desc: "Delete"},
			{Key: "q", Desc: "Quit"},
		},
	}
}

// SetWidth updates the status bar width.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// SetHints replaces the keybinding hints.
func (m *Model) SetHints(hints []Hint) {
	m.hints = hints
}

// SetError sets an error message to display.
func (m *Model) SetError(err string) {
	m.err = err
}

// ClearError clears any error message.
func (m *Model) ClearError() {
	m.err = ""
}

// SetUsername sets the displayed username.
func (m *Model) SetUsername(u string) {
	m.username = u
}

// SetConnected sets the connection status indicator.
func (m *Model) SetConnected(c bool) {
	m.connected = c
}

// SetStyles updates the styles.
func (m *Model) SetStyles(s styles.Styles) {
	m.styles = s
}

// View renders the status bar.
func (m Model) View() string {
	if m.err != "" {
		errStyle := m.styles.StatusBar.Foreground(m.styles.Theme.Error)
		return errStyle.Width(m.width).Render(fmt.Sprintf(" ERROR: %s", m.err))
	}

	// Key hints (left-aligned)
	var hintParts []string
	for _, h := range m.hints {
		key := m.styles.StatusKey.Render(h.Key)
		hintParts = append(hintParts, fmt.Sprintf("%s %s", key, h.Desc))
	}
	line := " " + strings.Join(hintParts, "  ")

	return m.styles.StatusBar.Width(m.width).Render(line)
}
