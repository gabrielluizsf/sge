package game

// Add puts o in the world and returns it. Its scripts start on the next
// step.
func (g *Game) Add(o *Object) *Object {
	o.game = g
	o.destroyed = false
	g.objects = append(g.objects, o)
	return o
}
