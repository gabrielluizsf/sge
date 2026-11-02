package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/render3d"
)

// liveStage draws the person on a device, live: uploaded once, its
// textures made again when colours change, drawn again when the view
// does. The studio's light is soft: a strong key from the front and left
// and a wide fill, no shadow, over a dark backdrop.
type liveStage struct {
	open    deviceMaker
	dev     render3d.Device
	closeFn func()
	ren     *render3d.Renderer

	person person
	model  *render3d.Model
	inst   *render3d.Instance
	scene  render3d.Scene

	img   *canvas.Canvas
	dirty bool
}

// newLiveStage is a stage drawing on the devices open makes.
func newLiveStage(open deviceMaker) *liveStage {
	s := &liveStage{open: open, dirty: true}
	s.scene.Sun = render3d.Sun{Toward: [3]float64{-0.45, -0.75, 0.55}, Color: [3]float64{0.8, 0.77, 0.72}, Ambient: [3]float64{0.42, 0.43, 0.47}}
	s.scene.Clear = [4]float32{0.045, 0.05, 0.06, 1}
	return s
}

// Height is the shown person's, 0 before there is one.
func (s *liveStage) Height() float64 {
	if s.person == nil {
		return 0
	}
	return heightOf(s.person.Model())
}

// Ready reports whether a person is shown.
func (s *liveStage) Ready() bool { return s.model != nil }

// Close lets go of the device.
func (s *liveStage) Close() {
	if s.closeFn != nil {
		s.closeFn()
		s.closeFn, s.dev, s.ren, s.model = nil, nil, nil, nil
	}
}

// Dirty says the picture is out of date: the view moved.
func (s *liveStage) Dirty() {
	if s != nil {
		s.dirty = true
	}
}

// Show puts u on the stage: a new person is uploaded, with the colours it
// has once u's plans are applied; the one shown has the materials the
// plans changed made again.
func (s *liveStage) Show(u *liveUpdate, w, h int) error {
	if err := s.ensure(w, h); err != nil {
		return err
	}
	var changed []int
	for _, a := range u.applies {
		changed = append(changed, a()...)
	}
	s.dirty = true
	if u.fresh {
		return s.upload(u.person)
	}
	for _, i := range changed {
		if err := s.ren.Retexture(s.model, i); err != nil {
			return err
		}
	}
	return nil
}

// upload shows p, in place of the one before.
func (s *liveStage) upload(p person) error {
	m, err := s.ren.Upload(p.Model())
	if err != nil {
		return err
	}
	if s.model != nil {
		s.model.Release(s.dev)
	}
	s.person, s.model = p, m
	s.inst = &render3d.Instance{Model: render3d.One(m), Place: render3d.At([3]float64{}, 0, 1), Animation: human.Stand.Name}
	s.scene.Instances = []*render3d.Instance{s.inst}
	return nil
}

// ensure has a device of w by h.
func (s *liveStage) ensure(w, h int) error {
	if s.dev == nil {
		dev, closeFn, err := s.open(w, h)
		if err != nil {
			return err
		}
		s.dev, s.closeFn, s.ren = dev, closeFn, render3d.NewRenderer(dev)
		return nil
	}
	if cw, ch := s.dev.Size(); cw != w || ch != h {
		s.dev.Resize(w, h)
		s.dirty = true
	}
	return nil
}
