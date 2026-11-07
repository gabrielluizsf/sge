package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
)

// background bakes tasks one frame at a time on a single core, resting
// between them, while the game plays, until stop is set: the villagers' other stages of life
// and motions. A sprite it has not baked yet is drawn as the nearest one
// that is.
func background(b *bakes.Baker, tasks []bakes.Task, stop *atomic.Bool) {
	b.Workers = 1
	b.Pause = 50 * time.Millisecond
	for _, t := range tasks {
		if stop.Load() {
			return
		}
		if err := b.Do(t); err != nil {
			fmt.Fprintln(os.Stderr, "farm:", t.Name+":", err)
		}
	}
	if err := b.Shelf.Save(); err != nil {
		fmt.Fprintln(os.Stderr, "farm:", err)
	}
}
