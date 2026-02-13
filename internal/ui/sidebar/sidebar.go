package sidebar

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/lima/internal/ui/styles"
	"github.com/ggfevans/lima/internal/util"
)

// Folder represents a sidebar folder entry.
type Folder struct {
	Name  string
	Count int
}

// Model represents the sidebar panel.
type Model struct {
	styles   styles.Styles
	width    int
	height   int
	focused  bool
	selected int
	folders  []Folder
}

// New creates a new sidebar model.
func New(s styles.Styles) Model {
	return Model{
		styles: s,
		folders: []Folder{
			{Name: "Inbox", Count: 0},
			{Name: "Unread", Count: 0},
		},
	}
}

// SetSize updates the sidebar dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Focus gives focus to the sidebar.
func (m *Model) Focus() { m.focused = true }

// Blur removes focus from the sidebar.
func (m *Model) Blur() { m.focused = false }

// Focused returns whether the sidebar is focused.
func (m Model) Focused() bool { return m.focused }

// Selected returns the index of the selected folder.
func (m Model) Selected() int { return m.selected }

// SetStyles updates the styles.
func (m *Model) SetStyles(s styles.Styles) {
	m.styles = s
}

// SetCounts updates the folder counts.
func (m *Model) SetCounts(inbox, unread int) {
	m.folders[0].Count = inbox
	m.folders[1].Count = unread
}

// MoveDown moves selection down.
func (m *Model) MoveDown() {
	if m.selected < len(m.folders)-1 {
		m.selected++
	}
}

// MoveUp moves selection up.
func (m *Model) MoveUp() {
	if m.selected > 0 {
		m.selected--
	}
}

// View renders the sidebar.
func (m Model) View() string {
	border := m.styles.BorderNormal
	if m.focused {
		border = m.styles.BorderFocused
	}

	title := m.styles.AccentText.Render("FOLDERS")
	content := title + "\n\n"

	accentBar := lipgloss.NewStyle().Foreground(m.styles.Theme.Secondary).Render("▎")

	for i, f := range m.folders {
		name := f.Name
		countStr := ""
		if f.Count > 0 {
			countStr = fmt.Sprintf(" %d", f.Count)
		}

		w := m.width - 4 // account for border + padding
		if w < 1 {
			w = 1
		}

		prefix := "  "
		if i == m.selected && m.focused {
			prefix = accentBar + " "
		}

		line := fmt.Sprintf("%s%s%s", prefix, name, countStr)
		lineStyle := lipgloss.NewStyle().Width(w)

		if i == m.selected && m.focused {
			lineStyle = lineStyle.
				Foreground(m.styles.Theme.Foreground)
		} else if i == m.selected {
			lineStyle = lineStyle.
				Foreground(m.styles.Theme.Foreground).
				Bold(true)
		}

		content += lineStyle.Render(line) + "\n"
	}

	innerHeight := m.height - 2 // subtract top/bottom border
	content = util.PadToHeight(content, innerHeight)
	return border.
		Width(m.width - 2). // subtract border width
		Render(content)
}
