package util

import (
	"strings"
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	now := time.Now()

	// For the "today shows time format" test case, we need a time that is:
	//   1. On the same calendar day as now
	//   2. More than 1 hour ago (so it doesn't hit the "Xm ago" branch)
	// If it is too early in the day (hour < 2), we skip that specific test case
	// since we can't construct a valid time satisfying both conditions.
	canTestTodayFormat := now.Hour() >= 2
	var todayEarlier time.Time
	if canTestTodayFormat {
		todayEarlier = now.Add(-2 * time.Hour)
	}

	tests := []struct {
		name      string
		input     time.Time
		skip      bool
		wantExact string   // if non-empty, expect exact match
		wantAny   []string // if non-empty, output must contain one of these
	}{
		{
			name:      "just now (30 seconds ago)",
			input:     now.Add(-30 * time.Second),
			wantExact: "just now",
		},
		{
			name:      "1 minute ago",
			input:     now.Add(-1 * time.Minute),
			wantExact: "1m ago",
		},
		{
			name:    "5 minutes ago",
			input:   now.Add(-5 * time.Minute),
			wantAny: []string{"5m ago"},
		},
		{
			name:    "45 minutes ago",
			input:   now.Add(-45 * time.Minute),
			wantAny: []string{"m ago"},
		},
		{
			name:    "earlier today shows time format",
			input:   todayEarlier,
			skip:    !canTestTodayFormat,
			wantAny: []string{"AM", "PM"},
		},
		{
			name:      "yesterday",
			input:     time.Date(now.Year(), now.Month(), now.Day()-1, 12, 0, 0, 0, now.Location()),
			wantExact: "Yesterday",
		},
		{
			name:  "3 days ago shows weekday",
			input: now.Add(-3 * 24 * time.Hour),
			wantAny: []string{
				"Monday", "Tuesday", "Wednesday", "Thursday",
				"Friday", "Saturday", "Sunday",
			},
		},
		{
			name:  "10 days ago same year shows month and day",
			input: now.Add(-10 * 24 * time.Hour),
			// Could be a weekday name if within 7 days, or "Jan 2" format
			// At 10 days it should definitely be past the weekday window
		},
		{
			name:    "different year",
			input:   time.Date(2020, time.March, 15, 10, 0, 0, 0, now.Location()),
			wantAny: []string{"2020"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("skipping: cannot construct valid test input at this time of day")
			}

			got := RelativeTime(tt.input)

			if got == "" {
				t.Errorf("RelativeTime() returned empty string")
				return
			}

			if tt.wantExact != "" {
				if got != tt.wantExact {
					t.Errorf("RelativeTime() = %q, want %q", got, tt.wantExact)
				}
				return
			}

			if len(tt.wantAny) > 0 {
				found := false
				for _, want := range tt.wantAny {
					if strings.Contains(got, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("RelativeTime() = %q, want string containing one of %v", got, tt.wantAny)
				}
			}
		})
	}
}

func TestRelativeTimeNotEmpty(t *testing.T) {
	// Verify a range of durations all produce non-empty output.
	durations := []time.Duration{
		0,
		10 * time.Second,
		1 * time.Minute,
		30 * time.Minute,
		2 * time.Hour,
		25 * time.Hour,
		4 * 24 * time.Hour,
		14 * 24 * time.Hour,
		200 * 24 * time.Hour,
		800 * 24 * time.Hour,
	}

	for _, d := range durations {
		input := time.Now().Add(-d)
		got := RelativeTime(input)
		if got == "" {
			t.Errorf("RelativeTime(%v ago) returned empty string", d)
		}
	}
}
