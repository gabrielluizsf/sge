package game

// Step advances the world dt seconds: objects new to it start, every
// script updates, animations move on, touching objects are told about
// each other, and destroyed objects leave. Objects added during the step
// take part from the next one.
func (g *Game) Step(dt float64) {
	g.time += dt
	objs := g.Objects()
	for _, o := range objs {
		if !o.started {
			o.started = true
			for _, s := range o.scripts {
				if st, ok := s.(Starter); ok {
					st.Start(o)
				}
			}
		}
	}
	for _, o := range objs {
		for _, s := range o.scripts {
			if !o.destroyed {
				s.Update(o, dt)
			}
		}
		if o.Animator != nil {
			o.Animator.Update(dt)
		}
	}
	g.collide(objs)
	kept := g.objects[:0]
	for _, o := range g.objects {
		if o.destroyed {
			o.game = nil
			continue
		}
		kept = append(kept, o)
	}
	clear(g.objects[len(kept):])
	g.objects = kept
}
