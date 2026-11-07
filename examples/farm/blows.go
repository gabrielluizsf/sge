package farm

// The harm done: what bare hands do, what an axe adds for each rank of
// it, what a villager's axe does to the player, and the health a villager
// has. Reach is how far in front a blow lands; npcAnger how long a
// villager goes after the one who struck them; npcLeash how far the player
// must get for them to give up.
const (
	fistHarm  = 5.0
	axeHarm   = 14.0
	axeRank   = 6.0
	npcHarm   = 7.0
	npcHealth = 40.0
	reachBlow = 1.3 * tile
	npcAnger  = 20.0
	npcLeash  = 9 * tile
	npcEvery  = 1.4
	// wolfBlow is how far from the player a blow lands on a wolf: further
	// than its bite, which throws them back from it.
	wolfBlow = 1.9 * tile
)

// blow is how much harm k does with what they hold: more with an axe in
// hand, and more the better it is.
func (k *kit) blow() float64 {
	if k.slot == slotAxe && k.axe != noAxe {
		return axeHarm + axeRank*float64(k.axe)
	}
	return fistHarm
}
