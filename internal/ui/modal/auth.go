package modal

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ggfevans/endorse/internal/ui/styles"
)

// AuthSubmitMsg is sent when the user submits auth credentials.
type AuthSubmitMsg struct {
	Cookie       string
	PageInstance string
	XLiTrack     string
}

// AuthModel is the auth form modal.
type AuthModel struct {
	styles     styles.Styles
	width      int
	height     int
	inputs     []textinput.Model
	focusIndex int
	err        string
	loading    bool
}

const (
	authInputCookie = iota
	authInputCount
)

// NewAuth creates a new auth modal.
func NewAuth(s styles.Styles) AuthModel {
	inputs := make([]textinput.Model, authInputCount)

	inputs[authInputCookie] = textinput.New()
	inputs[authInputCookie].Placeholder = "Paste full Cookie header or just li_at value"
	inputs[authInputCookie].Focus()
	inputs[authInputCookie].CharLimit = 0 // unlimited
	inputs[authInputCookie].Width = 60

	return AuthModel{
		styles: s,
		inputs: inputs,
	}
}

// SetSize updates the modal dimensions.
func (m *AuthModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	inputW := w - 12
	if inputW < 30 {
		inputW = 30
	}
	if inputW > 80 {
		inputW = 80
	}
	for i := range m.inputs {
		m.inputs[i].Width = inputW
	}
}

// SetError sets an error message.
func (m *AuthModel) SetError(err string) {
	m.err = err
	m.loading = false
}

// SetLoading sets loading state.
func (m *AuthModel) SetLoading(loading bool) {
	m.loading = loading
}

// SetStyles updates styles.
func (m *AuthModel) SetStyles(s styles.Styles) {
	m.styles = s
}

// Update handles tea messages.
func (m AuthModel) Update(msg tea.Msg) (AuthModel, tea.Cmd) {
	if m.loading {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.focusIndex = (m.focusIndex + 1) % authInputCount
			return m, m.updateFocus()
		case "shift+tab", "up":
			m.focusIndex = (m.focusIndex - 1 + authInputCount) % authInputCount
			return m, m.updateFocus()
		case "enter":
			cookie := strings.TrimSpace(m.inputs[authInputCookie].Value())
			if cookie == "" {
				m.err = "Cookie header is required"
				return m, nil
			}
			m.err = ""
			m.loading = true
			return m, func() tea.Msg {
				return AuthSubmitMsg{
					Cookie: cookie,
				}
			}
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m AuthModel) updateFocus() tea.Cmd {
	cmds := make([]tea.Cmd, authInputCount)
	for i := range m.inputs {
		if i == m.focusIndex {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

// View renders the auth modal.
func (m AuthModel) View() string {
	titleStyle := m.styles.AccentText.Align(lipgloss.Center)
	mutedStyle := m.styles.Muted

	var b strings.Builder
	b.WriteString(titleStyle.Render("Endorse Authentication"))
	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("Extract cookies from your browser DevTools:"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("1. Open LinkedIn in your browser, open DevTools (F12)"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("2. Go to Network tab, click any request to linkedin.com"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("3. In Request Headers, right-click the Cookie value -> Copy"))
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render("   (Or: paste full header, or just the li_at=... value)"))
	b.WriteString("\n\n")

	labelStyle := lipgloss.NewStyle().Foreground(m.styles.Theme.Accent).Bold(true)
	b.WriteString(labelStyle.Render("Cookie Header:"))
	b.WriteString("\n")
	b.WriteString(m.inputs[authInputCookie].View())
	b.WriteString("\n\n")

	if m.err != "" {
		errStyle := lipgloss.NewStyle().Foreground(m.styles.Theme.Error)
		b.WriteString(errStyle.Render("Error: " + m.err))
		b.WriteString("\n\n")
	}

	if m.loading {
		b.WriteString(m.styles.AccentText.Render("Validating credentials..."))
	} else {
		b.WriteString(mutedStyle.Render("Press Enter to authenticate"))
	}

	contentWidth := m.width - 8
	if contentWidth < 40 {
		contentWidth = 40
	}
	if contentWidth > 90 {
		contentWidth = 90
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Theme.Accent).
		Padding(1, 3).
		Width(contentWidth)

	box := boxStyle.Render(b.String())

	// Center the box
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
