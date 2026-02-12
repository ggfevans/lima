package util

import (
	"fmt"
	"time"
)

// RelativeTime formats a timestamp into a human-friendly relative string.
func RelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		m := int(diff.Minutes())
		if m == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", m)
	case isToday(t, now):
		return t.Format("3:04 PM")
	case isYesterday(t, now):
		return "Yesterday"
	case diff < 7*24*time.Hour:
		return t.Format("Monday")
	case t.Year() == now.Year():
		return t.Format("Jan 2")
	default:
		return t.Format("Jan 2, 2006")
	}
}

func isToday(t, now time.Time) bool {
	ty, tm, td := t.Date()
	ny, nm, nd := now.Date()
	return ty == ny && tm == nm && td == nd
}

func isYesterday(t, now time.Time) bool {
	yesterday := now.AddDate(0, 0, -1)
	return isToday(t, yesterday)
}
