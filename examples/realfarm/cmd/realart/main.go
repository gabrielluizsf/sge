// Command realart bakes the realistic farm's sprites: every asset of the
// manifest fetched, checked against its pin, read as a glTF model and
// rendered from the game's corner, turned to a few headings so no two
// trees look alike, a model of variants each variant alone, into assets/
// as PNGs, with sprites.json saying where each stands. -crops bakes the
// crops flora grows; -people the villagers, in every motion; -buildings
// the houses and things package building builds, clad in the materials
// the manifest pins; -water the shores of the lake and the sea, every
// mask of them, in looping frames, over the grounds they lie on. The
// game bakes the same way by itself when it starts; this is for baking
// ahead, or into a folder of one's own.
//
//	go run ./examples/realfarm/cmd/realart -manifest examples/realfarm/assets.json -humans examples/realfarm/humans.json -out assets
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

func main() {
	file := flag.String("manifest", "assets.json", "the manifest of assets")
	humans := flag.String("humans", "humans.json", "the pins of the people's data")
	out := flag.String("out", "assets", "the directory to write sprites to")
	only := flag.String("only", "", "bake only the tasks whose name has this in it")
	kinds := map[string]*bool{
		"crop":     flag.Bool("crops", false, "bake the crops flora grows, and no asset"),
		"person":   flag.Bool("people", false, "bake the villagers, made as humans.json pins, and no asset"),
		"building": flag.Bool("buildings", false, "bake the houses and things package building builds, and no asset"),
		"tree":     flag.Bool("trees", false, "bake the trees package tree grows, every species, size and seed, and no asset"),
		"water":    flag.Bool("water", false, "bake the shores of the lake and the sea, and no asset"),
		"import":   flag.Bool("imports", false, "bake the models of your own in SGE_IMPORTS the farm draws, and no asset"),
	}
	workers := flag.Int("workers", 1, "how many frames are baked at once; each takes a core")
	flag.Parse()
	b, err := baker(*file, *humans, *out)
	if err != nil {
		fail(err)
	}
	b.Workers = *workers
	prefixes := chosen(kinds)
	for _, t := range b.Tasks() {
		if !wanted(t.Name, prefixes) || !strings.Contains(t.Name, *only) {
			continue
		}
		if err := b.Do(t); err != nil {
			fail(err)
		}
		fmt.Println(t.Name, "baked")
	}
	if err := b.Shelf.Save(); err != nil {
		fail(err)
	}
}

// baker is a baker onto the folder out, from the manifest and humans
// files.
func baker(manifest, humans, out string) (*bakes.Baker, error) {
	m, err := kit.ReadManifest(manifest)
	if err != nil {
		return nil, err
	}
	h, err := kit.ReadHumans(humans)
	if err != nil {
		return nil, err
	}
	shelf, err := bakes.Open(out)
	if err != nil {
		return nil, err
	}
	return &bakes.Baker{Shelf: shelf, Manifest: m, Humans: h}, nil
}

// chosen are the kinds of task the flags ask for, or the assets' own,
// the grounds and models, when none.
func chosen(kinds map[string]*bool) []string {
	var out []string
	for k, on := range kinds {
		if *on {
			out = append(out, k+" ")
		}
	}
	if len(out) == 0 {
		return []string{"ground ", "model "}
	}
	return out
}

// wanted reports whether the task name is one of the kinds in prefixes.
func wanted(name string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "realart:", err)
	os.Exit(1)
}
