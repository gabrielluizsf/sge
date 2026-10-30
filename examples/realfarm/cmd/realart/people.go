package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// Each motion is baked facing eight ways, a frame at a time.
const (
	facings = 8
	frames  = 8
)

// water is the pond a swimmer is baked in: dark, clear a hand deep.
var water = &bake.Water{Color: [3]float64{0.01, 0.06, 0.07}, Clear: 0.35}

// frame is one picture of a villager to bake: its name and shot.
type frame struct {
	name string
	shot bake.Shot
	wet  bool
}

// bakePeople bakes every villager, made as humans.json pins, at every
// stage of life, in its motions, facing every way, frame by frame, into
// out, on every core; only that villager, if only is not "".
func bakePeople(out string, sprites map[string]sprite, only string) {
	h, err := realfarm.ReadHumans("humans.json")
	if err != nil {
		fail(err)
	}
	kit, done, err := h.Kit(realfarm.Fetcher())
	if err != nil {
		fail(err)
	}
	defer done()
	var mu sync.Mutex
	record := func(name string, s bake.Sprite) {
		b := s.Image.Bounds()
		mu.Lock()
		sprites[name] = sprite{W: b.Dx(), H: b.Dy(), PivotX: s.Pivot.X, PivotY: s.Pivot.Y}
		mu.Unlock()
	}
	for _, v := range realfarm.Villagers {
		if only != "" && v.ID != only {
			continue
		}
		if err := onlyCC0(kit, v); err != nil {
			fail(err)
		}
		for _, g := range v.Stages() {
			bakeStage(stageOf{kit, v, g}, out, record)
		}
		fmt.Println(v.ID, "baked")
	}
}

// stageOf is a villager at a stage of life, made from kit.
type stageOf struct {
	kit *human.Kit
	v   realfarm.Villager
	g   realfarm.Grown
}

// bakeStage bakes a villager at a stage of life, in its motions, every
// frame on every core, into out, recording each sprite.
func bakeStage(st stageOf, out string, record func(string, bake.Sprite)) {
	m, err := st.kit.Make(st.v.At(st.g.Years), st.v.Look)
	if err == nil {
		err = human.Animate(m, st.g.Motions...)
	}
	if err != nil {
		fail(fmt.Errorf("%s as a %s: %w", st.v.ID, st.g.Stage, err))
	}
	jobs := make(chan frame)
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range jobs {
				s := shoot(m, f)
				if err := writePNG(filepath.Join(out, f.name+".png"), s); err != nil {
					fail(err)
				}
				record(f.name, s)
			}
		}()
	}
	for _, f := range framesOf(st.v.ID+"_"+st.g.Stage.String(), st.g.Motions) {
		jobs <- f
	}
	close(jobs)
	wg.Wait()
}

// framesOf are every frame of id, a villager at a stage of life, in
// motions: person_<id>_<stage>_<motion>_d<facing>_f<frame>.
func framesOf(id string, motions []human.Motion) []frame {
	var out []frame
	for _, mo := range motions {
		for d := range facings {
			for k := range frames {
				out = append(out, frame{
					name: fmt.Sprintf("person_%s_%s_d%d_f%d", id, mo.Name, d, k),
					shot: bake.Shot{Animation: mo.Name, Time: mo.Length * float64(k) / frames, Yaw: float64(d) * 360 / facings},
					wet:  mo.Name == human.Swim.Name || mo.Name == human.Tread.Name,
				})
			}
		}
	}
	return out
}

// shoot bakes frame f of m, in the water if f is wet.
func shoot(m *gltf.Model, f frame) bake.Sprite {
	o := options()
	if f.wet {
		o.Water = water
	}
	return bake.Render(m, f.shot, o)
}

// onlyCC0 refuses a villager made of any file that does not say it is
// CC0, so a game sold with them owes no one: a pack's name is no promise.
func onlyCC0(kit *human.Kit, v realfarm.Villager) error {
	lics, err := kit.Licences(v.Body, v.Look)
	if err != nil {
		return fmt.Errorf("%s: %w", v.ID, err)
	}
	for name, lic := range lics {
		if lic != "CC0" {
			return fmt.Errorf("%s: %s is under %q, not CC0", v.ID, name, lic)
		}
	}
	return nil
}
