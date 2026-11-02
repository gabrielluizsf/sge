package farm

// sowing returns the crop whose seeds k sows: the ones in hand if they
// have some, else, for all but one who chose their seeds, any crop they
// have seeds of, which then are in hand. It is false when they have none.
func (k *kit) sowing() (int, bool) {
	if k.bag.Count(seedsOf(k.seed)) > 0 {
		return k.seed, true
	}
	if k.exact {
		return 0, false
	}
	return k.anySeeds()
}

// anySeeds returns a crop k has seeds of, which then are in hand.
func (k *kit) anySeeds() (int, bool) {
	for c := range crops {
		if k.bag.Count(seedsOf(c)) > 0 {
			k.seed = c
			return c, true
		}
	}
	return 0, false
}
