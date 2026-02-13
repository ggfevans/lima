package modal

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/lima/internal/ui/styles"
)

// ConfirmModel is a centered confirmation dialog.
type ConfirmModel struct {
	styles  styles.Styles
	width   int
	height  int
	message string
	active  bool
}

// NewConfirm creates a new confirmation modal.
func NewConfirm(s styles.Styles) ConfirmModel {
	return ConfirmModel{styles: s}
}

// Show displays the modal with the given message.
func (m *ConfirmModel) Show(message string) {
	m.message = message
	m.active = true
}

// Hide dismisses the modal.
func (m *ConfirmModel) Hide() {
	m.active = false
	m.message = ""
}

// Active returns whether the modal is showing.
func (m ConfirmModel) Active() bool {
	return m.active
}

// SetSize updates the modal dimensions.
func (m *ConfirmModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// View renders the confirmation modal centered on screen.
func (m ConfirmModel) View() string {
	if !m.active {
		return ""
	}

	title := m.styles.AccentText.Render("Confirm Delete")
	msg := lipgloss.NewStyle().Foreground(m.styles.Theme.Foreground).Render(m.message)
	hint := m.styles.Muted.Render("Enter to confirm  |  Esc to cancel")

	content := fmt.Sprintf("%s\n\n%s\n\n%s", title, msg, hint)

	boxWidth := 50
	if m.width > 0 && m.width < boxWidth+10 {
		boxWidth = m.width - 10
	}
	if boxWidth < 30 {
		boxWidth = 30
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Theme.Error).
		Padding(1, 3).
		Width(boxWidth)

	box := boxStyle.Render(content)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
