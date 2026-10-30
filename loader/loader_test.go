package loader

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// job is a job making its name, slept so long, and counting its runs.
func job(name, group string, bytes int64, sleep time.Duration, runs *[]string) Job {
	return Job{Name: name, Group: group, Bytes: bytes, Make: func() (any, error) {
		time.Sleep(sleep)
		*runs = append(*runs, name)
		return name + "!", nil
	}}
}

// Jobs run most urgent first, in the order asked among equals; asking
// again only hurries one.
func TestJobsRunByPriority(t *testing.T) {
	var runs []string
	var l Loader
	l.Request(job("far", "", 1, 0, &runs), 5)
	l.Request(job("near", "", 1, 0, &runs), 1)
	l.Request(job("also near", "", 1, 0, &runs), 1)
	l.Request(job("far", "", 1, 0, &runs), 0)
	l.Drain()
	if got := fmt.Sprint(runs); got != "[far near also near]" {
		t.Errorf("ran %s", got)
	}
	if v, err, ok := l.Get("near"); !ok || err != nil || v != "near!" {
		t.Errorf("near is %v %v %v", v, err, ok)
	}
	l.Request(job("near", "", 1, 0, &runs), 0)
	if l.Pending() != 0 {
		t.Error("an asset made was asked for again")
	}
}

// A frame spends no more than its budget on jobs, but at least one.
func TestAFrameKeepsToItsBudget(t *testing.T) {
	var runs []string
	l := &Loader{Budget: 5 * time.Millisecond}
	for i := range 10 {
		l.Request(job(fmt.Sprint(i), "boot", 1, 3*time.Millisecond, &runs), 0)
	}
	l.Update(nil, 0)
	if n := len(runs); n < 1 || n > 3 {
		t.Errorf("a frame ran %d jobs of 3ms in 5ms", n)
	}
	l.Budget = time.Microsecond
	before := len(runs)
	l.Update(nil, 0)
	if len(runs) != before+1 {
		t.Error("a frame over budget ran no job")
	}
	if done, asked := l.Progress("boot"); done != len(runs) || asked != 10 {
		t.Errorf("progress %d of %d", done, asked)
	}
}

// Past the memory budget the least lately used asset is let go, and made
// again when asked for.
func TestTheLeastLatelyUsedIsLetGo(t *testing.T) {
	var runs []string
	l := &Loader{Memory: 25}
	l.Request(job("a", "", 10, 0, &runs), 0)
	l.Request(job("b", "", 10, 0, &runs), 0)
	l.Drain()
	l.Get("a") // b is now the least lately used
	l.Request(job("c", "", 10, 0, &runs), 0)
	l.Drain()
	if _, _, ok := l.Get("b"); ok {
		t.Error("b was kept past the budget")
	}
	if _, _, ok := l.Get("a"); !ok {
		t.Error("a, used lately, was let go")
	}
	l.Request(job("b", "", 10, 0, &runs), 0)
	l.Drain()
	if _, _, ok := l.Get("b"); !ok || len(runs) != 4 {
		t.Errorf("b was not made again: %v", runs)
	}
}

// A job's error is kept as its result.
func TestAFailedJobKeepsItsError(t *testing.T) {
	var l Loader
	boom := errors.New("boom")
	l.Request(Job{Name: "x", Make: func() (any, error) { return nil, boom }}, 0)
	l.Drain()
	if _, err, ok := l.Get("x"); !ok || err != boom {
		t.Errorf("error %v", err)
	}
}

// The loading screen covers the game while its group loads, its bar as
// full as the group is done, and shows nothing once it is.
func TestTheLoadingScreenShowsProgress(t *testing.T) {
	var runs []string
	l := &Loader{}
	for i := range 4 {
		l.Request(job(fmt.Sprint(i), "boot", 1, 0, &runs), 0)
	}
	s := &Screen{Loader: l, Group: "boot", Back: canvas.RGB(0, 0, 0), Track: canvas.RGB(50, 50, 50), Fill: canvas.RGB(255, 255, 255)}
	o := game.NewObject("loading")
	o.AddScript(s)
	cv, _ := canvas.NewCanvas(100, 60)
	l.runNext()
	l.runNext()
	cv.Clear(canvas.RGB(9, 9, 9))
	s.Paint(o, cv)
	y := 60 * 2 / 3
	if cv.At(26, y) != canvas.RGB(255, 255, 255) || cv.At(74, y) != canvas.RGB(50, 50, 50) || cv.At(1, 1) != canvas.RGB(0, 0, 0) {
		t.Errorf("half done: %v %v %v", cv.At(26, y), cv.At(74, y), cv.At(1, 1))
	}
	l.Drain()
	cv.Clear(canvas.RGB(9, 9, 9))
	s.Paint(o, cv)
	if s.Loading() || cv.At(1, 1) != canvas.RGB(9, 9, 9) {
		t.Error("the screen covers the game after loading")
	}
}
