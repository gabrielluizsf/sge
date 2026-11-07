package farm

// takePeople are the villagers, in order.
func takePeople(f *farm) []villagerState {
	out := make([]villagerState, 0, len(f.people))
	for _, v := range f.people {
		b := v.bond.State()
		out = append(out, villagerState{
			Name: v.name, Woman: v.woman, Role: v.role, Genes: v.who.Genes, Kit: takeKit(v.kit),
			Town: takeRegion(v.town), Home: homeIndex(v), Order: takeOrder(v.order),
			Wait: v.wait, At: point{X: v.pos.X, Y: v.pos.Y}, Inside: v.inside, HomeFor: v.homeFor, Fuel: v.fuel,
			ID: v.id, Life: lifeState{Born: v.life.Born, Span: v.life.Span},
			Bond:    bondState{Affection: b.Affection, Courted: int(b.Courted), Gifted: b.Gifted, Talked: b.Talked, Wedding: b.Wedding},
			Partner: indexOfPerson(f, v.partner), Expecting: takeExpecting(v.expecting), Ours: v.ours, Luck: uint32(*v.luck),
		})
	}
	return out
}

// homeIndex is which house of their town is v's, or -1.
func homeIndex(v *villager) int {
	if v.town == nil || v.home == nil {
		return -1
	}
	for i := range v.town.Houses {
		if &v.town.Houses[i] == v.home {
			return i
		}
	}
	return -1
}
