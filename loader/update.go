package loader

import (
	"time"

	"github.com/gabrielluizsf/sge/game"
)

// Update runs the most urgent jobs until the frame's budget is spent, at
// least one a frame, so a job longer than the budget still gets done.
func (l *Loader) Update(_ *game.Object, _ float64) {
	budget := l.Budget
	if budget <= 0 {
		budget = 4 * time.Millisecond
	}
	start := time.Now()
	for len(l.queue) > 0 {
		l.runNext()
		if time.Since(start) >= budget {
			return
		}
	}
}

// Drain runs every job waiting, as a test or a tool would, not a frame.
func (l *Loader) Drain() {
	for len(l.queue) > 0 {
		l.runNext()
	}
}

// runNext runs the most urgent job, keeps what it made, and lets go of
// the least lately used assets while over the memory budget.
func (l *Loader) runNext() {
	best := 0
	for i, q := range l.queue {
		b := l.queue[best]
		if q.priority < b.priority || q.priority == b.priority && q.seq < b.seq {
			best = i
		}
	}
	q := l.queue[best]
	l.queue = append(l.queue[:best], l.queue[best+1:]...)
	v, err := q.job.Make()
	l.clock++
	l.assets[q.job.Name] = &asset{value: v, err: err, bytes: q.job.Bytes, used: l.clock, group: q.job.Group}
	l.used += q.job.Bytes
	l.groups[q.job.Group].done++
	l.evict(q.job.Name)
}

// evict lets go of the least lately used assets, but keep, while what is
// kept holds more than the memory budget.
func (l *Loader) evict(keep string) {
	for l.Memory > 0 && l.used > l.Memory {
		oldest := ""
		for name, a := range l.assets {
			if name != keep && (oldest == "" || a.used < l.assets[oldest].used || a.used == l.assets[oldest].used && name < oldest) {
				oldest = name
			}
		}
		if oldest == "" {
			return
		}
		a := l.assets[oldest]
		l.used -= a.bytes
		delete(l.assets, oldest)
		if g := l.groups[a.group]; g != nil {
			g.done--
			g.asked--
		}
	}
}
