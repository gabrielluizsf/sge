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
	look   Look
	// retints are the materials whose colours Person.Recolour changes.
	retints []retint
	// fitted are the proxies' meshes and places on the body, once worn.
	fitted []fitted
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
	d.look = l
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
	skin, err := d.skin(b)
	if err != nil {
		return err
	}
	d.body(hide, skin)
	if err := d.fitAll(proxies); err != nil {
		return err
	}
	for i, p := range proxies {
		if err := d.wear(i, p); err != nil {
			return fmt.Errorf("%s: %w", l.Wear[i], err)
		}
	}
	return nil
}

// skin is the body's material: of the look's .mhmat, or blended for b,
// in the look's tone, if it has one.
func (d *dresser) skin(b Body) (int, error) {
	if name := d.look.Skin; name != "" {
		mat, err := d.material(name)
		if err == nil {
			d.track(mat, role{kind: roleSkin}, false)
			if d.look.SkinTone > 0 {
				d.recolour(mat, skinTint(d.look.SkinTone))
			}
		}
		return mat, err
	}
	img, err := d.kit.blendSkin(skins(b))
	if err != nil {
		return 0, err
	}
	raw := img
	if d.look.SkinTone > 0 {
		img = tinted(img, skinTint(d.look.SkinTone))
	}
	mat := d.textured("skin", material{color: [3]float64{1, 1, 1}}, maps{d.image("skin", img), -1, -1})
	r := retint{mat: mat, role: role{kind: roleSkin}, src: raw, orig: -1, own: d.model.Materials[mat].Texture, color: d.model.Materials[mat].BaseColor}
	if d.look.SkinTone <= 0 {
		r.orig, r.own = r.own, -1
	}
	t, ok := tintFor(r.role, d.look)
	r.want, r.wantSrc = keyOf(t, ok), r.src
	d.retints = append(d.retints, r)
	return mat, nil
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

// wear dresses the person in the proxy the look wears at i, fitted to
// the body, in its own material or, for eyes, the look's, recoloured as
// the look asks.
func (d *dresser) wear(i int, p *proxy) error {
	name, eyes := d.look.Wear[i], d.look.Eyes
	m, at := d.fitted[i].mesh, d.fitted[i].at
	var err error
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
	d.track(mat, roleOf(name, i), greys(name))
	if t, ok := d.tintOf(i); ok {
		d.recolour(mat, t)
	}
	if greys(name) {
		d.greyOut(mat)
	}
	if insideMouth(name) {
		shadeMouth(&d.model.Materials[mat])
	}
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
