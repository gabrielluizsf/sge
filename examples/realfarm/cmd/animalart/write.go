package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

// facings are the ways each motion is baked, as the people are.
const facings = bakes.Facings

// writer puts the frames in a folder, and says how they move.
type writer struct {
	shelf *bakes.Shelf
	info  map[string]*info
}

// info is how an animal moves, for the game: the coats it comes in, and
// each motion's frames, period and speed.
type info struct {
	Unit    float64               `json:"unit"`
	Coats   []string              `json:"coats"`
	Motions map[string]motionInfo `json:"motions"`
}

type motionInfo struct {
	Frames int     `json:"frames"`
	Length float64 `json:"length"`
	Speed  float64 `json:"speed,omitempty"`
	Once   bool    `json:"once,omitempty"`
}

func newWriter(dir string) (*writer, error) {
	shelf, err := bakes.Open(dir)
	if err != nil {
		return nil, err
	}
	w := &writer{shelf: shelf, info: map[string]*info{}}
	if data, err := os.ReadFile(filepath.Join(dir, "animals.json")); err == nil {
		json.Unmarshal(data, &w.info) // what earlier runs baked, kept
	}
	return w, nil
}

// frameName is the name of a frame: animal_coat_motion_d<facing>_f<frame>.
func frameName(animal, coat, motion string, facing, frame int) string {
	return fmt.Sprintf("%s_%s_%s_d%d_f%d", animal, coat, motion, facing, frame)
}

// put bakes and writes one frame.
func (w *writer) put(name string, sp bake.Sprite) error { return w.shelf.Put(name, sp) }

// finish writes sprites.json and animals.json, and drops the baker's
// bookkeeping the game does not read.
func (w *writer) finish() error {
	if err := w.shelf.Save(); err != nil {
		return err
	}
	os.Remove(filepath.Join(w.shelf.Dir, "done.json"))
	data, err := json.MarshalIndent(w.info, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(w.shelf.Dir, "animals.json"), append(data, '\n'), 0o644)
}
