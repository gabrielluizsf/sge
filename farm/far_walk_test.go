//go:build linux

package farm

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// walkStats is the cost of the frames of one stretch of a long walk.
type walkStats struct {
	step, draw []time.Duration
}

func (w walkStats) String() string {
	return fmt.Sprintf("step avg %.2f p99 %.2f max %.2f | draw avg %.2f p99 %.2f max %.2f ms | frames over 16.6ms: %d",
		avg(w.step), pct(w.step, 0.99), pct(w.step, 1), avg(w.draw), pct(w.draw, 0.99), pct(w.draw, 1), w.slow())
}

// slow is how many frames took longer than a 60 Hz frame.
func (w walkStats) slow() int {
	n := 0
	for i := range w.step {
		if w.step[i]+w.draw[i] > 16600*time.Microsecond {
			n++
		}
	}
	return n
}

func avg(d []time.Duration) float64 {
	var s time.Duration
	for _, x := range d {
		s += x
	}
	return float64(s) / float64(max(1, len(d))) / 1e6
}

func pct(d []time.Duration, p float64) float64 {
	c := append([]time.Duration(nil), d...)
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })
	if len(c) == 0 {
		return 0
	}
	return float64(c[min(len(c)-1, int(p*float64(len(c))))]) / 1e6
}

// cpuNow is the CPU time this thread has used: a frame measured by it is
// not made longer by other programs sharing the machine.
func cpuNow() time.Duration {
	var ru syscall.Rusage
	if err := syscall.Getrusage(1, &ru); err != nil { // RUSAGE_THREAD
		return 0
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

// run plays frames frames, moving the farmer by pace(i) pixels east
// before each, and times the steps and the draws.
func run(g *game.Game, cv *canvas.Canvas, frames int, pace func(i int) float64) walkStats {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := g.Find("farmer")
	var w walkStats
	for i := range frames {
		o.Transform.Position.X += pace(i)
		t0 := cpuNow()
		g.Step(frame)
		t1 := cpuNow()
		g.Draw(cv)
		w.step = append(w.step, t1-t0)
		w.draw = append(w.draw, cpuNow()-t1)
	}
	return w
}

// liveHeap is the megabytes of the heap still held after a collection.
func liveHeap() uint64 {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return m.HeapAlloc >> 20
}

// TestWalkFar runs the farmer east, stretch after stretch, through town
// after town, and reports the frames of each, after the frames of the
// farmer pacing to and fro near the farm: FARM_FAR=stretches turns it on.
func TestWalkFar(t *testing.T) {
	stretches, _ := strconv.Atoi(os.Getenv("FARM_FAR"))
	if stretches == 0 {
		t.Skip("FARM_FAR not set")
	}
	g, f := newGame(idle{})
	cv, _ := canvas.NewCanvas(1280, 720)
	g.Draw(cv)
	to := func(i int) float64 { return runSpeed / 60 * float64(1-i/100%2*2) } // 100 frames each way
	run(g, cv, 600, to)
	fmt.Printf("near the farm    : %s\n", run(g, cv, 1500, to))
	o := g.Find("farmer")
	for s := range stretches {
		w := run(g, cv, 1500, func(int) float64 { return runSpeed / 60 })
		fmt.Printf("stretch %2d x=%5.0f: %s | live heap %dMB, objects %d, towns %d, warm queue %d\n",
			s, o.Transform.Position.X/tile, w, liveHeap(), len(g.Objects()), len(f.settled), f.warm.Waiting())
	}
	if path := os.Getenv("FARM_HEAP"); path != "" { // what the walk still holds
		if out, err := os.Create(path); err == nil {
			pprof.Lookup("heap").WriteTo(out, 0)
			out.Close()
		}
	}
	runtime.KeepAlive(g)
}
