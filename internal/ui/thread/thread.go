package thread

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/lima/internal/ui/styles"
	"github.com/ggfevans/lima/internal/util"
)

// Message represents a single message in a thread.
type Message struct {
	ID        string
	Sender    string
	Body      string
	Timestamp string
	IsOwn     bool
}

// Model represents the message thread panel.
type Model struct {
	styles         styles.Styles
	width          int
	height         int
	focused        bool
	subject        string
	messages       []Message
	conversationID string
	viewport       viewport.Model
}

// New creates a new thread model.
func New(s styles.Styles) Model {
	vp := viewport.New(0, 0)
	return Model{styles: s, viewport: vp}
}

// SetSize updates dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	contentWidth := w - 4 // border + padding
	if contentWidth < 1 {
		contentWidth = 1
	}
	visibleH := h - 2 - 1 // border + title line
	if visibleH < 1 {
		visibleH = 1
	}
	m.viewport.Width = contentWidth
	m.viewport.Height = visibleH
	m.refreshContent()
}

// Focus gives focus.
func (m *Model) Focus() { m.focused = true }

// Blur removes focus.
func (m *Model) Blur() { m.focused = false }

// Focused returns focus state.
func (m Model) Focused() bool { return m.focused }

// SetStyles updates the styles.
func (m *Model) SetStyles(s styles.Styles) {
	m.styles = s
}

// SetConversation sets the current conversation.
func (m *Model) SetConversation(id, subject string) {
	m.conversationID = id
	m.subject = subject
	m.messages = nil
	m.viewport.SetContent("")
	m.viewport.GotoTop()
}

// ConversationID returns the current conversation ID.
func (m Model) ConversationID() string {
	return m.conversationID
}

// SetMessages replaces the message list.
func (m *Model) SetMessages(msgs []Message) {
	m.messages = msgs
	m.refreshContent()
	m.viewport.GotoBottom()
}

// AppendMessage adds a message at the end.
func (m *Model) AppendMessage(msg Message) {
	m.messages = append(m.messages, msg)
	m.refreshContent()
	m.viewport.GotoBottom()
}

// ScrollUp scrolls the view up.
func (m *Model) ScrollUp(lines int) {
	m.viewport.LineUp(lines)
}

// ScrollDown scrolls the view down.
func (m *Model) ScrollDown(lines int) {
	m.viewport.LineDown(lines)
}

// View renders the thread panel.
func (m Model) View() string {
	border := m.styles.BorderNormal
	if m.focused {
		border = m.styles.BorderFocused
	}

	contentWidth := m.width - 4
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Build title with scroll percentage
	title := m.styles.AccentText.Render("MESSAGES")
	if m.subject != "" {
		titleText := util.Truncate(m.subject, contentWidth)
		if len(m.messages) > 0 && m.viewport.TotalLineCount() > m.viewport.Height {
			pct := int(m.viewport.ScrollPercent() * 100)
			suffix := fmt.Sprintf("  %d%%", pct)
			titleText = util.Truncate(m.subject, contentWidth-len(suffix)) + suffix
		}
		title = m.styles.AccentText.Render(titleText)
	}

	if m.conversationID == "" {
		// No conversation selected — show placeholder
		placeholder := "\n" + m.styles.Muted.Render("  Select a conversation")
		innerHeight := m.height - 2
		content := title + "\n" + placeholder
		content = util.PadToHeight(content, innerHeight)
		return border.Width(m.width - 2).Render(content)
	}

	if len(m.messages) == 0 {
		// Conversation selected but no messages
		placeholder := "\n" + m.styles.Muted.Render("  No messages")
		innerHeight := m.height - 2
		content := title + "\n" + placeholder
		content = util.PadToHeight(content, innerHeight)
		return border.Width(m.width - 2).Render(content)
	}

	// Render viewport content
	innerHeight := m.height - 2 // subtract top/bottom border
	content := title + "\n" + m.viewport.View()
	content = util.PadToHeight(content, innerHeight)
	return border.Width(m.width - 2).Render(content)
}

// refreshContent rebuilds the viewport content string from messages.
func (m *Model) refreshContent() {
	if len(m.messages) == 0 {
		m.viewport.SetContent("")
		return
	}

	contentWidth := m.width - 4
	if contentWidth < 1 {
		contentWidth = 1
	}

	accentBar := lipgloss.NewStyle().Foreground(m.styles.Theme.OwnSender).Render("▎")
	divider := m.styles.Muted.Render(strings.Repeat("─", contentWidth-2))

	var lines []string
	var prevSender string
	for _, msg := range m.messages {
		if msg.Sender != prevSender {
			if prevSender != "" {
				lines = append(lines, divider)
			}

			senderStyle := m.styles.SenderName
			if msg.IsOwn {
				senderStyle = m.styles.OwnSenderName
			}

			header := fmt.Sprintf("%s  %s",
				senderStyle.Render(msg.Sender),
				m.styles.Timestamp.Render(msg.Timestamp),
			)
			lines = append(lines, header)
			prevSender = msg.Sender
		}

		prefix := " "
		if msg.IsOwn {
			prefix = accentBar
		}
		lines = append(lines, prefix+msg.Body)
	}

	m.viewport.SetContent(strings.Join(lines, "\n"))
}

// Clear resets the thread.
func (m *Model) Clear() {
	m.conversationID = ""
	m.subject = ""
	m.messages = nil
	m.viewport.SetContent("")
	m.viewport.GotoTop()
}

// HasConversation returns whether a conversation is loaded.
func (m Model) HasConversation() bool {
	return m.conversationID != ""
}

// AtTop returns true if scrolled to the top.
func (m Model) AtTop() bool {
	return m.viewport.AtTop()
}

// MessageCount returns the number of messages.
func (m Model) MessageCount() int {
	return len(m.messages)
}

// VisibleHeight returns the viewport height.
func (m Model) VisibleHeight() int {
	if m.viewport.Height < 1 {
		return 1
	}
	return m.viewport.Height
}

// ScrollPercent returns the viewport scroll percentage.
func (m Model) ScrollPercent() float64 {
	return m.viewport.ScrollPercent()
}
