//go:build linux && amd64

package gl

import "github.com/gabrielluizsf/sge/render3d"

// op is what a recorded command does.
type op uint8

const (
	opWrite op = iota
	opBegin
	opPipeline
	opVertices
	opIndices
	opBindings
	opDraw
)

// cmd is a command recorded for the device's thread: a Write of count
// floats from from on into buf at at, or a pass's.
type cmd struct {
	op              op
	buf             render3d.Buffer
	pipe            render3d.Pipeline
	at, from, count int
	pass            render3d.PassDesc
	bind            render3d.Bindings
	draw            render3d.Draw
}

// pass records a pass's commands on its device, for End to hand over.
type pass struct{ d *Device }

// Begin starts a pass into p's targets, cleared as p says; a pass whose
// targets do not exist draws nothing.
func (d *Device) Begin(p render3d.PassDesc) render3d.Pass {
	d.queue = append(d.queue, cmd{op: opBegin, pass: p})
	return pass{d}
}

func (p pass) SetPipeline(pl render3d.Pipeline) {
	p.d.queue = append(p.d.queue, cmd{op: opPipeline, pipe: pl})
}

func (p pass) SetVertices(b render3d.Buffer) {
	p.d.queue = append(p.d.queue, cmd{op: opVertices, buf: b})
}

func (p pass) SetIndices(b render3d.Buffer) {
	p.d.queue = append(p.d.queue, cmd{op: opIndices, buf: b})
}

func (p pass) SetBindings(b render3d.Bindings) {
	p.d.queue = append(p.d.queue, cmd{op: opBindings, bind: b})
}

func (p pass) Draw(dr render3d.Draw) {
	p.d.queue = append(p.d.queue, cmd{op: opDraw, draw: dr})
}

// End hands what the pass recorded to the device's thread, and waits for
// it to be issued.
func (p pass) End() { p.d.th.do(p.d.run) }

// runQueue issues what is recorded, in order, on the thread.
func (d *Device) runQueue() {
	for i := range d.queue {
		c := &d.queue[i]
		switch c.op {
		case opWrite:
			d.write(c)
		case opBegin:
			d.begin(&c.pass)
		case opPipeline:
			d.setPipeline(c.pipe)
		case opVertices:
			d.st.verts = c.buf
		case opIndices:
			d.st.idx = c.buf
			call(d.gl.bindBuffer, n(glElementArrayBuffer), n(d.bufferName(c.buf)))
		case opBindings:
			d.st.bind, d.st.bound = c.bind, false
		case opDraw:
			d.drawElements(c.draw)
		}
	}
	d.queue, d.floats = d.queue[:0], d.floats[:0]
}
