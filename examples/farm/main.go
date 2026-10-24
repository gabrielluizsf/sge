// Command farm is a small farming game on sge's tile maps: buy seeds at
// the stall, sow them, water them, pick what grows and sell it for more
// seeds. Every picture in it, the HUD's too, is drawn from SVG.
//
// FARM_DEMO=n plays n frames by itself, without a window, at a fixed step,
// and prints how the farm did, so a run in Go and one on Deno can be
// compared; FARM_SHOT=file.png also saves the last frame, of a window
// FARM_VIEW=widthxheight in size.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gabrielluizsf/antui/canvas"
)

func main() {
	frames, _ := strconv.Atoi(os.Getenv("FARM_DEMO"))
	var ctl controller = keyboard{}
	if frames > 0 {
		ctl = &bot{tends: 10}
	}
	g, f := newGame(ctl)
	w, h := f.ground.Cols*tile, f.ground.Rows*tile
	if frames <= 0 {
		if err := g.Run("sge farm", w, h); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	fmt.Sscanf(os.Getenv("FARM_VIEW"), "%dx%d", &w, &h)
	cv, err := canvas.NewCanvas(w, h)
	if err != nil {
		panic(err)
	}
	g.Draw(cv) // the view's size, as the window gives it to Run
	for range frames {
		g.Step(1.0 / 60)
	}
	fmt.Printf("money: %d, sold: %d, seeds: %v, growing: %d\n", f.money, f.sold, f.seeds, growing(f))
	if shot := os.Getenv("FARM_SHOT"); shot != "" {
		g.Draw(cv)
		if err := writePNG(cv, shot); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// growing is how many plots have a crop on them.
func growing(f *farm) int {
	n := 0
	for _, p := range f.plots {
		if p.crop >= 0 {
			n++
		}
	}
	return n
}
