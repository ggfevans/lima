package compose

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/endorse/internal/ui/styles"
)

// Model represents the compose/reply textarea.
type Model struct {
	styles    styles.Styles
	textarea  textarea.Model
	width     int
	height    int
	focused   bool
	recipient string
}

// New creates a new compose model.
func New(s styles.Styles) Model {
	ta := textarea.New()
	ta.Placeholder = "Type a message..."
	ta.CharLimit = 8000
	ta.ShowLineNumbers = false
	ta.SetHeight(3)

	// Remove the default full-line highlight
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.BlurredStyle.CursorLine = lipgloss.NewStyle()

	// Use teal cursor
	ta.Cursor.Style = lipgloss.NewStyle().Foreground(s.Theme.OwnSender)

	// Muted placeholder
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(s.Theme.Subtle)
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(s.Theme.Subtle)

	return Model{
		styles:   s,
		textarea: ta,
		height:   5,
	}
}

// SetSize updates dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.textarea.SetWidth(w)
	m.textarea.SetHeight(h)
}

// Focus gives focus to the compose box and returns the cursor blink cmd.
func (m *Model) Focus() tea.Cmd {
	m.focused = true
	return m.textarea.Focus()
}

// Blur removes focus.
func (m *Model) Blur() {
	m.focused = false
	m.textarea.Blur()
}

// Deactivate resets the compose box.
func (m *Model) Deactivate() {
	m.focused = false
	m.textarea.Blur()
	m.textarea.Reset()
	m.recipient = ""
}

// Focused returns focus state.
func (m Model) Focused() bool { return m.focused }

// SetRecipient sets who we're replying to.
func (m *Model) SetRecipient(name string) {
	m.recipient = name
	m.textarea.Placeholder = "Reply to " + name + "..."
}

// SetStyles updates the styles.
func (m *Model) SetStyles(s styles.Styles) {
	m.styles = s
}

// Value returns the current text content.
func (m Model) Value() string {
	return m.textarea.Value()
}

// Reset clears the textarea.
func (m *Model) Reset() {
	m.textarea.Reset()
}

// Update handles tea messages for the textarea.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// View renders the compose textarea content (no border — thread handles that).
func (m Model) View() string {
	return m.textarea.View()
}

// ComposeHeight returns the compose box height.
func (m Model) ComposeHeight() int {
	return m.height
}

// CursorLine returns the cursor line position in the textarea.
func (m Model) CursorLine() int {
	return m.textarea.Line()
}

// SetPlaceholder sets a custom placeholder.
func (m *Model) SetPlaceholder(p string) {
	m.textarea.Placeholder = p
}

// SetFocusStyle allows customizing focus visual feedback.
func (m *Model) SetFocusStyle(focused lipgloss.Style) {
	m.textarea.FocusedStyle.Base = focused
}
