package main

import (
	"os"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

// bakeAnimal builds the coats of a and bakes every motion of each, in
// every facing.
func bakeAnimal(root *os.Root, a animalDef, w *writer) error {
	coats, err := a.build(root)
	if err != nil {
		return err
	}
	in := &info{Unit: a.unit, Motions: map[string]motionInfo{}}
	w.info[a.name] = in
	o := bakes.Options()
	o.Unit = kit.Cell * a.unit
	for _, c := range coats {
		in.Coats = append(in.Coats, c.name)
		for _, m := range c.motions {
			in.Motions[m.name] = motionInfo{Frames: len(m.times), Length: length(m), Speed: m.speed, Once: m.once}
			for d := range facings {
				for f, t := range m.times {
					s := bake.Render(c.model, bake.Shot{Animation: m.clip, Time: t, Yaw: float64(d) * 360 / facings}, o)
					if err := w.put(frameName(a.name, c.name, m.name, d, f), s); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// length is the seconds a motion runs: for a loop, its period.
func length(m motion) float64 {
	n := len(m.times)
	if n < 2 {
		return 0
	}
	if m.once {
		return m.times[n-1]
	}
	return m.times[1] * float64(n)
}
