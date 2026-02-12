package compose

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/linkedin-tui/internal/ui/styles"
)

// Model represents the compose/reply textarea.
type Model struct {
	styles    styles.Styles
	textarea  textarea.Model
	width     int
	height    int
	focused   bool
	active    bool
	recipient string
}

// New creates a new compose model.
func New(s styles.Styles) Model {
	ta := textarea.New()
	ta.Placeholder = "Type a message..."
	ta.CharLimit = 8000
	ta.ShowLineNumbers = false
	ta.SetHeight(3)

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
	m.textarea.SetWidth(w - 4)
	m.textarea.SetHeight(h - 2)
}

// Focus gives focus and activates the compose box.
func (m *Model) Focus() {
	m.focused = true
	m.active = true
	m.textarea.Focus()
}

// Blur removes focus.
func (m *Model) Blur() {
	m.focused = false
	m.textarea.Blur()
}

// Deactivate hides the compose box.
func (m *Model) Deactivate() {
	m.active = false
	m.focused = false
	m.textarea.Blur()
	m.textarea.Reset()
	m.recipient = ""
}

// Focused returns focus state.
func (m Model) Focused() bool { return m.focused }

// Active returns whether the compose box is visible/active.
func (m Model) Active() bool { return m.active }

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
	if !m.active || !m.focused {
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// View renders the compose box.
func (m Model) View() string {
	if !m.active {
		return ""
	}

	border := m.styles.BorderNormal
	if m.focused {
		border = m.styles.BorderFocused
	}

	content := m.textarea.View()

	return border.
		Width(m.width - 2).
		Render(content)
}

// Height returns the compose box height when active.
func (m Model) ComposeHeight() int {
	if !m.active {
		return 0
	}
	return m.height
}

// CursorLine returns the cursor line position in the textarea.
func (m Model) CursorLine() int {
	return m.textarea.Line()
}

// Placeholder sets a custom placeholder.
func (m *Model) SetPlaceholder(p string) {
	m.textarea.Placeholder = p
}

// SetFocusStyle allows customizing focus visual feedback.
func (m *Model) SetFocusStyle(focused lipgloss.Style) {
	m.textarea.FocusedStyle.Base = focused
}
