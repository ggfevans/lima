package util

import (
	"strings"
)

// Initials extracts up to 2 initials from a name.
func Initials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "?"
	}

	var initials []rune
	for _, p := range parts {
		if len(initials) >= 2 {
			break
		}
		for _, r := range p {
			initials = append(initials, r)
			break
		}
	}
	return strings.ToUpper(string(initials))
}

// Truncate truncates a string to maxLen, adding ellipsis if needed.
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// PadToHeight pads or truncates content to exactly the given number of lines.
func PadToHeight(content string, height int) string {
	if height <= 0 {
		return ""
	}
	lines := strings.Split(content, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
