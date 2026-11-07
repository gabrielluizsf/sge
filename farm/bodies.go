package farm

import "github.com/gabrielluizsf/sge/game"

// bodies are the player and the villagers awake, as they stand.
func (f *farm) bodies() []*game.Object {
	var all []*game.Object
	if f.player != nil {
		all = append(all, f.player)
	}
	for _, v := range f.awake() {
		all = append(all, v.obj)
	}
	return all
}

// awake are the villagers whose towns are awake, in the game.
func (f *farm) awake() []*villager {
	var all []*villager
	for _, v := range f.people {
		if v.obj != nil && f.keeper.Awake(v.town) {
			all = append(all, v)
		}
	}
	return all
}
