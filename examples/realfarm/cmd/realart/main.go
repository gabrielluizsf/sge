// Command realart bakes the realistic farm's sprites: every asset of the
// manifest fetched, checked against its pin, read as a glTF model and
// rendered from the game's corner, turned to a few headings so no two
// trees look alike, a model of variants each variant alone, into assets/ as PNGs, with sprites.json saying where
// each stands.
//
//	go run ./cmd/realart -manifest assets.json -out assets
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"os"
	"path"
	"path/filepath"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm"
	"github.com/gabrielluizsf/sge/flora"
	"github.com/gabrielluizsf/sge/gltf"
)

// sprite is where a baked picture stands: the pixel its model's foot is.
type sprite struct {
	W, H   int
	PivotX int
	PivotY int
}

// patches are the cells a ground's texture spans: two by two metres, so
// its four tiles laid by the cells' parity meet without a seam.
var patches = []bake.Patch{{Span: 2, Col: 0, Row: 0}, {Span: 2, Col: 1, Row: 0}, {Span: 2, Col: 0, Row: 1}, {Span: 2, Col: 1, Row: 1}}

// headings are the turns each model is baked at.
var headings = []float64{45, 135, 225, 315}

func main() {
	file := flag.String("manifest", "assets.json", "the manifest of assets")
	out := flag.String("out", "assets", "the directory to write sprites to")
	only := flag.String("only", "", "bake only this asset")
	crops := flag.Bool("crops", false, "bake the crops flora grows, and no asset")
	flag.Parse()
	m, err := realfarm.ReadManifest(*file)
	if err != nil {
		fail(err)
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fail(err)
	}
	sprites := map[string]sprite{}
	index := filepath.Join(*out, "sprites.json")
	if data, err := os.ReadFile(index); err == nil {
		json.Unmarshal(data, &sprites)
	}
	if *crops {
		bakeCrops(*out, sprites)
		writeIndex(index, sprites)
		return
	}
	fe := realfarm.Fetcher()
	for _, a := range m.Assets {
		if *only != "" && a.ID != *only {
			continue
		}
		model, err := load(fe.Cache, a)
		if err != nil {
			fail(fmt.Errorf("%s: %w", a.ID, err))
		}
		if a.Kind == "ground" {
			for i, p := range patches {
				s := bake.Ground(model, 0, p, options())
				name := fmt.Sprintf("%s_g%d", a.ID, i)
				if err := writePNG(filepath.Join(*out, name+".png"), s); err != nil {
					fail(err)
				}
				b := s.Image.Bounds()
				sprites[name] = sprite{W: b.Dx(), H: b.Dy(), PivotX: s.Pivot.X, PivotY: s.Pivot.Y}
			}
			fmt.Println(a.ID, "baked")
			continue
		}
		for _, sh := range shots(model) {
			s := bake.Render(model, sh.Shot, options())
			if err := writePNG(filepath.Join(*out, a.ID+sh.suffix+".png"), s); err != nil {
				fail(err)
			}
			b := s.Image.Bounds()
			sprites[a.ID+sh.suffix] = sprite{W: b.Dx(), H: b.Dy(), PivotX: s.Pivot.X, PivotY: s.Pivot.Y}
		}
		fmt.Println(a.ID, "baked")
	}
	writeIndex(index, sprites)
}

// cropSeeds are how many plants of each crop and stage are baked, so a
// field of them is not one plant over and over.
const cropSeeds = 3

// bakeCrops bakes every crop flora grows, at every stage, a few plants of
// each, into out.
func bakeCrops(out string, sprites map[string]sprite) {
	for _, c := range flora.Crops {
		for stage := range c.Stages {
			for seed := range cropSeeds {
				s := bake.Render(c.Grow(stage, uint32(seed+1)), bake.Shot{Yaw: float64(seed) * 120}, options())
				name := fmt.Sprintf("crop_%s_s%d_%d", c.Name, stage, seed)
				if err := writePNG(filepath.Join(out, name+".png"), s); err != nil {
					fail(err)
				}
				b := s.Image.Bounds()
				sprites[name] = sprite{W: b.Dx(), H: b.Dy(), PivotX: s.Pivot.X, PivotY: s.Pivot.Y}
			}
		}
		fmt.Println(c.Name, "grown")
	}
}

// shot is a picture to bake of a model, and the suffix of its name.
type shot struct {
	bake.Shot
	suffix string
}

// shots are the pictures baked of m: each heading of the whole model, or,
// for a model of variants side by side, each variant alone at two.
func shots(m *gltf.Model) []shot {
	var out []shot
	if len(m.Roots) == 1 {
		for i, yaw := range headings {
			out = append(out, shot{bake.Shot{Yaw: yaw}, fmt.Sprintf("_%d", i)})
		}
		return out
	}
	for r := range m.Roots {
		for i, yaw := range headings[:2] {
			out = append(out, shot{bake.Shot{Yaw: yaw, Only: true, Root: r}, fmt.Sprintf("_v%d_%d", r, i)})
		}
	}
	return out
}

// load fetches asset a's files, checked, and reads its model.
func load(cache string, a realfarm.Asset) (*gltf.Model, error) {
	fe := realfarm.Fetcher()
	for _, f := range a.Files {
		if _, err := fe.Get(f); err != nil {
			return nil, err
		}
	}
	root, err := os.OpenRoot(cache)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return gltf.Open(root, path.Clean(a.Model), gltf.DefaultLimits)
}

// options are the farm's light: a morning sun from the upper left.
func options() bake.Options {
	o := bake.DefaultOptions(realfarm.Projection)
	o.Unit = realfarm.Cell
	return o
}

func writePNG(name string, s bake.Sprite) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	if err := png.Encode(f, s.Image); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// writeIndex writes sprites as JSON; encoding/json sorts a map's keys.
func writeIndex(name string, sprites map[string]sprite) {
	data, _ := json.MarshalIndent(sprites, "", "  ")
	if err := os.WriteFile(name, append(data, '\n'), 0o644); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "realart:", err)
	os.Exit(1)
}
