package thread

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/linkedin-tui/internal/ui/styles"
	"github.com/ggfevans/linkedin-tui/internal/util"
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
	scrollOffset   int
	conversationID string
}

// New creates a new thread model.
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

// SetStyles updates the styles.
func (m *Model) SetStyles(s styles.Styles) {
	m.styles = s
}

// SetConversation sets the current conversation.
func (m *Model) SetConversation(id, subject string) {
	m.conversationID = id
	m.subject = subject
	m.messages = nil
	m.scrollOffset = 0
}

// ConversationID returns the current conversation ID.
func (m Model) ConversationID() string {
	return m.conversationID
}

// SetMessages replaces the message list.
func (m *Model) SetMessages(msgs []Message) {
	m.messages = msgs
	// Scroll to bottom on new messages
	m.scrollToBottom()
}

// AppendMessage adds a message at the end.
func (m *Model) AppendMessage(msg Message) {
	m.messages = append(m.messages, msg)
	m.scrollToBottom()
}

// ScrollUp scrolls the view up.
func (m *Model) ScrollUp(lines int) {
	m.scrollOffset -= lines
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// ScrollDown scrolls the view down.
func (m *Model) ScrollDown(lines int) {
	m.scrollOffset += lines
	maxOffset := m.maxScrollOffset()
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
}

func (m *Model) scrollToBottom() {
	m.scrollOffset = m.maxScrollOffset()
}

func (m Model) maxScrollOffset() int {
	totalLines := m.renderedLineCount()
	visible := m.height - 2 - 1 // borders + title line
	if visible < 1 {
		visible = 1
	}
	if totalLines <= visible {
		return 0
	}
	return totalLines - visible
}

func (m Model) renderedLineCount() int {
	if len(m.messages) == 0 {
		return 1
	}
	count := 0
	var prevSender string
	for _, msg := range m.messages {
		if msg.Sender != prevSender {
			if prevSender != "" {
				count++ // divider
			}
			count++ // header
			prevSender = msg.Sender
		}
		count++ // body
	}
	return count
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

	title := m.styles.AccentText.Render("MESSAGES")
	if m.subject != "" {
		title = m.styles.AccentText.Render(util.Truncate(m.subject, contentWidth))
	}
	content := title + "\n"

	if m.conversationID == "" {
		content += "\n" + m.styles.Muted.Render("  Select a conversation")
	} else if len(m.messages) == 0 {
		content += "\n" + m.styles.Muted.Render("  No messages")
	} else {
		accentBar := lipgloss.NewStyle().Foreground(m.styles.Theme.OwnSender).Render("▎")
		divider := m.styles.Muted.Render(strings.Repeat("─", contentWidth-2))

		var lines []string
		var prevSender string
		for _, msg := range m.messages {
			// Group header: only when sender changes
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

			// Body line with accent bar for own messages
			prefix := " "
			if msg.IsOwn {
				prefix = accentBar
			}
			lines = append(lines, prefix+msg.Body)
		}

		// Apply scroll offset
		visible := m.height - 2 - 1 // borders + title line
		if visible < 1 {
			visible = 1
		}
		start := m.scrollOffset
		if start > len(lines) {
			start = len(lines)
		}
		end := start + visible
		if end > len(lines) {
			end = len(lines)
		}

		content += strings.Join(lines[start:end], "\n")
	}

	innerHeight := m.height - 2 // subtract top/bottom border
	content = util.PadToHeight(content, innerHeight)
	return border.
		Width(m.width - 2).
		Render(content)
}

// Clear resets the thread.
func (m *Model) Clear() {
	m.conversationID = ""
	m.subject = ""
	m.messages = nil
	m.scrollOffset = 0
}

// HasConversation returns whether a conversation is loaded.
func (m Model) HasConversation() bool {
	return m.conversationID != ""
}

// AtTop returns true if scrolled to the top.
func (m Model) AtTop() bool {
	return m.scrollOffset == 0
}

// MessageCount returns the number of messages.
func (m Model) MessageCount() int {
	return len(m.messages)
}

func (m Model) VisibleHeight() int {
	v := m.height - 2 - 1 // borders + title line
	if v < 1 {
		return 1
	}
	return v
}

func (m Model) ScrollPercent() float64 {
	maxOff := m.maxScrollOffset()
	if maxOff == 0 {
		return 1.0
	}
	return float64(m.scrollOffset) / float64(maxOff)
}
