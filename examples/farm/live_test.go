package main

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/sge/human"
)

func waitFor(t *testing.T, what string, ok func() bool) {
	t.Helper()
	for end := time.Now().Add(5 * time.Second); time.Now().Before(end); time.Sleep(2 * time.Millisecond) {
		if ok() {
			return
		}
	}
	t.Fatalf("timed out: %s", what)
}

// The frame takes what is ready and never waits for a person being made.
func TestTheFrameNeverWaitsForABuild(t *testing.T) {
	f := &fakeMaker{gate: make(chan struct{})}
	p := newLivePreview(f, 0)
	defer p.Close()
	p.Ask(someone(1))
	start := time.Now()
	for range 1000 {
		if _, ok := p.Take(); ok {
			t.Fatal("something ready before the build ended")
		}
	}
	if d := time.Since(start); d > 20*time.Millisecond {
		t.Errorf("1000 takes waited %v", d)
	}
	if !p.Busy() {
		t.Error("not busy while it is made")
	}
	close(f.gate)
	waitFor(t, "the person", func() bool { _, ok := p.Take(); return ok })
	waitFor(t, "not busy", func() bool { return !p.Busy() })
}

// A change of colours only does not make the person again: it plans the
// textures of the one shown.
func TestAColourChangeDoesNotMakeThePersonAgain(t *testing.T) {
	f := &fakeMaker{}
	p := newLivePreview(f, 0)
	defer p.Close()
	c := someone(2)
	p.Ask(c)
	var u *liveUpdate
	waitFor(t, "the first person", func() bool { var ok bool; u, ok = p.Take(); return ok })
	if !u.fresh || f.makes() != 1 {
		t.Fatalf("fresh %v, made %d", u.fresh, f.makes())
	}
	for i, edit := range []func(*human.Look){
		func(l *human.Look) { l.SkinTone = 0.9 },
		func(l *human.Look) { l.HairColour = human.RGB(0.9, 0.1, 0.1) },
		func(l *human.Look) { l.EyeColour = human.RGB(0.1, 0.1, 0.9) },
		func(l *human.Look) { l.ClothesColours = []human.Colour{{}, {}, human.RGB(0.2, 0.9, 0.2)} },
		func(l *human.Look) { l.Eyes = "eyes/materials/grey.mhmat" },
	} {
		edit(&c.Look)
		p.Ask(c)
		waitFor(t, "recoloured", func() bool { return !p.Busy() })
		u, ok := p.Take()
		if !ok || u.fresh || len(u.applies) == 0 {
			t.Fatalf("edit %d: ready %v, update %+v", i, ok, u)
		}
	}
	if f.makes() != 1 {
		t.Errorf("made %d people for colour changes, want 1", f.makes())
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.plans != 5 {
		t.Errorf("%d plans, want 5", f.plans)
	}
}

// Colour plans that were not taken yet are all kept, in order, for the
// frame that does take them.
func TestUntakenColourPlansAreKeptInOrder(t *testing.T) {
	f := &fakeMaker{}
	p := newLivePreview(f, 0)
	defer p.Close()
	c := someone(3)
	p.Ask(c)
	waitFor(t, "the person", func() bool { return !p.Busy() })
	p.Take()
	for _, tone := range []float64{0.2, 0.4, 0.6} {
		c.Look.SkinTone = tone
		p.Ask(c)
		waitFor(t, "planned", func() bool { return !p.Busy() })
	}
	u, ok := p.Take()
	if !ok || len(u.applies) != 3 {
		t.Fatalf("ready %v, %d plans", ok, len(u.applies))
	}
}

// Another body or other clothes make the person again.
func TestAChangeOfClothesMakesThePersonAgain(t *testing.T) {
	f := &fakeMaker{}
	p := newLivePreview(f, 0)
	defer p.Close()
	c := someone(4)
	p.Ask(c)
	waitFor(t, "first", func() bool { return !p.Busy() })
	c.Look.Wear = append(append([]string(nil), c.Look.Wear...), "hats/x/x.mhclo")
	p.Ask(c)
	waitFor(t, "second", func() bool { return !p.Busy() && f.makes() == 2 })
	if u, ok := p.Take(); !ok || !u.fresh {
		t.Errorf("update %+v, %v", u, ok)
	}
}

// Of many asks in a row, the newest is what is shown at the end, and the
// ones it replaced are not all made.
func TestTheNewestAskWins(t *testing.T) {
	f := &fakeMaker{delay: 20 * time.Millisecond}
	p := newLivePreview(f, 30*time.Millisecond)
	defer p.Close()
	var last string
	for i := range 8 {
		c := someone(uint32(100 + i))
		last = shapeKey(c)
		p.Ask(c)
		time.Sleep(3 * time.Millisecond)
	}
	waitFor(t, "settled", func() bool { return !p.Busy() })
	u, ok := p.Take()
	if !ok || !u.fresh {
		t.Fatalf("ready %v", ok)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if got := f.made[len(f.made)-1]; got != last {
		t.Errorf("the last made is not the newest ask")
	}
	if len(f.made) > 2 {
		t.Errorf("%d people made for a burst of asks", len(f.made))
	}
}
