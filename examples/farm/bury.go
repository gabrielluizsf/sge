package farm

import "github.com/gabrielluizsf/sge/bond"

// bury is villager v dying of old age: gone from their town and the
// game, dead in the family tree, their partner or the player widowed;
// their house stays. The player hears of it if they knew them.
func (f *farm) bury(v *villager) {
	f.family.Die(v.id)
	f.drop(v)
	if v.partner != nil {
		v.partner.partner = nil
	}
	if f.self.spouse == v {
		f.self.spouse = nil
	}
	if v.ours || v.bond.Affection >= bond.Acquaint {
		f.notify(v.name + " morreu de velhice.")
	}
}

// drop takes villager v out of the farm's people, and their object, out
// of doors or at home, out of the game.
func (f *farm) drop(v *villager) {
	kept := f.people[:0]
	for _, p := range f.people {
		if p != v {
			kept = append(kept, p)
		}
	}
	f.people = kept
	if v.obj != nil {
		v.obj.Destroy()
		v.obj = nil
	}
	home := f.residents[:0]
	for _, r := range f.residents {
		if r.v == v {
			r.obj.Destroy()
			continue
		}
		home = append(home, r)
	}
	f.residents = home
	if f.talk != nil && f.talk.v == v {
		f.talk = nil
	}
}
