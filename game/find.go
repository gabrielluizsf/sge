package game

// Find is the first object of that name, or nil.
func (g *Game) Find(name string) *Object {
	for _, o := range g.objects {
		if o.Name == name && !o.destroyed {
			return o
		}
	}
	return nil
}
