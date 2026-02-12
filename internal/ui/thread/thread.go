package thread

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/endorse/internal/ui/styles"
	"github.com/ggfevans/endorse/internal/util"
)

// typingDots is a custom spinner for the typing indicator.
var typingDots = spinner.Spinner{
	Frames: []string{"●··", "·●·", "··●", "·●·"},
	FPS:    300 * time.Millisecond,
}

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
	typingName     string        // who is typing ("" = nobody)
	typingSpinner  spinner.Model // animation driver
	composeView    string        // pre-rendered compose view
	hasCompose     bool          // whether compose is embedded
}

// New creates a new thread model.
func New(s styles.Styles) Model {
	vp := viewport.New(0, 0)
	sp := spinner.New()
	sp.Spinner = typingDots
	sp.Style = lipgloss.NewStyle().Foreground(s.Theme.Accent)
	return Model{styles: s, viewport: vp, typingSpinner: sp}
}

// SetComposeView sets the pre-rendered compose view for embedded rendering.
func (m *Model) SetComposeView(view string) {
	m.composeView = view
	m.hasCompose = true
}

// SetSize updates dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	contentWidth := w - 4 // border + padding
	if contentWidth < 1 {
		contentWidth = 1
	}
	composeH := 0
	if m.conversationID != "" && m.hasCompose {
		composeH = 4 // textarea (3 lines) + divider (1 line)
	}
	visibleH := h - 2 - 1 - composeH // border + title line - compose
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
	m.typingName = ""
	m.refreshContent()
	m.viewport.GotoTop()
}

// SetTyping shows the typing indicator and starts animation.
// Returns a Cmd to start the spinner ticker.
func (m *Model) SetTyping(name string) tea.Cmd {
	wasTyping := m.typingName != ""
	m.typingName = name
	m.refreshContent()
	if !wasTyping {
		m.viewport.GotoBottom()
	}
	return m.typingSpinner.Tick
}

// ClearTyping hides the typing indicator.
func (m *Model) ClearTyping() {
	if m.typingName == "" {
		return
	}
	m.typingName = ""
	m.refreshContent()
}

// IsTyping returns whether a typing indicator is active.
func (m Model) IsTyping() bool {
	return m.typingName != ""
}

// Update handles spinner tick messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.typingName == "" {
		return m, nil
	}
	var cmd tea.Cmd
	m.typingSpinner, cmd = m.typingSpinner.Update(msg)
	m.refreshContent()
	return m, cmd
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
		// No conversation selected — show placeholder, no compose
		placeholder := "\n" + m.styles.Muted.Render("  Select a conversation")
		innerHeight := m.height - 2
		content := title + "\n" + placeholder
		content = util.PadToHeight(content, innerHeight)
		return border.Width(m.width - 2).Render(content)
	}

	// Build content: title + viewport + compose
	innerHeight := m.height - 2 // subtract top/bottom border
	composeSection := ""
	composeH := 0
	if m.hasCompose {
		composeDivider := m.styles.Muted.Render(strings.Repeat("─", contentWidth))
		composeSection = "\n" + composeDivider + "\n" + m.composeView
		composeH = 4 // divider + textarea
	}

	content := title + "\n" + m.viewport.View()
	content = util.PadToHeight(content, innerHeight-composeH)
	content += composeSection
	return border.Width(m.width - 2).Render(content)
}

// refreshContent rebuilds the viewport content string from messages.
func (m *Model) refreshContent() {
	if len(m.messages) == 0 {
		if m.conversationID != "" {
			m.viewport.SetContent(m.styles.Muted.Render("  No messages"))
		} else {
			m.viewport.SetContent("")
		}
		return
	}

	contentWidth := m.width - 4
	if contentWidth < 1 {
		contentWidth = 1
	}

	accentBar := lipgloss.NewStyle().Foreground(m.styles.Theme.OwnSender).Render("▎")
	divider := m.styles.Muted.Render(strings.Repeat("─", contentWidth-2))

	// Word-wrap style for body text (account for prefix character)
	bodyWidth := contentWidth - 2
	if bodyWidth < 1 {
		bodyWidth = 1
	}
	wrapStyle := lipgloss.NewStyle().Width(bodyWidth)

	// Skip first sender header if it matches the conversation title (1:1 chat)
	skipFirstSender := len(m.messages) > 0 && m.messages[0].Sender == m.subject

	var lines []string
	var prevSender string
	for _, msg := range m.messages {
		if msg.Sender != prevSender {
			if prevSender == "" && skipFirstSender {
				// First sender matches title — skip redundant header
				prevSender = msg.Sender
			} else {
				if prevSender != "" {
					lines = append(lines, divider)
				}

				senderStyle := m.styles.SenderName
				if msg.IsOwn {
					senderStyle = m.styles.OwnSenderName
				}

				header := senderStyle.Render(msg.Sender)
				lines = append(lines, header)
				prevSender = msg.Sender
			}
		}

		prefix := " "
		if msg.IsOwn {
			prefix = accentBar
		}
		wrapped := wrapStyle.Render(msg.Body)
		for i, line := range strings.Split(wrapped, "\n") {
			if i == 0 {
				lines = append(lines, prefix+line)
			} else {
				lines = append(lines, " "+line)
			}
		}
		lines = append(lines, " "+m.styles.Timestamp.Render(msg.Timestamp))
	}

	if m.typingName != "" {
		if len(lines) > 0 {
			lines = append(lines, divider)
		}
		typingText := wrapStyle.Render(m.styles.Muted.Render(m.typingName+" is typing ") + m.typingSpinner.View())
		for _, line := range strings.Split(typingText, "\n") {
			lines = append(lines, " "+line)
		}
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
