package human

import (
	"image"
	"sync"

	"github.com/gabrielluizsf/sge/gltf"
)

// Person is a made person whose colours can be changed after: Recolour
// retints their materials' textures, the mesh left as it is.
type Person struct {
	Model *gltf.Model
	kit   *Kit
	mats  []retint
	grey  float64

	mu sync.Mutex // guards mats, which Plan works on off the frame
}

// retint is how one material is recoloured: what it was before any
// tint, and what it is now.
type retint struct {
	open    bool // set by track, until settle
	mat     int
	role    role
	grey    bool        // greys with age
	src     image.Image // its texture before tinting, nil for none; planned: wantSrc
	orig    int         // its texture's index in the model before tinting, -1 if it had none there
	own     int         // the index of its tinted texture, -1 if it has none yet
	color   [4]float64  // its base colour before tinting
	want    key         // the tint last planned for it
	wantSrc image.Image // and the texture it is planned from
}

// key is the tint a material wears, to tell when it changes.
type key struct {
	ok      bool
	to      [3]float64
	keep    float64
	greying float64
	eyes    string // for the eyes, the material file they are of
}

func keyOf(t tint, ok bool) key {
	if !ok {
		return key{}
	}
	return key{ok: true, to: t.to, keep: t.keep}
}

// MakePerson is Make with the means to Recolour the person after.
func (k *Kit) MakePerson(b Body, l Look) (*Person, error) {
	d, err := k.dressed(b, l)
	if err != nil {
		return nil, err
	}
	return &Person{Model: d.done(), kit: k, mats: d.retints, grey: d.grey}, nil
}

// Recolouring is a person's colours changed, made but not yet applied:
// see Plan.
type Recolouring struct {
	p     *Person
	items []change
}

// change is one material's new texture or colour.
type change struct {
	at      int // an index of the person's mats
	src     image.Image
	img     image.Image
	color   [4]float64
	want    key         // the tint last planned for it
	wantSrc image.Image // and the texture it is planned from
}

// Plan works out what changing the person's skin tone, hair, eye and
// clothes colours to l's takes: the textures tinted, which is most of the
// time, and does not touch the model, so it may run while the model is
// drawn. Apply it to change the model. l must wear what the person
// wears, with the same body; only the eyes' material may differ.
func (p *Person) Plan(l Look) *Recolouring {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := &Recolouring{p: p}
	for i := range p.mats {
		r := &p.mats[i]
		t, ok := tintFor(r.role, l)
		k := keyOf(t, ok)
		if r.grey {
			k.greying = p.grey
		}
		src := r.wantSrc
		if r.role.kind == roleEyes {
			k.eyes = r.want.eyes
			if l.Eyes != "" && l.Eyes != k.eyes {
				if s, err := p.kit.eyesTexture(l.Eyes); err == nil {
					src, k.eyes = s, l.Eyes
				}
			}
		}
		if k == r.want {
			continue
		}
		r.want, r.wantSrc = k, src
		c := change{at: i, src: src, color: r.color}
		if src == nil {
			if ok {
				c.color = [4]float64{t.to[0], t.to[1], t.to[2], r.color[3]}
			}
		} else {
			c.img = src
			if n, isNRGBA := src.(*image.NRGBA); isNRGBA {
				if ok {
					n = tinted(n, t)
				}
				if r.grey && p.grey > 0 {
					n = greyed(n, p.grey)
				}
				c.img = n
			}
		}
		out.items = append(out.items, c)
	}
	return out
}

// Apply puts r's textures and colours on the model, which takes no time,
// and says which materials changed. The plans of a person are applied in
// the order they were made, each one, or the model ends up as an older
// plan had it.
func (p *Person) Apply(r *Recolouring) []int {
	p.mu.Lock()
	defer p.mu.Unlock()
	var changed []int
	for _, c := range r.items {
		rt := &p.mats[c.at]
		if c.src != rt.src {
			rt.orig = -1 // the texture it had is another's now
		}
		rt.src = c.src
		m := &p.Model.Materials[rt.mat]
		switch {
		case c.src == nil:
			m.BaseColor = c.color
		case rt.own >= 0:
			p.Model.Images[rt.own] = c.img
			m.Texture = rt.own
		case c.img == c.src && rt.orig >= 0:
			m.Texture = rt.orig
		default:
			p.Model.Images = append(p.Model.Images, c.img)
			rt.own = len(p.Model.Images) - 1
			m.Texture = rt.own
		}
		changed = append(changed, rt.mat)
	}
	return changed
}

// Recolour is Plan, then Apply.
func (p *Person) Recolour(l Look) []int { return p.Apply(p.Plan(l)) }
