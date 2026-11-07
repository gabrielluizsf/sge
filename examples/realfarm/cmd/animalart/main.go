// Command animalart bakes the realistic farm's animals: the cow, the
// sheep, the wolf and the horse, each a model by someone else (see the
// SOURCE.txt it writes), read through package gltf, given the bones and the
// walks it lacks, and rendered from the game's corner in eight facings, a
// few frames each, into examples/realfarm/animals as PNGs, with
// sprites.json saying where each stands and animals.json how each moves.
//
//	go run ./examples/realfarm/cmd/animalart -models <folder of the extracted models>
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	models := flag.String("models", "", "the folder holding each model's own folder, as extracted")
	out := flag.String("out", "examples/realfarm/animals", "the folder to write the sprites to")
	only := flag.String("only", "", "bake only the animals whose name has this in it")
	flag.Parse()
	if *models == "" {
		fail(fmt.Errorf("-models is needed"))
	}
	root, err := os.OpenRoot(*models)
	if err != nil {
		fail(err)
	}
	defer root.Close()
	if *only == "stride" {
		v, err := measureHorse(root)
		fmt.Println("horse body speed, model units a second:", v, err)
		return
	}
	w, err := newWriter(*out)
	if err != nil {
		fail(err)
	}
	for _, a := range animals {
		if !strings.Contains(a.name, *only) {
			continue
		}
		if err := bakeAnimal(root, a, w); err != nil {
			fail(fmt.Errorf("%s: %w", a.name, err))
		}
		fmt.Println(a.name, "baked")
	}
	if err := w.finish(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "animalart:", err)
	os.Exit(1)
}
