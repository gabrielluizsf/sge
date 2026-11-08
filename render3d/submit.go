package render3d

// Submit draws what Prepare laid out: it hands the device the frame's
// uniforms, draws the sun's shadow map when the sun casts one, and then
// the scene on the screen.
func (r *Renderer) Submit(s *Scene) error {
	if err := r.upload(); err != nil {
		return err
	}
	shadow, err := r.shadowMap(s.Sun.Shadow)
	if err != nil {
		return err
	}
	if shadow != 0 {
		p := r.dev.Begin(PassDesc{Depth: shadow, DepthOnly: true, ClearDepth: true})
		err = r.issue(p, passOf{shadow: true})
		p.End()
		if err != nil {
			return err
		}
	}
	p := r.dev.Begin(PassDesc{Target: Screen, Clear: true, ClearColor: s.Clear, ClearDepth: true})
	err = r.issue(p, passOf{shadowMap: shadow})
	p.End()
	return err
}

// upload writes the staging uniforms to the device's buffer, made larger
// first when they no longer fit.
func (r *Renderer) upload() error {
	if len(r.staging) > r.capacity {
		if r.uniforms != 0 {
			r.dev.Release(r.uniforms)
			r.uniforms, r.capacity = 0, 0
		}
		n := aligned(max(len(r.staging)*3/2, 4096))
		b, err := r.dev.NewBuffer(BufferDesc{Kind: Uniforms, Size: n})
		if err != nil {
			return err
		}
		r.uniforms, r.capacity = b, n
	}
	r.dev.Write(r.uniforms, 0, r.staging)
	return nil
}

// passOf is which pass draws are issued for: the shadow's, or the
// scene's, with the shadow map it reads.
type passOf struct {
	shadow    bool
	shadowMap Texture
}

// issue issues the frame's draws in p, setting only what changes.
func (r *Renderer) issue(p Pass, of passOf) error {
	var pipe Pipeline
	var verts, idx Buffer
	frame := Range{Buffer: r.uniforms, Count: FrameFloats}
	for _, d := range r.draws {
		m := d.model
		pt := m.parts[d.part]
		mt := m.materials[pt.material]
		pl, err := r.pipeline(kindOf(pt, of.shadow), mt.double)
		if err != nil {
			return err
		}
		if pl != pipe {
			p.SetPipeline(pl)
			pipe = pl
		}
		v := m.static
		if pt.skin >= 0 {
			v = m.skinned
		}
		if v != verts {
			p.SetVertices(v)
			verts = v
		}
		if m.indices != idx {
			p.SetIndices(m.indices)
			idx = m.indices
		}
		p.SetBindings(Bindings{
			Frame:    frame,
			Material: Range{Buffer: m.uniforms, At: mt.at, Count: MaterialFloats},
			Objects:  Range{Buffer: r.uniforms, At: d.objects, Count: d.count * ObjectFloats},
			Texture:  mt.texture,
			Shadow:   of.shadowMap,
		})
		p.Draw(Draw{First: pt.first, Count: pt.count, Instances: d.instances})
	}
	return nil
}
