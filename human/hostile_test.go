package human

import (
	"strings"
	"testing"
	"testing/fstest"
)

// After CVE-2020-28589, an OBJ index outside the mesh, or a relative one,
// is refused, not read past the end.
func TestOBJIndicesOutsideTheMeshAreRefused(t *testing.T) {
	for name, obj := range map[string]string{
		"negative":      "v 0 0 0\nv 1 0 0\nv 0 1 0\nf -1 -2 -3\n",
		"zero":          "v 0 0 0\nv 1 0 0\nv 0 1 0\nf 0 1 2\n",
		"past the end":  "v 0 0 0\nv 1 0 0\nv 0 1 0\nf 1 2 4\n",
		"texture past":  "v 0 0 0\nv 1 0 0\nv 0 1 0\nvt 0 0\nf 1/1 2/2 3/1\n",
		"huge":          "v 0 0 0\nf 1 1 99999999999\n",
		"too many":      "v 0 0 0\nf" + strings.Repeat(" 1", maxCorners+1) + "\n",
		"two corners":   "v 0 0 0\nf 1 1\n",
		"not a number":  "v 0 x 0\n",
		"not finite":    "v 0 NaN 0\n",
		"some textures": "v 0 0 0\nvt 0 0\nf 1/1 1 1\n",
	} {
		if _, err := parseOBJ([]byte(obj)); err == nil {
			t.Errorf("%s: read", name)
		}
	}
}

// A target, skeleton, weights or proxy naming a vertex outside the body
// is refused.
func TestVerticesOutsideTheBodyAreRefused(t *testing.T) {
	const n = 7
	for name, data := range map[string]string{"negative": "-1 0 0 0\n", "past": "7 0 0 0\n", "short": "1 0 0\n"} {
		if _, err := parseTarget([]byte(data), n); err == nil {
			t.Errorf("target %s: read", name)
		}
	}
	bad := map[string]string{
		"joint past":     strings.Replace(tinySkeleton, "[4, 5]", "[4, 7]", 1),
		"joint negative": strings.Replace(tinySkeleton, "[4, 5]", "[4, -1]", 1),
		"no joints":      strings.Replace(tinySkeleton, "[4, 5]", "[]", 1),
		"joint missing":  strings.Replace(tinySkeleton, `"spine____tail": [4, 5]`, `"x": [4]`, 1),
		"cycle":          strings.Replace(tinySkeleton, `"parent": null`, `"parent": "spine"`, 1),
		"parent missing": strings.Replace(tinySkeleton, `"parent": "root"`, `"parent": "nobody"`, 1),
	}
	for name, data := range bad {
		if _, err := parseSkeleton([]byte(data), n); err == nil {
			t.Errorf("skeleton %s: read", name)
		}
	}
	s, err := parseSkeleton([]byte(tinySkeleton), n)
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{
		"bone missing": `{"weights": {"tail": [[0, 1]]}}`,
		"vertex past":  `{"weights": {"root": [[7, 1]]}}`,
		"vertex half":  `{"weights": {"root": [[1.5, 1]]}}`,
		"too much":     `{"weights": {"root": [[1, 2]]}}`,
		"negative":     `{"weights": {"root": [[1, -0.5]]}}`,
	} {
		if _, err := parseWeights([]byte(data), s, n); err == nil {
			t.Errorf("weights %s: read", name)
		}
	}
	for name, data := range map[string]string{
		"ref past":    "obj_file a.obj\nverts 0\n1 2 7 0.3 0.3 0.4 0 0 0\n",
		"single past": "obj_file a.obj\nverts 0\n9\n",
		"scale past":  "obj_file a.obj\nx_scale 0 70 1\nverts 0\n1\n",
		"scale zero":  "obj_file a.obj\nx_scale 0 1 0\nverts 0\n1\n",
		"hide past":   "obj_file a.obj\ndelete_verts\n1 - 70\nverts 0\n1\n",
		"hide back":   "obj_file a.obj\ndelete_verts\n5 - 2\nverts 0\n1\n",
		"no mesh":     "verts 0\n1\n",
		"odd ref":     "obj_file a.obj\nverts 0\n1 2\n",
	} {
		if _, err := parseProxy([]byte(data), n); err == nil {
			t.Errorf("proxy %s: read", name)
		}
	}
	p, err := parseProxy([]byte("obj_file a.obj\ndelete_verts\n1 - 3 5\nverts 0\n1\n"), n)
	if err != nil || len(p.hides) != 4 {
		t.Errorf("delete_verts 1 - 3 5: %v %v", p, err)
	}
}

// After zip slip, no name a pack gives leads out of the kit: not up past
// its root, not from its root, and a file over the limit is not read.
func TestNoNameLeadsOutOfTheKit(t *testing.T) {
	for _, rel := range []string{"../../../../etc/passwd", "/etc/passwd", "../x.png"} {
		if name, err := beside("skins/plain/plain.mhmat", rel); err == nil && rel != "../x.png" {
			t.Errorf("%q from skins/plain: %q", rel, name)
		}
	}
	if name, err := beside("eyes/low-poly/low-poly.mhclo", "../materials/brown.mhmat"); err != nil || name != "eyes/materials/brown.mhmat" {
		t.Errorf("a neighbour folder: %q %v", name, err)
	}
	fsys := tinyKit(t)
	fsys["skins/evil/evil.mhmat"] = &fstest.MapFile{Data: []byte("diffuseTexture ../../../../../../etc/passwd\n")}
	k := loaded(t, fsys)
	if _, err := k.Make(Average, Look{Skin: "skins/evil/evil.mhmat"}); err == nil || !strings.Contains(err.Error(), "out of the kit") {
		t.Errorf("a texture out of the kit: %v", err)
	}
	if _, err := k.lim.read(fsys, "../base.obj"); err == nil {
		t.Error("read a name climbing out")
	}
	small := Limits{Bytes: 10, Pixels: DefaultLimits.Pixels, Texture: 64}
	if _, err := small.read(fsys, baseMesh); err == nil || !strings.Contains(err.Error(), "over") {
		t.Errorf("read a file over the limit: %v", err)
	}
}

// A texture larger than the limit is refused from its header, before its
// pixels are made, and one kept is shrunk.
func TestTexturesAreBounded(t *testing.T) {
	fsys := tinyKit(t)
	fsys["big.png"] = &fstest.MapFile{Data: pngOf(t, 64, colorOf(10))}
	k := loaded(t, fsys)
	k.lim.Pixels = 63 * 63
	if _, err := k.texture("big.png"); err == nil {
		t.Error("decoded a texture over the limit")
	}
	k.lim.Pixels, k.lim.Texture = 64*64, 16
	img, err := k.texture("big.png")
	if err != nil || img.Rect.Dx() != 16 || img.NRGBAAt(3, 3) != colorOf(10) {
		t.Errorf("shrunk: %v %v", img.Rect, err)
	}
}
