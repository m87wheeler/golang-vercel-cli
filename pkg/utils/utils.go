package utils

import (
	"fmt"
	"time"
)

func ElapsedTime(sinceUnix int64) string {
	// Convert the Unix timestamp to a time.Time
	since := time.Unix(sinceUnix, 0)

	// Get the current time
	now := time.Now()

	// Calculate the difference in minutes, hours, and days
	diff := now.Sub(since)
	seconds := int(diff.Seconds())
	minutes := int(diff.Minutes())
	hours := int(diff.Hours())
	days := int(diff.Hours() / 24)

	// Return a single unit based on the elapsed time
	if seconds < 60 {
		return fmt.Sprintf("%ds ago", seconds)
	} else if minutes < 60 {
		return fmt.Sprintf("%dm ago", minutes)
	} else if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	} else {
		return fmt.Sprintf("%dd ago", days)
	}
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ToggleState(states []string, state string) []string {
	// Check if the state exists in the slice
	for i, s := range states {
		if s == state {
			// State exists, remove it
			return append(states[:i], states[i+1:]...)
		}
	}

	// State doesn't exist, add it
	return append(states, state)
}
