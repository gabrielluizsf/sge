package human

import (
	"fmt"
	"image"
	"strings"

	"github.com/gabrielluizsf/sge/gltf"
)

// dresser builds a person: the body shaped, the parts made so far, and
// the textures read, by name, as indices of the model's images.
type dresser struct {
	kit    *Kit
	pos    [][3]float64
	model  *gltf.Model
	parts  []part
	images map[string]int
	grey   float64 // how grey their hair is
}

// part is a mesh of the person in MakeHuman's units, before it is stood
// on the ground: its vertices, texture places, influences, triangles and
// material.
type part struct {
	pos    [][3]float64
	uv     [][2]float64
	infl   [][]influence
	tris   []int32
	smooth []int32 // the vertex each is, for smooth normals across seams
	mat    int
}

// dress makes the body, in skin of l's, and every proxy l wears.
func (d *dresser) dress(b Body, l Look) error {
	proxies := make([]*proxy, len(l.Wear))
	hide := map[int32]bool{}
	for i, name := range l.Wear {
		data, err := d.kit.lim.read(d.kit.fsys, name)
		if err != nil {
			return err
		}
		if proxies[i], err = parseProxy(data, len(d.pos)); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		for _, v := range proxies[i].hides {
			hide[v] = true
		}
	}
	skin, err := d.skin(b, l.Skin)
	if err != nil {
		return err
	}
	d.body(hide, skin)
	for i, p := range proxies {
		if err := d.wear(l.Wear[i], p, l.Eyes); err != nil {
			return fmt.Errorf("%s: %w", l.Wear[i], err)
		}
	}
	return nil
}

// skin is the body's material: of the .mhmat name, or blended for b.
func (d *dresser) skin(b Body, name string) (int, error) {
	if name != "" {
		return d.material(name)
	}
	img, err := d.kit.blendSkin(skins(b))
	if err != nil {
		return 0, err
	}
	return d.textured("skin", material{color: [3]float64{1, 1, 1}}, maps{d.image("skin", img), -1, -1}), nil
}

// body is the base mesh's body, not its helpers, less the faces of
// vertices hidden under what is worn.
func (d *dresser) body(hide map[int32]bool, mat int) {
	m := d.kit.base
	var faces []face
	for _, f := range m.faces {
		if m.groups[f.group] != "body" {
			continue
		}
		shown := true
		for _, v := range f.v {
			shown = shown && !hide[v]
		}
		if shown {
			faces = append(faces, f)
		}
	}
	d.parts = append(d.parts, split(surface{faces, m.uv, d.pos, func(v int32) []influence { return d.kit.infl[v] }, mat}))
}

// wear fits the proxy in file name to the body, in its own material or,
// for eyes, a proxy of the eyes folder, those of eyes.
func (d *dresser) wear(name string, p *proxy, eyes string) error {
	objName, err := beside(name, p.obj)
	if err != nil {
		return err
	}
	data, err := d.kit.lim.read(d.kit.fsys, objName)
	if err != nil {
		return err
	}
	m, err := parseOBJ(data)
	if err != nil {
		return err
	}
	if len(m.pos) != len(p.refs) {
		return fmt.Errorf("human: %d vertices placed for a mesh of %d", len(p.refs), len(m.pos))
	}
	matName := eyes
	if eyes == "" || !strings.HasPrefix(name, "eyes/") {
		if matName, err = beside(name, p.material); err != nil {
			return err
		}
	}
	mat, err := d.material(matName)
	if err != nil {
		return err
	}
	if greys(name) {
		d.greyOut(mat)
	}
	at := p.fit(d.pos)
	infl := func(v int32) []influence { r := p.refs[v]; return mixed(d.kit.infl, r.v, r.w) }
	d.parts = append(d.parts, split(surface{m.faces, m.uv, at, infl, mat}))
	return nil
}

// material is the .mhmat in file name as a material of the model.
func (d *dresser) material(name string) (int, error) {
	data, err := d.kit.lim.read(d.kit.fsys, name)
	if err != nil {
		return 0, err
	}
	m := parseMaterial(data)
	tex := maps{-1, -1, -1}
	for _, t := range []struct {
		file string
		to   *int
	}{{m.diffuse, &tex.color}, {m.normal, &tex.normal}, {m.occlusion, &tex.occlusion}} {
		if t.file == "" {
			continue
		}
		file, err := beside(name, t.file)
		if err != nil {
			return 0, err
		}
		if *t.to, err = d.load(file); err != nil {
			return 0, err
		}
	}
	return d.textured(name, m, tex), nil
}

// load is the texture in file as an image of the model, read once.
func (d *dresser) load(file string) (int, error) {
	if i, ok := d.images[file]; ok {
		return i, nil
	}
	img, err := d.kit.texture(file)
	if err != nil {
		return 0, err
	}
	return d.image(file, img), nil
}

func (d *dresser) image(name string, img image.Image) int {
	d.model.Images = append(d.model.Images, img)
	d.images[name] = len(d.model.Images) - 1
	return d.images[name]
}

// maps are a material's textures, indices of the model's images: its
// colour, normal and occlusion maps, -1 for none.
type maps struct{ color, normal, occlusion int }

// textured adds material m, named name, of textures tex: m's colour alone
// where it has no colour map.
func (d *dresser) textured(name string, m material, tex maps) int {
	c := m.color
	if tex.color >= 0 {
		c = [3]float64{1, 1, 1}
	}
	g := gltf.Material{Name: name, BaseColor: [4]float64{c[0], c[1], c[2], 1}, Texture: tex.color, MetalRough: -1,
		Normal: tex.normal, NormalScale: 1, Occlusion: tex.occlusion, OcclusionStrength: 1,
		Roughness: 0.6, AlphaMode: "OPAQUE", AlphaCutoff: 0.5, DoubleSided: m.twoSided}
	if m.transparent {
		g.AlphaMode, g.AlphaCutoff, g.DoubleSided = "MASK", 0.4, true
	}
	d.model.Materials = append(d.model.Materials, g)
	return len(d.model.Materials) - 1
}
