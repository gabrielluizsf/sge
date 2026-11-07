// Package animals holds the realistic farm's baked animals: PNG sprites of
// the cow, sheep, wolf and horse in eight facings, a few frames a motion,
// made by cmd/animalart from models by others (see SOURCE.txt), with
// sprites.json saying where each stands and animals.json how each moves.
package animals

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

//go:embed *.png sprites.json animals.json SOURCE.txt
var files embed.FS

// Motion is how an animal moves in one motion: its frames, the seconds a
// loop takes, the speed that keeps its feet from sliding, in cells a
// second, and whether it is run through once.
type Motion struct {
	Frames int     `json:"frames"`
	Length float64 `json:"length"`
	Speed  float64 `json:"speed"`
	Once   bool    `json:"once"`
}

// Info is an animal: how many cells a unit of its model is, its coats and
// its motions.
type Info struct {
	Unit    float64           `json:"unit"`
	Coats   []string          `json:"coats"`
	Motions map[string]Motion `json:"motions"`
}

// Sprite is where a baked frame stands: its size, and the pixel its
// animal's foot is.
type Sprite struct{ W, H, PivotX, PivotY int }

// Frame is the name of a frame: animal_coat_motion_d<facing>_f<frame>.
func Frame(animal, coat, motion string, facing, frame int) string {
	return fmt.Sprintf("%s_%s_%s_d%d_f%d", animal, coat, motion, facing, frame)
}

// Infos are the animals' motions, by name.
func Infos() (map[string]Info, error) {
	var out map[string]Info
	return out, read("animals.json", &out)
}

// Sprites are every frame's place, by name.
func Sprites() (map[string]Sprite, error) {
	var out map[string]Sprite
	return out, read("sprites.json", &out)
}

// PNG is the bytes of frame name's picture.
func PNG(name string) ([]byte, error) { return fs.ReadFile(files, name+".png") }

func read(name string, v any) error {
	data, err := fs.ReadFile(files, name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
