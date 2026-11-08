// Package progress is an experience counter with levels: deeds add XP, and
// each level takes more than the one before. It is plain integer maths, so
// it plays the same on every platform.
package progress

import (
	"fmt"
	"math"
)

// MaxLevel is the last level; there is no XP past it.
const MaxLevel = 99

// Progress is a level and the XP gained toward the next one.
type Progress struct {
	Level int
	XP    int
}

// New is a beginning: level 1 and no XP.
func New() Progress { return Progress{Level: 1} }

// Need is the XP that takes level l to the next: 100 x l^1.5, rounded.
func Need(l int) int { return int(math.Round(100 * float64(l) * math.Sqrt(float64(l)))) }

// Add gains n XP, carrying what is over into the next levels, and says how
// many levels it gained. A negative n gains nothing.
func (p *Progress) Add(n int) int {
	gained := 0
	for n > 0 && p.Level < MaxLevel {
		room := Need(p.Level) - p.XP
		if n < room {
			p.XP += n
			return gained
		}
		n -= room
		p.Level++
		p.XP = 0
		gained++
	}
	return gained
}

// Valid is nil if p is a state Add can have made.
func (p Progress) Valid() error {
	if p.Level < 1 || p.Level > MaxLevel || p.XP < 0 || p.XP >= p.span() {
		return fmt.Errorf("progress: level %d with %d XP", p.Level, p.XP)
	}
	return nil
}

// span is how much XP p's level holds: what it needs, or nothing at all at
// the last level, where XP stays 0.
func (p Progress) span() int {
	if p.Level >= MaxLevel {
		return 1
	}
	return Need(p.Level)
}
