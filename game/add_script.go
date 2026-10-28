package game

// AddScript attaches a script to the object, after those it has.
func (o *Object) AddScript(s Script) *Object {
	o.scripts = append(o.scripts, s)
	if _, ok := s.(CollisionHandler); ok {
		o.handles = true
	}
	return o
}
