package main

import "time"

// run waits for asks and settles them.
func (p *livePreview) run() {
	for {
		select {
		case <-p.quit:
			return
		case <-p.wake:
		}
		for p.settle() {
		}
	}
}

// settle shows the latest ask: planned if only colours changed, made
// again after the asks have been quiet for wait if not. It reports whether
// to look again, a newer ask having come in.
func (p *livePreview) settle() bool {
	p.mu.Lock()
	g, c, at, done := p.gen, p.want, p.asked, p.done
	p.mu.Unlock()
	if g == done {
		return false
	}
	key := shapeKey(c)
	if p.cur != nil && key == p.curKey {
		apply := p.cur.Plan(c)
		p.publish(g, &liveUpdate{person: p.cur, applies: []func() []int{apply}}, nil)
		return !p.closed()
	}
	if d := p.wait - time.Since(at); d > 0 {
		select {
		case <-time.After(d):
		case <-p.quit:
			return false
		case <-p.wake:
		}
		return true
	}
	pr, err := p.mk.Make(c)
	if err == nil {
		p.cur, p.curKey = pr, key
		p.publish(g, &liveUpdate{person: pr, fresh: true}, nil)
	} else {
		p.publish(g, nil, err)
	}
	return !p.closed()
}

// publish makes u what Take gives next, after what is waiting to be
// taken if u only recolours the same person, instead of it if it is new;
// g is the ask it is of.
func (p *livePreview) publish(g int, u *liveUpdate, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.done, p.err = g, err
	switch {
	case u == nil:
	case p.ready != nil && !u.fresh && p.ready.person == u.person:
		p.ready.applies = append(p.ready.applies, u.applies...)
	default:
		p.ready = u
	}
}

func (p *livePreview) closed() bool {
	select {
	case <-p.quit:
		return true
	default:
		return false
	}
}
