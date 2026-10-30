package farm

import (
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/life"
)

// Person is who as package art draws them, as tall as their time of life
// makes them: the founding farmer grown up as the assets draw them, an
// heir in the farmer's hat and overalls with their own skin and hair, or
// a villager dressed as who they are says; grey-haired once old.
func (a *svgArt) Person(who Who) Body {
	size := sizes[who.Stage]
	if who.Player && who.Genes == ([5]int{}) && !who.Woman && who.Stage < life.Elder {
		return svgBody{frame: func(n string) string { return spriteSVG("farmer_" + n) }, size: size}
	}
	l := svgLook(who)
	if who.Player {
		l = heirLook(who)
	}
	frames := art.PersonFrames(l)
	return svgBody{frame: func(n string) string { return frames[n] }, size: size}
}

// sizes are how tall each time of life is drawn, as a part of a grown
// person's height: a baby small, a child about two thirds.
var sizes = map[life.Stage]float64{life.Baby: 0.45, life.Child: 0.65, life.Teen: 0.86, life.Adult: 1, life.Elder: 1}

// grey is the hair of the old.
var grey = art.Hairs[len(art.Hairs)-1]

// svgLook is how package art dresses villager who: by their five numbers,
// the smith in overalls, and the old grey-haired.
func svgLook(who Who) art.Look {
	l := art.Villager(who.Woman, who.Genes)
	if who.Role == roleSmith {
		l.Overalls, l.Pants = true, art.Trousers[2]
	}
	if who.Stage == life.Elder {
		l.Hair = grey
	}
	return l
}

// heirLook is the player as the farmer's look dresses them: the straw hat
// and overalls, with their own skin and hair, a woman's long; grey once
// old.
func heirLook(who Who) art.Look {
	l := art.FarmerLook
	if who.Genes != ([5]int{}) {
		v := art.Villager(who.Woman, who.Genes)
		l.Skin, l.Hair = v.Skin, v.Hair
	}
	if who.Woman {
		l.Style = art.Long
	}
	if who.Stage == life.Elder {
		l.Hair = grey
	}
	return l
}
