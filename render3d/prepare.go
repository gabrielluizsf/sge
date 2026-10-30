package render3d

import (
	"cmp"
	"slices"
)

// seen is an instance the camera sees, at the level of detail chosen, as
// deep as its middle is.
type seen struct {
	in    *Instance
	model *Model
	depth float64
}

// Prepare does a frame's work on the CPU: it culls the instances the
// camera does not see, chooses each other's detail, poses them, lays out
// their matrices and palettes for the device, and sorts their draws:
// opaque first, by pipeline and material, front to back, then those cut
// by alpha. Submit then draws them.
func (r *Renderer) Prepare(s *Scene) {
	r.Stats = Stats{}
	r.staging = r.staging[:0]
	r.reserve(FrameFloats)
	r.draws = r.draws[:0]
	r.cull(s)
	slices.SortFunc(r.visible, func(a, b seen) int {
		if c := cmp.Compare(a.model.id, b.model.id); c != 0 {
			return c
		}
		return cmp.Compare(a.depth, b.depth)
	})
	for i := 0; i < len(r.visible); {
		j := i + 1
		for j < len(r.visible) && r.visible[j].model == r.visible[i].model {
			j++
		}
		r.group(r.visible[i:j])
		i = j
	}
	r.order()
	r.light = sunCamera(&s.Sun)
	r.frame(s)
}

// cull keeps in visible the instances s's camera sees, each at the level
// of detail its size on the screen calls for.
func (r *Renderer) cull(s *Scene) {
	r.visible = r.visible[:0]
	f := frustumOf(&s.Camera)
	w, h := r.dev.Size()
	for _, in := range s.Instances {
		if in == nil || in.Model == nil || len(in.Model.Levels) == 0 {
			continue
		}
		lv := in.Model.Levels
		in.level = min(max(in.level, 0), len(lv)-1)
		b := lv[in.level].Model.Bounds.placed(in.Place)
		if !f.sees(b) {
			r.Stats.Culled++
			continue
		}
		px, depth := onScreen(&s.Camera, b, [2]int{w, h})
		in.level = in.Model.choose(in.level, px)
		r.visible = append(r.visible, seen{in: in, model: lv[in.level].Model, depth: depth})
	}
	r.Stats.Shown = len(r.visible)
}

// frame writes the frame's uniforms, at the start of the staging ones.
func (r *Renderer) frame(s *Scene) {
	r.put(FrameViewProj, &s.Camera.ViewProj)
	r.put(FrameLightViewProj, &r.light.ViewProj)
	sun := s.Sun
	t := unit(sun.Toward)
	vals := [...][4]float64{
		{t[0], t[1], t[2], sun.Bias},
		{sun.Color[0], sun.Color[1], sun.Color[2], 1},
		{sun.Ambient[0], sun.Ambient[1], sun.Ambient[2], 1},
		{s.Camera.Eye[0], s.Camera.Eye[1], s.Camera.Eye[2], 1},
	}
	for i, v := range vals {
		for k := range v {
			r.staging[FrameSun+4*i+k] = float32(v[k])
		}
	}
}
