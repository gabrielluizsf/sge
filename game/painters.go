package game

// versions appends to v the version of each script of o that draws, as a
// Painter or a Spreader, and reports whether all of them have one.
func (o *Object) versions(v []uint64) ([]uint64, bool) {
	for _, s := range o.scripts {
		_, paints := s.(Painter)
		_, spreads := s.(Spreader)
		if !paints && !spreads {
			continue
		}
		vs, ok := s.(Versioned)
		if !ok {
			return v, false
		}
		v = append(v, vs.Version())
	}
	return v, true
}

// draws reports whether any script of o is a Painter or a Spreader.
func (o *Object) draws() bool {
	for _, s := range o.scripts {
		_, paints := s.(Painter)
		_, spreads := s.(Spreader)
		if paints || spreads {
			return true
		}
	}
	return false
}
