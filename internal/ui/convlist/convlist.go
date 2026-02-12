package convlist

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/linkedin-tui/internal/ui/styles"
	"github.com/ggfevans/linkedin-tui/internal/util"
)

// Conversation represents a conversation list item.
type Conversation struct {
	ID           string
	Name         string
	LastMessage  string
	Timestamp    string
	Unread       bool
	UnreadCount  int
}

// Model represents the conversation list panel.
type Model struct {
	styles        styles.Styles
	width         int
	height        int
	focused       bool
	selected      int
	conversations []Conversation
	offset        int // scroll offset

	// Filter tabs
	filterTab    int // 0=Inbox, 1=Unread
	inboxCount   int
	unreadCount  int
}

// New creates a new conversation list model.
func New(s styles.Styles) Model {
	return Model{styles: s}
}

// SetSize updates dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Focus gives focus.
func (m *Model) Focus() { m.focused = true }

// Blur removes focus.
func (m *Model) Blur() { m.focused = false }

// Focused returns focus state.
func (m Model) Focused() bool { return m.focused }

// Selected returns the selected index.
func (m Model) Selected() int { return m.selected }

// SelectedConversation returns the selected conversation, if any.
func (m Model) SelectedConversation() (Conversation, bool) {
	if len(m.conversations) == 0 {
		return Conversation{}, false
	}
	return m.conversations[m.selected], true
}

// SetConversations replaces the conversation list.
func (m *Model) SetConversations(convs []Conversation) {
	m.conversations = convs
	if m.selected >= len(convs) {
		m.selected = max(len(convs)-1, 0)
	}
}

// SetStyles updates the styles.
func (m *Model) SetStyles(s styles.Styles) {
	m.styles = s
}

// FilterTab returns the active filter tab index.
func (m Model) FilterTab() int { return m.filterTab }

// ToggleFilter switches between Inbox and Unread tabs.
func (m *Model) ToggleFilter() {
	m.filterTab = 1 - m.filterTab
	m.selected = 0
	m.offset = 0
}

// SetFilterCounts updates the tab counts.
func (m *Model) SetFilterCounts(inbox, unread int) {
	m.inboxCount = inbox
	m.unreadCount = unread
}

// MoveDown moves selection down.
func (m *Model) MoveDown() {
	if m.selected < len(m.conversations)-1 {
		m.selected++
		m.ensureVisible()
	}
}

// MoveUp moves selection up.
func (m *Model) MoveUp() {
	if m.selected > 0 {
		m.selected--
		m.ensureVisible()
	}
}

// MoveToTop jumps to the first conversation.
func (m *Model) MoveToTop() {
	m.selected = 0
	m.offset = 0
}

// MoveToBottom jumps to the last conversation.
func (m *Model) MoveToBottom() {
	if len(m.conversations) > 0 {
		m.selected = len(m.conversations) - 1
		m.ensureVisible()
	}
}

func (m *Model) ensureVisible() {
	visible := m.visibleEntries()
	if visible < 1 {
		visible = 1
	}
	if m.selected < m.offset {
		m.offset = m.selected
	}
	if m.selected >= m.offset+visible {
		m.offset = m.selected - visible + 1
	}
}

// View renders the conversation list.
func (m Model) View() string {
	border := m.styles.BorderNormal
	if m.focused {
		border = m.styles.BorderFocused
	}

	contentWidth := m.width - 4 // border + padding
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Render filter tabs
	inboxLabel := fmt.Sprintf("Inbox %d", m.inboxCount)
	unreadLabel := fmt.Sprintf("Unread %d", m.unreadCount)
	sep := m.styles.Muted.Render(" · ")

	tabSelected := lipgloss.NewStyle().Foreground(m.styles.Theme.Secondary).Bold(true)
	tabNormal := m.styles.Muted

	var tabBar string
	if m.filterTab == 0 {
		tabBar = tabSelected.Render(inboxLabel) + sep + tabNormal.Render(unreadLabel)
	} else {
		tabBar = tabNormal.Render(inboxLabel) + sep + tabSelected.Render(unreadLabel)
	}
	content := tabBar + "\n"

	if len(m.conversations) == 0 {
		content += "\n" + m.styles.Muted.Render("  No conversations")
	} else {
		visible := m.visibleEntries()
		if visible < 1 {
			visible = 1
		}
		end := m.offset + visible
		if end > len(m.conversations) {
			end = len(m.conversations)
		}

		accentBar := lipgloss.NewStyle().Foreground(m.styles.Theme.Secondary).Render("▎")

		for i := m.offset; i < end; i++ {
			c := m.conversations[i]
			name := util.Truncate(c.Name, contentWidth-2)
			preview := util.Truncate(c.LastMessage, contentWidth-2)
			ts := c.Timestamp

			nameStyle := lipgloss.NewStyle().Foreground(m.styles.Theme.Foreground)
			prefix := "  "
			if c.Unread && i == m.selected {
				nameStyle = nameStyle.Bold(true)
				prefix = accentBar + m.styles.Unread.Render("●")
			} else if i == m.selected {
				prefix = accentBar + " "
			} else if c.Unread {
				nameStyle = nameStyle.Bold(true)
				prefix = m.styles.Unread.Render("● ")
			}

			line := prefix + nameStyle.Render(name)
			if ts != "" {
				gap := contentWidth - lipgloss.Width(prefix) - lipgloss.Width(name) - lipgloss.Width(ts)
				if gap < 1 {
					gap = 1
				}
				line = prefix + nameStyle.Render(name) + lipgloss.NewStyle().Width(gap).Render("") + m.styles.Muted.Render(ts)
			}

			previewLine := "  " + m.styles.Muted.Render(preview)

			entry := line + "\n" + previewLine
			entryStyle := lipgloss.NewStyle().Width(contentWidth)

			content += "\n" + entryStyle.Render(entry)
		}
	}

	innerHeight := m.height - 2 // subtract top/bottom border
	content = util.PadToHeight(content, innerHeight)
	return border.
		Width(m.width - 2).
		Render(content)
}

func (m Model) visibleEntries() int {
	// borders(2) + title line(1) + gap(1) = 4 lines of overhead
	visibleLines := m.height - 4
	if visibleLines < 1 {
		return 1
	}
	// Each entry takes ~3 lines: name+timestamp, preview, gap
	entries := visibleLines / 3
	if entries < 1 {
		entries = 1
	}
	return entries
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Count returns the number of conversations.
func (m Model) Count() int {
	return len(m.conversations)
}

// UnreadCount returns the number of unread conversations.
func (m Model) UnreadCount() int {
	count := 0
	for _, c := range m.conversations {
		if c.Unread {
			count++
		}
	}
	return count
}

// Conversations returns the full conversation list.
func (m Model) Conversations() []Conversation {
	return m.conversations
}

func (m Model) ConversationIndex(id string) int {
	for i, c := range m.conversations {
		if c.ID == id {
			return i
		}
	}
	return -1
}
