package thread

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/ggfevans/li-cli/internal/config"
	"github.com/ggfevans/li-cli/internal/ui/styles"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripAnsi(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func countRenderedLines(s string) int {
	return len(strings.Split(s, "\n"))
}

func newTestThread() Model {
	theme := config.ThemeByName("")
	s := styles.New(theme)
	return New(s)
}

func sampleMessages() []Message {
	return []Message{
		{ID: "m1", Sender: "Alice Johnson", Body: "Hello there!", Timestamp: "10:30 AM", IsOwn: false},
		{ID: "m2", Sender: "Me", Body: "Hi Alice, how are you?", Timestamp: "10:31 AM", IsOwn: true},
		{ID: "m3", Sender: "Alice Johnson", Body: "Doing great, thanks!", Timestamp: "10:32 AM", IsOwn: false},
	}
}

func TestDimension(t *testing.T) {
	m := newTestThread()
	m.SetSize(60, 20)
	output := m.View()
	lines := countRenderedLines(output)
	if lines != 20 {
		t.Errorf("expected 20 lines, got %d", lines)
	}
}

func TestEmptyPlaceholder(t *testing.T) {
	m := newTestThread()
	m.SetSize(60, 20)

	output := stripAnsi(m.View())
	if !strings.Contains(output, "Select a conversation") {
		t.Errorf("expected 'Select a conversation' in empty state, got:\n%s", output)
	}
}

func TestNoMessages(t *testing.T) {
	m := newTestThread()
	m.SetSize(60, 20)
	m.SetConversation("conv-1", "Test Subject")

	output := stripAnsi(m.View())
	if !strings.Contains(output, "No messages") {
		t.Errorf("expected 'No messages' when conversation has no messages, got:\n%s", output)
	}
}

func TestContent(t *testing.T) {
	m := newTestThread()
	m.SetSize(60, 20)
	m.SetConversation("conv-1", "Test Subject")
	m.SetMessages(sampleMessages())

	output := stripAnsi(m.View())
	if !strings.Contains(output, "Alice Johnson") {
		t.Errorf("expected view to contain sender name 'Alice Johnson', got:\n%s", output)
	}
	if !strings.Contains(output, "Hello there!") {
		t.Errorf("expected view to contain message body 'Hello there!', got:\n%s", output)
	}
}

func TestScrolling(t *testing.T) {
	m := newTestThread()
	m.SetSize(60, 20)
	m.SetConversation("conv-1", "Chat")

	// Create many messages so they overflow the visible area
	var msgs []Message
	for i := 0; i < 30; i++ {
		msgs = append(msgs, Message{
			ID:        fmt.Sprintf("m%d", i),
			Sender:    fmt.Sprintf("User %d", i),
			Body:      fmt.Sprintf("Message number %d with some content", i),
			Timestamp: "12:00",
			IsOwn:     i%2 == 0,
		})
	}
	m.SetMessages(msgs)

	// SetMessages scrolls to bottom, so AtTop should be false
	if m.AtTop() {
		t.Error("expected AtTop()=false after SetMessages with many messages (scrolled to bottom)")
	}

	// ScrollUp enough to reach top
	m.ScrollUp(1000)
	if !m.AtTop() {
		t.Error("expected AtTop()=true after scrolling up fully")
	}

	// ScrollDown moves away from top
	m.ScrollDown(5)
	if m.AtTop() {
		t.Error("expected AtTop()=false after ScrollDown")
	}
}

func TestScrollPercent(t *testing.T) {
	m := newTestThread()
	m.SetSize(60, 20)
	m.SetConversation("conv-1", "Chat")

	var msgs []Message
	for i := 0; i < 50; i++ {
		msgs = append(msgs, Message{
			ID:        fmt.Sprintf("m%d", i),
			Sender:    fmt.Sprintf("User %d", i),
			Body:      fmt.Sprintf("Message %d", i),
			Timestamp: "12:00",
			IsOwn:     false,
		})
	}
	m.SetMessages(msgs)

	// At bottom, scroll percent should be 1.0
	pct := m.ScrollPercent()
	if pct < 0.99 {
		t.Errorf("expected ScrollPercent()~1.0 at bottom, got %f", pct)
	}

	// At top, scroll percent should be 0.0
	m.ScrollUp(10000)
	pct = m.ScrollPercent()
	if pct > 0.01 {
		t.Errorf("expected ScrollPercent()~0.0 at top, got %f", pct)
	}
}

func TestClear(t *testing.T) {
	m := newTestThread()
	m.SetSize(60, 20)
	m.SetConversation("conv-1", "Test Subject")
	m.SetMessages(sampleMessages())

	if !m.HasConversation() {
		t.Error("expected HasConversation()=true after SetConversation")
	}
	if m.MessageCount() != 3 {
		t.Errorf("expected MessageCount()=3, got %d", m.MessageCount())
	}

	m.Clear()

	if m.HasConversation() {
		t.Error("expected HasConversation()=false after Clear()")
	}
	if m.ConversationID() != "" {
		t.Errorf("expected ConversationID()='', got %q", m.ConversationID())
	}
	if m.MessageCount() != 0 {
		t.Errorf("expected MessageCount()=0 after Clear(), got %d", m.MessageCount())
	}
}
