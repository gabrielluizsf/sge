package farm

// sowing returns the crop whose seeds k sows: the ones in hand if they
// have some, else any crop they have seeds of, which then are in hand.
// It is false when they have none.
func (k *kit) sowing() (int, bool) {
	if k.bag.Count(seedsOf(k.seed)) > 0 {
		return k.seed, true
	}
	for c := range crops {
		if k.bag.Count(seedsOf(c)) > 0 {
			k.seed = c
			return c, true
		}
	}
	return 0, false
}

// nextSeed puts in hand the seeds of the next crop k has some of, or, if
// they have none, the next crop.
func (k *kit) nextSeed() {
	for i := 1; i <= len(crops); i++ {
		c := (k.seed + i) % len(crops)
		if k.bag.Count(seedsOf(c)) > 0 {
			k.seed = c
			return
		}
	}
	k.seed = (k.seed + 1) % len(crops)
}
