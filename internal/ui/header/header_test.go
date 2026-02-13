package header

import (
	"regexp"
	"strings"
	"testing"

	"github.com/ggfevans/lima/internal/config"
	"github.com/ggfevans/lima/internal/ui/styles"
)

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripAnsi(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

func countRenderedLines(s string) int {
	return len(strings.Split(s, "\n"))
}

func newTestHeader() Model {
	theme := config.ThemeByName("")
	s := styles.New(theme)
	return New(s)
}

func TestDefaultState(t *testing.T) {
	m := newTestHeader()
	m.SetWidth(80)

	output := stripAnsi(m.View())
	if !strings.Contains(output, "Li-CLI") {
		t.Errorf("expected view to contain 'Li-CLI', got:\n%s", output)
	}
	if !strings.Contains(output, "★") {
		t.Errorf("expected view to contain '★' indicator, got:\n%s", output)
	}
}

func TestConnected(t *testing.T) {
	m := newTestHeader()
	m.SetWidth(80)
	m.SetConnected(true)

	output := stripAnsi(m.View())
	if !strings.Contains(output, "★") {
		t.Errorf("expected view to contain '★' indicator, got:\n%s", output)
	}
}

func TestUsername(t *testing.T) {
	m := newTestHeader()
	m.SetWidth(80)
	m.SetUsername("testuser")

	output := stripAnsi(m.View())
	if !strings.Contains(output, "testuser") {
		t.Errorf("expected view to contain 'testuser', got:\n%s", output)
	}
}

func TestWidth(t *testing.T) {
	m := newTestHeader()
	m.SetWidth(80)

	output := m.View()
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one line in header output")
	}

	// Strip ANSI to measure visible width of first line
	firstLine := stripAnsi(lines[0])
	runeWidth := len([]rune(firstLine))

	// The rendered first line should be around the set width (lipgloss may
	// add padding but should not greatly exceed the target width).
	if runeWidth > 85 {
		t.Errorf("expected first line width <= ~85 (for width=80), got %d", runeWidth)
	}
	if runeWidth < 60 {
		t.Errorf("expected first line width >= 60 (for width=80), got %d", runeWidth)
	}
}
