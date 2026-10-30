package farm

import "github.com/gabrielluizsf/sge/examples/farm/art"

// Person is the player as the assets draw them, or a villager dressed by
// package art as who they are says.
func (a *svgArt) Person(who Who) Body {
	if who.Player {
		return svgBody{frame: func(n string) string { return spriteSVG("farmer_" + n) }}
	}
	frames := art.PersonFrames(svgLook(who))
	return svgBody{frame: func(n string) string { return frames[n] }}
}

// svgLook is how package art dresses villager who: by their five numbers,
// and the smith in overalls.
func svgLook(who Who) art.Look {
	l := art.Villager(who.Woman, who.Genes)
	if who.Role == roleSmith {
		l.Overalls, l.Pants = true, art.Trousers[2]
	}
	return l
}
