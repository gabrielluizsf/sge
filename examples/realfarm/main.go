// Command realfarm is the farm drawn from real things: its ground, trees,
// crops, houses and people are baked from CC0 models into sprites when
// the game first starts (a few minutes, with a loading screen), kept in
// the user's cache and reused after. What is not needed to begin is baked
// while the game plays.
//
//	go run ./examples/realfarm
package main

import (
	"fmt"
	"os"

	"github.com/gabrielluizsf/sge/examples/farm"
)

func main() {
	b, err := newBaker()
	if err != nil {
		fail(err)
	}
	core, rest := plan(b)
	if os.Getenv("FARM_DEMO") == "" { // a demo has no window to load in: it plays what is baked
		if err := boot(b, core); err != nil {
			fail(err)
		}
		go background(b, rest)
	}
	farm.Main(newArt(b.Shelf), "sge realfarm")
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "realfarm:", err)
	os.Exit(1)
}
