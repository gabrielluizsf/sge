package farm

import "github.com/gabrielluizsf/sge/vector"

// spot is where villager v stands as the player sees them: in the room
// the player is in, out of doors, or where they were last.
func (f *farm) spot(v *villager) vector.Vector3 {
	for _, r := range f.residents {
		if r.v == v {
			return r.obj.Transform.Position
		}
	}
	if v.obj != nil {
		return v.obj.Transform.Position
	}
	return v.pos
}
