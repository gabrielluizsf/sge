// Command farm is a small farming game on sge's tile maps: buy seeds at
// the stall, sow them, water them, pick what grows and sell it for more
// seeds. Every picture in it, the HUD's too, is drawn from SVG.
//
// FARM_DEMO=n plays n frames by itself, without a window, at a fixed step,
// and prints how the farm did, so a run in Go and one on Deno can be
// compared; FARM_SHOT=file.png also saves the last frame, of a window
// FARM_VIEW=widthxheight in size.
package main

import "github.com/gabrielluizsf/sge/examples/farm"

func main() {
	farm.Main(farm.SVGArt(), "sge farm")
}
