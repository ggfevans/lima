package convlist

import (
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

func newTestConvList() Model {
	theme := config.ThemeByName("")
	s := styles.New(theme)
	return New(s)
}

func sampleConversations() []Conversation {
	return []Conversation{
		{ID: "1", Name: "Alice Johnson", LastMessage: "Hey there!", Timestamp: "10:30", Unread: true, UnreadCount: 2},
		{ID: "2", Name: "Bob Smith", LastMessage: "See you tomorrow", Timestamp: "09:15", Unread: false, UnreadCount: 0},
		{ID: "3", Name: "Carol White", LastMessage: "Thanks!", Timestamp: "Yesterday", Unread: true, UnreadCount: 1},
	}
}

func TestDimension(t *testing.T) {
	m := newTestConvList()
	m.SetSize(30, 20)
	output := m.View()
	lines := countRenderedLines(output)
	if lines != 20 {
		t.Errorf("expected 20 lines, got %d", lines)
	}
}

func TestEmptyState(t *testing.T) {
	m := newTestConvList()
	m.SetSize(30, 20)

	output := stripAnsi(m.View())
	if !strings.Contains(output, "No conversations") {
		t.Errorf("expected 'No conversations' in empty state, got:\n%s", output)
	}
}

func TestContent(t *testing.T) {
	m := newTestConvList()
	m.SetSize(30, 20)
	m.SetConversations(sampleConversations())

	output := stripAnsi(m.View())
	for _, name := range []string{"Alice Johnson", "Bob Smith", "Carol White"} {
		if !strings.Contains(output, name) {
			t.Errorf("expected view to contain %q, got:\n%s", name, output)
		}
	}
}

func TestNavigation(t *testing.T) {
	m := newTestConvList()
	m.SetSize(30, 20)
	m.SetConversations(sampleConversations())

	// Initial selection is 0
	if m.Selected() != 0 {
		t.Errorf("expected initial Selected()=0, got %d", m.Selected())
	}

	// MoveDown
	m.MoveDown()
	if m.Selected() != 1 {
		t.Errorf("expected Selected()=1 after MoveDown, got %d", m.Selected())
	}

	// MoveDown again
	m.MoveDown()
	if m.Selected() != 2 {
		t.Errorf("expected Selected()=2 after second MoveDown, got %d", m.Selected())
	}

	// MoveDown at end should not exceed bounds
	m.MoveDown()
	if m.Selected() != 2 {
		t.Errorf("expected Selected()=2 at upper bound, got %d", m.Selected())
	}

	// MoveUp
	m.MoveUp()
	if m.Selected() != 1 {
		t.Errorf("expected Selected()=1 after MoveUp, got %d", m.Selected())
	}

	// MoveUp to 0
	m.MoveUp()
	if m.Selected() != 0 {
		t.Errorf("expected Selected()=0 after second MoveUp, got %d", m.Selected())
	}

	// MoveUp at 0 should not go below 0
	m.MoveUp()
	if m.Selected() != 0 {
		t.Errorf("expected Selected()=0 at lower bound, got %d", m.Selected())
	}

	// MoveToBottom
	m.MoveToBottom()
	if m.Selected() != 2 {
		t.Errorf("expected Selected()=2 after MoveToBottom, got %d", m.Selected())
	}

	// MoveToTop
	m.MoveToTop()
	if m.Selected() != 0 {
		t.Errorf("expected Selected()=0 after MoveToTop, got %d", m.Selected())
	}
}

func TestSelection(t *testing.T) {
	m := newTestConvList()
	m.SetSize(30, 20)
	convs := sampleConversations()
	m.SetConversations(convs)

	// Initially selects first conversation
	c, ok := m.SelectedConversation()
	if !ok {
		t.Fatal("expected SelectedConversation to return ok=true")
	}
	if c.ID != "1" {
		t.Errorf("expected selected conversation ID='1', got %q", c.ID)
	}

	// After MoveDown, selects second
	m.MoveDown()
	c, ok = m.SelectedConversation()
	if !ok {
		t.Fatal("expected SelectedConversation to return ok=true")
	}
	if c.ID != "2" {
		t.Errorf("expected selected conversation ID='2', got %q", c.ID)
	}

	// After MoveToBottom, selects last
	m.MoveToBottom()
	c, ok = m.SelectedConversation()
	if !ok {
		t.Fatal("expected SelectedConversation to return ok=true")
	}
	if c.ID != "3" {
		t.Errorf("expected selected conversation ID='3', got %q", c.ID)
	}

	// Empty list returns ok=false
	m2 := newTestConvList()
	_, ok = m2.SelectedConversation()
	if ok {
		t.Error("expected SelectedConversation to return ok=false for empty list")
	}
}

func TestCountAndUnreadCount(t *testing.T) {
	m := newTestConvList()
	m.SetConversations(sampleConversations())

	if m.Count() != 3 {
		t.Errorf("expected Count()=3, got %d", m.Count())
	}

	// sampleConversations has 2 unread (Alice and Carol)
	if m.UnreadCount() != 2 {
		t.Errorf("expected UnreadCount()=2, got %d", m.UnreadCount())
	}
}
