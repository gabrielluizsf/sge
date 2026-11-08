package human

import "testing"

// Of every top and bottom a person may wear together, the outer layer lies
// clear of the inner one: no vertex of a top is inside the bottom's
// surface by more than a centimetre, so nothing pokes through at the
// waist.
func TestLayersLieClearOfEachOther(t *testing.T) {
	k := realKit(t)
	for _, sex := range []Sex{Man, Woman} {
		b := baselines[0].body
		if sex == Woman {
			b = baselines[1].body
		}
		tops, bottoms := suited(tops, sex, 30, 0), suited(bottoms, sex, 30, 0)
		for _, top := range tops {
			for _, bottom := range bottoms {
				check(t, k, b, []string{top.item, bottom.item})
			}
		}
	}
}

// What is worn under a garment does not come through it: a waistband or the
// fold of a pair of trousers stood out through a shirt's hem.
func TestWhatIsUnderAGarmentStaysUnder(t *testing.T) {
	k := realKit(t)
	for _, sex := range []Sex{Man, Woman} {
		b := baselines[0].body
		if sex == Woman {
			b = baselines[1].body
		}
		for _, top := range suited(tops, sex, 30, 0) {
			for _, bottom := range suited(bottoms, sex, 30, 0) {
				throughTop(t, k, b, []string{top.item, bottom.item})
			}
		}
	}
}

// throughTop fails if the vertices of the bottom wear[1] that the top wear[0]
// covers lie out through it.
func throughTop(t *testing.T, k *Kit, b Body, wear []string) {
	t.Helper()
	d, err := k.dressed(b, Look{Eyes: "eyes/materials/brown.mhmat", Wear: wear})
	if err != nil {
		t.Fatal(wear, err)
	}
	data, _ := k.lim.read(k.fsys, wear[0])
	top, _ := parseProxy(data, len(d.pos))
	g := newGrid(d.skinOf(0, top, d.kit.base.bodyNormals(d.pos)))
	bad := 0
	for _, q := range d.fitted[1].at {
		if on, out, ok := g.closest(q); ok && covers(q, on, out) && dot(sub(q, on), out) > 0 {
			bad++
		}
	}
	if bad > 0 {
		t.Errorf("%v: %d vertices of the bottom come through the top", wear, bad)
	}
}

// check makes the person and tests every garment against each under it.
func check(t *testing.T, k *Kit, b Body, wear []string) {
	t.Helper()
	d, err := k.dressed(b, Look{Eyes: "eyes/materials/brown.mhmat", Wear: wear})
	if err != nil {
		t.Fatal(wear, err)
	}
	proxies := make([]*proxy, len(wear))
	normals := d.kit.base.bodyNormals(d.pos)
	for i, name := range wear {
		data, _ := k.lim.read(k.fsys, name)
		proxies[i], _ = parseProxy(data, len(d.pos))
	}
	for i := range wear {
		for j := range wear {
			if layerOf(wear[j]) == 0 || layerOf(wear[j]) >= layerOf(wear[i]) || layerOf(wear[i]) == 0 {
				continue
			}
			g := newGrid(d.skinOf(j, proxies[j], normals))
			for v, q := range d.fitted[i].at {
				on, out, ok := g.closest(q)
				if !ok {
					continue
				}
				if h := dot(sub(q, on), out); h < -0.1 {
					t.Errorf("%s + %s: vertex %d is %.4f dm inside", wear[i], wear[j], v, h)
					break
				}
			}
		}
	}
}
