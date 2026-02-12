package sidebar

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

func newTestSidebar() Model {
	theme := config.ThemeByName("")
	s := styles.New(theme)
	return New(s)
}

func TestDimension(t *testing.T) {
	m := newTestSidebar()
	m.SetSize(14, 20)
	output := m.View()
	lines := countRenderedLines(output)
	if lines != 20 {
		t.Errorf("expected 20 lines, got %d", lines)
	}
}

func TestNavigation(t *testing.T) {
	m := newTestSidebar()
	m.SetSize(14, 20)

	// Initial selection is 0
	if m.Selected() != 0 {
		t.Errorf("expected initial Selected()=0, got %d", m.Selected())
	}

	// MoveDown increases Selected
	m.MoveDown()
	if m.Selected() != 1 {
		t.Errorf("expected Selected()=1 after MoveDown, got %d", m.Selected())
	}

	// MoveDown again
	m.MoveDown()
	if m.Selected() != 2 {
		t.Errorf("expected Selected()=2 after second MoveDown, got %d", m.Selected())
	}

	// MoveDown at max should not go above 2
	m.MoveDown()
	if m.Selected() != 2 {
		t.Errorf("expected Selected()=2 at upper bound, got %d", m.Selected())
	}

	// MoveUp decreases Selected
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
}

func TestFocus(t *testing.T) {
	m := newTestSidebar()
	m.SetSize(14, 20)

	// Initially not focused
	if m.Focused() {
		t.Error("expected Focused()=false initially")
	}

	// Focus sets Focused to true
	m.Focus()
	if !m.Focused() {
		t.Error("expected Focused()=true after Focus()")
	}

	// Blur sets Focused to false
	m.Blur()
	if m.Focused() {
		t.Error("expected Focused()=false after Blur()")
	}

	// Toggle back to focused
	m.Focus()
	if !m.Focused() {
		t.Error("expected Focused()=true after second Focus()")
	}

	// Verify the View can be rendered in both states without error
	m.Focus()
	focusedView := m.View()
	if focusedView == "" {
		t.Error("expected non-empty focused view")
	}

	m.Blur()
	blurredView := m.View()
	if blurredView == "" {
		t.Error("expected non-empty blurred view")
	}
}

func TestCounts(t *testing.T) {
	m := newTestSidebar()
	m.SetSize(14, 20)
	m.SetCounts(5, 3)

	output := stripAnsi(m.View())
	if !strings.Contains(output, "Inbox 5") {
		t.Errorf("expected view to contain 'Inbox 5', got:\n%s", output)
	}
}

func TestContent(t *testing.T) {
	m := newTestSidebar()
	m.SetSize(14, 20)

	output := stripAnsi(m.View())

	expected := []string{"FOLDERS", "Inbox", "Unread", "Archived"}
	for _, s := range expected {
		if !strings.Contains(output, s) {
			t.Errorf("expected view to contain %q, got:\n%s", s, output)
		}
	}
}
