package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModel(t *testing.T) {
	m := New(Options{})
	// New() should return a valid Model without panicking.
	// The initial state should be either StateAuth or StateLoading depending
	// on whether credentials are stored on disk.
	if m.state != StateAuth && m.state != StateLoading {
		t.Errorf("expected initial state to be StateAuth or StateLoading, got %d", m.state)
	}
}

func TestModelInit(t *testing.T) {
	m := New(Options{})
	// Init() should return without panicking. The returned command may be nil
	// (no stored credentials) or a validation command (stored credentials).
	_ = m.Init()
}

func TestModelView_BeforeReady(t *testing.T) {
	m := New(Options{})
	// Before receiving a WindowSizeMsg, ready is false, so View() should
	// return the loading placeholder.
	got := m.View()
	if got != "Loading..." {
		t.Errorf("expected View() before ready to return %q, got %q", "Loading...", got)
	}
}

func TestModelUpdate_WindowSize(t *testing.T) {
	m := New(Options{})
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	result, _ := m.Update(msg)
	// After processing a WindowSizeMsg, the model should be ready and
	// View() should no longer return the bare "Loading..." placeholder.
	view := result.View()
	if view == "Loading..." {
		t.Error("expected View() after WindowSizeMsg to differ from the loading placeholder")
	}
	if view == "" {
		t.Error("expected View() after WindowSizeMsg to return non-empty content")
	}
}

func TestAppStateConstants(t *testing.T) {
	// Verify the iota-generated constants have the expected values.
	if StateAuth != 0 {
		t.Errorf("expected StateAuth = 0, got %d", StateAuth)
	}
	if StateLoading != 1 {
		t.Errorf("expected StateLoading = 1, got %d", StateLoading)
	}
	if StateMessaging != 2 {
		t.Errorf("expected StateMessaging = 2, got %d", StateMessaging)
	}

	// Ensure all values are distinct.
	states := map[AppState]string{
		StateAuth:      "StateAuth",
		StateLoading:   "StateLoading",
		StateMessaging: "StateMessaging",
	}
	if len(states) != 3 {
		t.Error("expected 3 distinct AppState values")
	}
}

func TestFocusPanelConstants(t *testing.T) {
	// Verify the iota-generated constants have the expected values.
	if FocusSidebar != 0 {
		t.Errorf("expected FocusSidebar = 0, got %d", FocusSidebar)
	}
	if FocusConvList != 1 {
		t.Errorf("expected FocusConvList = 1, got %d", FocusConvList)
	}
	if FocusThread != 2 {
		t.Errorf("expected FocusThread = 2, got %d", FocusThread)
	}
	if FocusCompose != 3 {
		t.Errorf("expected FocusCompose = 3, got %d", FocusCompose)
	}

	// Ensure all values are distinct.
	panels := map[FocusedPanel]string{
		FocusSidebar:  "FocusSidebar",
		FocusConvList: "FocusConvList",
		FocusThread:   "FocusThread",
		FocusCompose:  "FocusCompose",
	}
	if len(panels) != 4 {
		t.Error("expected 4 distinct FocusedPanel values")
	}
}
