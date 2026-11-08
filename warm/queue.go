// Package warm makes the pictures a game will show soon ahead of time, a
// few a frame, so that none is drawn for the first time in the frame that
// shows it: a person met in a town far off, whose frames are each drawn
// from SVG, would hold that frame up for as long as all of them take.
package warm

import (
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
)

// Queue is a script that makes the sprites added to it, in the order
// added, PerFrame of them each frame (one if that is not set), twice as
// many while more than Backlog (sixteen times PerFrame if that is not set)
// wait: a count and not a time, so a run is the same on every machine. A
// sprite already made costs nothing.
type Queue struct {
	PerFrame int
	Backlog  int
	waiting  []*sprite.Sprite
}

// Add puts the sprites to be made behind those waiting.
func (q *Queue) Add(s ...*sprite.Sprite) {
	for _, p := range s {
		if p != nil && !p.IsReady() {
			q.waiting = append(q.waiting, p)
		}
	}
}

// Waiting is how many sprites are still to be made.
func (q *Queue) Waiting() int { return len(q.waiting) }

// Update makes this frame's share.
func (q *Queue) Update(_ *game.Object, _ float64) {
	made, n := 0, 0
	share := max(q.PerFrame, 1)
	if len(q.waiting) > q.backlog(share) {
		share *= 2
	}
	for n < len(q.waiting) && made < share {
		if s := q.waiting[n]; !s.IsReady() {
			s.Ready()
			made++
		}
		q.waiting[n] = nil
		n++
	}
	q.waiting = q.waiting[n:]
}

// backlog is how many sprites may wait before the share doubles.
func (q *Queue) backlog(share int) int {
	if q.Backlog > 0 {
		return q.Backlog
	}
	return 16 * share
}
