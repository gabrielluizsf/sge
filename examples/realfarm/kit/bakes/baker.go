package bakes

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/tree"
)

// Baker bakes the farm's sprites onto a Shelf from the manifest's assets
// and the humans' pins. Workers is how many frames of a task are baked at
// once, one when 0, so a bake never takes more than a core; Pause is how
// long it rests between tasks, so it does not peg that one.
type Baker struct {
	Shelf    *Shelf
	Manifest kit.Manifest
	Humans   kit.Humans
	Workers  int
	Pause    time.Duration
	// Hold, while true, keeps Do from starting a task: the game is busy
	// with something that needs the core, a character being made.
	Hold atomic.Bool

	paletteOnce sync.Once
	palette     building.Palette
	paletteErr  error

	barkOnce sync.Once
	barks    tree.Barks
	barkErr  error

	makeOnce sync.Once
	makers   *human.Kit
	makeErr  error

	mu   sync.Mutex
	last stageKey
	body *stageBody

	saved time.Time
}

// Task is one thing to bake, named: of a Kind ("ground", "model", "crop",
// "water", "building" or "person") and an ID within it, Core when the
// game cannot start without it.
type Task struct {
	Name, Kind, ID string
	Core           bool
	Run            func() error
}

// saveEvery is how often the shelf's index is written between tasks.
const saveEvery = 3 * time.Second

// holdWait is how long Do waits between looks at Hold.
const holdWait = 50 * time.Millisecond

// Do runs t unless it ran to its end before, and records that it has;
// the shelf is saved now and then, and at the end by Save.
func (b *Baker) Do(t Task) error {
	if b.Shelf.Finished(t.Name) {
		return nil
	}
	for b.Hold.Load() {
		time.Sleep(holdWait)
	}
	if err := t.Run(); err != nil {
		return err
	}
	b.Shelf.Mark(t.Name)
	time.Sleep(b.Pause)
	if time.Since(b.saved) < saveEvery {
		return nil
	}
	b.saved = time.Now()
	return b.Shelf.Save()
}
