package farm

import (
	"github.com/gabrielluizsf/sge/random"
)

// How long a chore takes, in hours of it, for each thing it makes of
// the world, and the most that is made of it in one spell.
const (
	fellEvery, fellMost   = 1.0, 3
	mineEvery, mineMost   = 1.5, 2
	tendEvery, tendMost   = 0.5, 4
	forgeEvery, forgeMost = 0.5, 3
)

// settleSpell settles hours hours of villager v doing chore c, in the spell
// key names: what they made of the world meanwhile, as they would have if
// watched, but not step by step. The same spell settles the same.
func (f *farm) settleSpell(v *villager, c chore, key uint32, hours float64) {
	src := random.New(key)
	switch c {
	case choreChop:
		f.fellFor(v, src, min(int(hours/fellEvery), fellMost))
	case choreMine:
		f.mineFor(v, src, min(int(hours/mineEvery), mineMost))
	case chorePlant:
		f.tendFor(v, min(int(hours/tendEvery), tendMost))
	case choreTrade:
		v.kit.sellOut()
	case choreForge:
		f.forgeFor(v, min(int(hours/forgeEvery), forgeMost))
	}
}

// fellFor has v fell up to n of the trees about their town, each a stump
// that grows back, and the wood theirs.
func (f *farm) fellFor(v *villager, src *random.Source, n int) {
	sites := f.sitesOf(f.chopQuarry(v), false)
	for ; n > 0 && len(sites) > 0; n-- {
		i := src.Intn(len(sites))
		at := sites[i].target
		sites = append(sites[:i], sites[i+1:]...)
		tree := f.thingAt(at.Col, at.Row)
		v.kit.bag.Add(wood, src.Between(1, axes[v.kit.axe].most))
		f.setThing(at.Col, at.Row, tileID["stump"])
		f.stumps = append(f.stumps, &stump{col: at.Col, row: at.Row, tree: tree, left: regrowIn})
	}
}

// mineFor has v break up to n of the rocks about their town, each leaving
// ground that grows a rock back, and the stone theirs.
func (f *farm) mineFor(v *villager, src *random.Source, n int) {
	sites := f.sitesOf(f.mineQuarry(v), false)
	if v.kit.bag.Count(itemPick) == 0 {
		v.kit.bag.Add(itemPick, 1)
	}
	for ; n > 0 && len(sites) > 0; n-- {
		i := src.Intn(len(sites))
		at := sites[i].target
		sites = append(sites[:i], sites[i+1:]...)
		rock := f.thingAt(at.Col, at.Row)
		v.kit.bag.Add(stone, src.Between(stoneLeast, stoneMost))
		f.setThing(at.Col, at.Row, leftBy(rock))
		f.stumps = append(f.stumps, &stump{col: at.Col, row: at.Row, tree: rock, left: rockRegrowIn})
	}
}

// tendFor has v do up to n of the tasks of their garden: pick what is ripe
// into their bag, clear what withered, sow what is bare. The water is the
// morning's, as it is for a garden that is not watched.
func (f *farm) tendFor(v *villager, n int) {
	if v.home == nil {
		return
	}
	p := f.plotAt(v.home.Garden.Col, v.home.Garden.Row)
	for ; n > 0 && p != nil; n-- {
		switch {
		case p.ripened():
			v.kit.bag.Add(cropOf(p.crop), 1)
			f.pick(p)
		case p.dead:
			p.clear()
		case p.crop < 0 && f.sow(v.kit, p):
		default:
			return
		}
		f.show(p)
	}
}

// forgeFor has v, a smith, forge up to n tools from their fuel, each sold
// for its price, and burn the wood they have into fuel.
func (f *farm) forgeFor(v *villager, n int) {
	k := v.kit
	if w := k.bag.Count(wood); w > 0 && v.smith() {
		k.bag.Take(wood, w)
		v.fuel += float64(w)
	}
	for ; n > 0 && v.fuel >= 1 && v.smith(); n-- {
		v.fuel--
		k.money += toolPrice
	}
}
