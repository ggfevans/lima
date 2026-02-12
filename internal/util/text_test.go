package util

import (
	"strings"
	"testing"
)

func TestInitials(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "two words", input: "John Doe", expected: "JD"},
		{name: "single lowercase word", input: "alice", expected: "A"},
		{name: "empty string", input: "", expected: "?"},
		{name: "three words takes first two", input: "First Middle Last", expected: "FM"},
		{name: "whitespace only", input: "   ", expected: "?"},
		{name: "single uppercase word", input: "Bob", expected: "B"},
		{name: "already uppercase", input: "JOHN DOE", expected: "JD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Initials(tt.input)
			if got != tt.expected {
				t.Errorf("Initials(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{name: "short string unchanged", input: "hello", maxLen: 10, expected: "hello"},
		{name: "exact length unchanged", input: "hello", maxLen: 5, expected: "hello"},
		{name: "long string truncated with ellipsis", input: "hello world", maxLen: 8, expected: "hello..."},
		{name: "maxLen zero returns empty", input: "hello", maxLen: 0, expected: ""},
		{name: "negative maxLen returns empty", input: "hello", maxLen: -1, expected: ""},
		{name: "maxLen 1 no ellipsis", input: "hello", maxLen: 1, expected: "h"},
		{name: "maxLen 2 no ellipsis", input: "hello", maxLen: 2, expected: "he"},
		{name: "maxLen 3 no ellipsis", input: "hello", maxLen: 3, expected: "hel"},
		{name: "maxLen 4 with ellipsis", input: "hello world", maxLen: 4, expected: "h..."},
		{name: "empty string", input: "", maxLen: 5, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.input, tt.maxLen)
			if got != tt.expected {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.expected)
			}
		})
	}
}

func TestPadToHeight(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		height         int
		expectedLines  int
		checkContent   bool
		expectedOutput string
	}{
		{
			name:          "content shorter than height gets padded",
			content:       "line1\nline2",
			height:        4,
			expectedLines: 4,
		},
		{
			name:          "content taller than height gets truncated",
			content:       "line1\nline2\nline3\nline4",
			height:        2,
			expectedLines: 2,
		},
		{
			name:           "height zero returns empty",
			content:        "line1\nline2",
			height:         0,
			checkContent:   true,
			expectedOutput: "",
		},
		{
			name:           "negative height returns empty",
			content:        "line1",
			height:         -1,
			checkContent:   true,
			expectedOutput: "",
		},
		{
			name:           "exact height unchanged",
			content:        "line1\nline2\nline3",
			height:         3,
			expectedLines:  3,
			checkContent:   true,
			expectedOutput: "line1\nline2\nline3",
		},
		{
			name:          "single line padded",
			content:       "only",
			height:        3,
			expectedLines: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PadToHeight(tt.content, tt.height)

			if tt.checkContent {
				if got != tt.expectedOutput {
					t.Errorf("PadToHeight(%q, %d) = %q, want %q", tt.content, tt.height, got, tt.expectedOutput)
				}
				return
			}

			lines := strings.Split(got, "\n")
			if len(lines) != tt.expectedLines {
				t.Errorf("PadToHeight(%q, %d) produced %d lines, want %d", tt.content, tt.height, len(lines), tt.expectedLines)
			}
		})
	}
}
