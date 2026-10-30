// Package human makes people from MakeHuman's CC0 assets: its base mesh
// shaped by a Body, woman to man, child to old, of any ancestry, build and
// height, as MakeHuman's own macro targets shape it; skinned to its
// skeleton; with a skin of their age, sex and ancestry; and dressed in
// proxies, clothes, hair, eyes and brows fitted to the body. A person is a
// gltf.Model, for bake to draw, and moves by Motions.
//
// A Kit reads MakeHuman's data, its "data" folder and asset packs, from an
// fs.FS: a folder, or a pack's zip read where it is, never unpacked. Its
// files are read as untrusted, after CVE-2020-28589 and zip slip: every
// index is checked against what it indexes, no name climbs out of the FS,
// and files and textures are bounded by Limits.
package human

import (
	"errors"
	"image"
	"io/fs"

	"github.com/gabrielluizsf/sge/gltf"
)

// Kit is MakeHuman's data: the base mesh, its skeleton and how its bones
// move it, and its targets and textures, read when first needed. A Kit is
// for one goroutine at a time.
type Kit struct {
	fsys     fs.FS
	lim      Limits
	base     *mesh
	skel     *skeleton
	infl     [][]influence
	targets  map[string]target
	pictures map[string]*image.NRGBA
}

// The files of a Kit, in MakeHuman's data folder.
const (
	baseMesh     = "3dobjs/base.obj"
	baseSkeleton = "rigs/default.mhskel"
	baseWeights  = "rigs/default_weights.mhw"
)

// Load reads MakeHuman's base mesh, skeleton and weights from fsys, which
// holds its data folder's files.
func Load(fsys fs.FS, lim Limits) (*Kit, error) {
	k := &Kit{fsys: fsys, lim: lim, targets: map[string]target{}, pictures: map[string]*image.NRGBA{}}
	data, err := lim.read(fsys, baseMesh)
	if err != nil {
		return nil, err
	}
	if k.base, err = parseOBJ(data); err != nil {
		return nil, err
	}
	n := len(k.base.pos)
	if data, err = lim.read(fsys, baseSkeleton); err != nil {
		return nil, err
	}
	if k.skel, err = parseSkeleton(data, n); err != nil {
		return nil, err
	}
	if data, err = lim.read(fsys, baseWeights); err != nil {
		return nil, err
	}
	if k.infl, err = parseWeights(data, k.skel, n); err != nil {
		return nil, err
	}
	return k, nil
}

// shape is the base mesh's vertices shaped as body b.
func (k *Kit) shape(b Body) ([][3]float64, error) {
	pos := append([][3]float64(nil), k.base.pos...)
	v := b.values()
	for _, name := range macroTargets() {
		w := weight(stem(name), v)
		if w < 1e-4 {
			continue
		}
		t, err := k.target("targets/" + name + ".target")
		if err != nil {
			return nil, err
		}
		t.apply(pos, w)
	}
	return pos, nil
}

// target is the target in file name, read once; one MakeHuman does not
// have is none, as it is to MakeHuman.
func (k *Kit) target(name string) (target, error) {
	if t, ok := k.targets[name]; ok {
		return t, nil
	}
	data, err := k.lim.read(k.fsys, name)
	if errors.Is(err, fs.ErrNotExist) {
		k.targets[name] = nil
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t, err := parseTarget(data, len(k.base.pos))
	if err != nil {
		return nil, err
	}
	k.targets[name] = t
	return t, nil
}

// Make is a person of body b, as l looks: their skin, eyes and what they
// wear, standing on the ground at the origin, facing +Z, in metres.
func (k *Kit) Make(b Body, l Look) (*gltf.Model, error) {
	pos, err := k.shape(b)
	if err != nil {
		return nil, err
	}
	d := &dresser{kit: k, pos: pos, model: &gltf.Model{}, images: map[string]int{}, grey: Greying(b.Years())}
	if err := d.dress(b, l); err != nil {
		return nil, err
	}
	return d.done(), nil
}
