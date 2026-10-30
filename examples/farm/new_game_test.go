package farm

import "github.com/gabrielluizsf/sge/game"

// newGame is the farm in its stylized look, played by ctl, as the tests
// play it.
func newGame(ctl Controller) (*game.Game, *farm) { return build(SVGArt(), ctl) }
