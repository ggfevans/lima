package app

import tea "github.com/charmbracelet/bubbletea"

// isQuitKey returns true for quit keybindings.
func isQuitKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "q", "ctrl+c":
		return true
	}
	return false
}

// isTabKey returns true for focus-cycle keys.
func isTabKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "tab", "right":
		return true
	}
	return false
}

// isShiftTabKey returns true for reverse focus-cycle keys.
func isShiftTabKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "shift+tab", "left":
		return true
	}
	return false
}

// isMarkReadKey returns true for mark-as-read/unread toggle.
func isMarkReadKey(msg tea.KeyMsg) bool {
	return msg.String() == "m"
}

// isDeleteKey returns true for delete action.
func isDeleteKey(msg tea.KeyMsg) bool {
	return msg.String() == "d"
}

// isDownKey returns true for downward navigation.
func isDownKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "j", "down":
		return true
	}
	return false
}

// isUpKey returns true for upward navigation.
func isUpKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "k", "up":
		return true
	}
	return false
}

// isEnterKey returns true for selection/confirm.
func isEnterKey(msg tea.KeyMsg) bool {
	return msg.String() == "enter"
}

// isEscapeKey returns true for cancel/back.
func isEscapeKey(msg tea.KeyMsg) bool {
	return msg.String() == "esc"
}

// isHelpKey returns true for help toggle.
func isHelpKey(msg tea.KeyMsg) bool {
	return msg.String() == "?"
}

// isReplyKey returns true for reply action.
func isReplyKey(msg tea.KeyMsg) bool {
	return msg.String() == "r"
}

// isFilterKey returns true for filter toggle.
func isFilterKey(msg tea.KeyMsg) bool {
	return msg.String() == "f"
}

// isSendKey returns true for message send (Ctrl+Enter in compose).
func isSendKey(msg tea.KeyMsg) bool {
	return msg.String() == "ctrl+s"
}

// isTopKey returns true for jump-to-top.
func isTopKey(msg tea.KeyMsg) bool {
	return msg.String() == "g"
}

// isBottomKey returns true for jump-to-bottom.
func isBottomKey(msg tea.KeyMsg) bool {
	return msg.String() == "G"
}

// isPageDown returns true for page-down scroll.
func isPageDown(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "ctrl+d", "pgdown":
		return true
	}
	return false
}

// isPageUp returns true for page-up scroll.
func isPageUp(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "ctrl+u", "pgup":
		return true
	}
	return false
}
