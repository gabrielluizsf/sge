package game

// collide tells every pair of touching objects about each other, through
// the scripts of each that handle collisions.
func (g *Game) collide(objs []*Object) {
	for i, a := range objs {
		for _, b := range objs[i+1:] {
			if a.destroyed || b.destroyed || !a.Collides(b) {
				continue
			}
			tell(a, b)
			tell(b, a)
		}
	}
}

func tell(o, other *Object) {
	for _, s := range o.scripts {
		if h, ok := s.(CollisionHandler); ok {
			h.Collision(o, other)
		}
	}
}
