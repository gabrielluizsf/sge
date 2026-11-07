package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/render3d"
)

// Render is the picture of the person, w by h, seen from v: drawn again
// only if something changed since the last, else the last.
func (s *liveStage) Render(v *liveView, w, h int) (*canvas.Canvas, error) {
	if s.model == nil {
		return nil, nil
	}
	if err := s.ensure(w, h); err != nil {
		return nil, err
	}
	if !s.dirty && s.img != nil {
		return s.img, nil
	}
	v.height = s.Height()
	s.inst.Place = render3d.At([3]float64{}, v.yaw, 1)
	s.scene.Camera = v.Camera(float64(w) / float64(max(h, 1)))
	if err := s.ren.Draw(&s.scene); err != nil {
		return nil, err
	}
	px, err := s.dev.Read(render3d.Screen)
	if err != nil {
		return nil, err
	}
	if s.img == nil || s.img.Width != w || s.img.Height != h {
		if s.img, err = canvas.NewCanvas(w, h); err != nil {
			return nil, err
		}
	}
	for y := range h {
		row := s.img.Pixels[y*s.img.Stride : y*s.img.Stride+w]
		for x := range row {
			p := px.Pix[4*(y*w+x) : 4*(y*w+x)+4]
			row[x] = canvas.RGB(p[0], p[1], p[2])
		}
	}
	s.dirty = false
	return s.img, nil
}
