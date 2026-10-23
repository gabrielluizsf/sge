package game

// Objects are the objects in the world, in the order they were added.
func (g *Game) Objects() []*Object {
	out := make([]*Object, 0, len(g.objects))
	for _, o := range g.objects {
		if !o.destroyed {
			out = append(out, o)
		}
	}
	return out
}
