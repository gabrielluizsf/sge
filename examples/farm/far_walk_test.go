//go:build linux

package main

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
	"github.com/gabrielluizsf/sge/farm"
)

// cpuNow is the CPU time this thread has used, so a frame measured by it
// is not lengthened by other programs on the machine.
func cpuNow() time.Duration {
	var ru syscall.Rusage
	if syscall.Getrusage(1, &ru) != nil { // RUSAGE_THREAD
		return 0
	}
	return time.Duration(ru.Utime.Nano() + ru.Stime.Nano())
}

func pct(d []time.Duration, p float64) float64 {
	c := append([]time.Duration(nil), d...)
	sort.Slice(c, func(i, j int) bool { return c[i] < c[j] })
	return float64(c[min(len(c)-1, int(p*float64(len(c))))]) / 1e6
}

func mean(d []time.Duration) float64 {
	var s time.Duration
	for _, x := range d {
		s += x
	}
	return float64(s) / float64(len(d)) / 1e6
}

// TestRealfarmWalkFar runs the farmer east through the real art, stretch
// after stretch, and reports the frames: REALFARM_FAR=stretches turns it
// on; the shelf in the cache must hold what is baked.
func TestRealfarmWalkFar(t *testing.T) {
	stretches, _ := strconv.Atoi(os.Getenv("REALFARM_FAR"))
	if stretches == 0 {
		t.Skip("REALFARM_FAR not set")
	}
	b, err := newBaker()
	if err != nil {
		t.Skip(err)
	}
	s := farm.New(newArt(b.Shelf), farm.Bot(0))
	cv, _ := canvas.NewCanvas(1280, 720)
	s.Game.Draw(cv)
	o := s.Game.Find("farmer")
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	pace := 216.0 / 60
	for range 300 {
		s.Game.Step(1.0 / 60)
		s.Game.Draw(cv)
	}
	for st := range stretches {
		var step, draw []time.Duration
		slow := 0
		for range 1500 {
			o.Transform.Position.X += pace
			t0 := cpuNow()
			s.Game.Step(1.0 / 60)
			t1 := cpuNow()
			s.Game.Draw(cv)
			t2 := cpuNow()
			step, draw = append(step, t1-t0), append(draw, t2-t1)
			if t2-t0 > 16600*time.Microsecond {
				slow++
			}
		}
		var m runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m)
		fmt.Printf("stretch %2d x=%5.0f: step avg %.2f p99 %.2f max %.2f | draw avg %.2f p99 %.2f max %.2f ms | over 16.6ms: %d | live heap %dMB\n",
			st, o.Transform.Position.X/32, mean(step), pct(step, .99), pct(step, 1), mean(draw), pct(draw, .99), pct(draw, 1), slow, m.HeapAlloc>>20)
	}
	if path := os.Getenv("REALFARM_HEAP"); path != "" { // what the walk still holds
		if out, err := os.Create(path); err == nil {
			pprof.Lookup("heap").WriteTo(out, 0)
			out.Close()
		}
	}
	runtime.KeepAlive(s)
}
