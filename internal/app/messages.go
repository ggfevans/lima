package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AppState represents the top-level application state.
type AppState int

const (
	StateAuth     AppState = iota
	StateLoading
	StateMessaging
)

// FocusedPanel identifies which panel currently has focus.
type FocusedPanel int

const (
	FocusSidebar FocusedPanel = iota
	FocusConvList
	FocusThread
	FocusCompose
)

// --- Internal tea.Msg types (not produced by linkedin client) ---

// ProgramRefMsg carries the tea.Program reference for real-time event bridging.
type ProgramRefMsg struct {
	Program *tea.Program
}

// ClearErrorMsg clears the error display.
type ClearErrorMsg struct{}

// clearErrorAfter returns a command that clears errors after a delay.
func clearErrorAfter() tea.Cmd {
	return tea.Tick(5*time.Second, func(_ time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}
