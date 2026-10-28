//go:build gots_overlay

// Deadlines as the gotsNet bridge takes them.

package net

import "time"

// waitMillis is how long a call may wait for deadline d: -1 for no
// deadline, 0 for one passed, otherwise at least a millisecond.
func waitMillis(d time.Time) int {
	if d.IsZero() {
		return -1
	}
	left := time.Until(d)
	if left <= 0 {
		return 0
	}
	return int((left + time.Millisecond - 1) / time.Millisecond)
}

// passed reports whether the deadline d has passed.
func passed(d time.Time) bool {
	return !d.IsZero() && !time.Now().Before(d)
}
