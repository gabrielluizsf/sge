package main

import (
	"sync"
	"time"

	"github.com/gabrielluizsf/sge/examples/farm/kit"
)

// livePreview makes the person of the character asked for, off the
// frame, in one goroutine: Ask returns at once and the frame takes what
// is ready with Take, which never waits. A change of colours alone only
// plans new textures for the person shown; a change of body or clothes
// makes the person again, once the asks have been quiet for wait. A
// character asked while another is being made takes its place: the
// latest wins.
type livePreview struct {
	mk   maker
	wait time.Duration

	mu    sync.Mutex
	want  kit.Character
	gen   int // which ask is the latest
	asked time.Time
	done  int // the ask the person shown, or about to be, is of
	ready *liveUpdate
	err   error

	cur    person // the worker's: the person last made
	curKey string

	wake, quit chan struct{}
	closeOnce  sync.Once
}

// newLivePreview starts the goroutine that makes people with mk.
func newLivePreview(mk maker, wait time.Duration) *livePreview {
	p := &livePreview{mk: mk, wait: wait, wake: make(chan struct{}, 1), quit: make(chan struct{})}
	go p.run()
	return p
}

// Ask wants c shown.
func (p *livePreview) Ask(c kit.Character) {
	p.mu.Lock()
	p.want, p.asked = c, time.Now()
	p.gen++
	p.mu.Unlock()
	select {
	case p.wake <- struct{}{}:
	default:
	}
}

// Take is what is ready since the last Take, none if nothing is: it
// never waits for a person being made.
func (p *livePreview) Take() (*liveUpdate, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	u := p.ready
	p.ready = nil
	return u, u != nil
}

// Busy reports whether the last character asked is not ready yet.
func (p *livePreview) Busy() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.gen != p.done
}

// Err is why the last person could not be made, nil if it could.
func (p *livePreview) Err() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.err
}

// Close stops the goroutine once what it is making is made; it does not
// wait for that.
func (p *livePreview) Close() { p.closeOnce.Do(func() { close(p.quit) }) }
