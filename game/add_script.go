package game

// AddScript attaches a script to the object, after those it has.
func (o *Object) AddScript(s Script) *Object {
	o.scripts = append(o.scripts, s)
	return o
}
