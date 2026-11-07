package main

import (
	_ "embed"
	"encoding/json"
	"time"

	"github.com/gabrielluizsf/sge/examples/farm/kit/bakes"
)

//go:embed assets.json
var assetsJSON []byte

//go:embed humans.json
var humansJSON []byte

// newBaker is a baker onto the cache, from the pins this program holds.
func newBaker() (*bakes.Baker, error) {
	b := &bakes.Baker{Workers: 1, Pause: 5 * time.Millisecond}
	if err := json.Unmarshal(assetsJSON, &b.Manifest); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(humansJSON, &b.Humans); err != nil {
		return nil, err
	}
	dir, err := cacheDir()
	if err != nil {
		return nil, err
	}
	shelf, err := bakes.Open(dir)
	b.Shelf = shelf
	return b, err
}

// plan splits what the farm needs and has not yet baked into what it
// cannot begin without, core, and the rest.
func plan(b *bakes.Baker) (core, rest []bakes.Task) {
	for _, t := range b.Tasks() {
		switch {
		case !needed(t) || b.Shelf.Finished(t.Name):
		case t.Core:
			core = append(core, t)
		default:
			rest = append(rest, t)
		}
	}
	return core, rest
}
