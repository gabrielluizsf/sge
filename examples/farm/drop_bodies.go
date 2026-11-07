package farm

import "github.com/gabrielluizsf/sge/town"

// dropBodies lets go of the bodies and objects of the people of town t, which
// hold the drawings of every frame of each; wake draws them anew, the same,
// from who they are and where they stood (v.pos), when the town comes back.
// Those who live with the player keep theirs.
func (f *farm) dropBodies(t *town.Town) {
	f.leftAt(t)
	for _, v := range f.people {
		if v.ours {
			continue
		}
		if v.town == t && v.guest == nil || v.guest == t {
			leaves(v.obj)
			v.body, v.obj = nil, nil
		}
		if v.guest == t {
			v.guest = nil
		}
	}
}
