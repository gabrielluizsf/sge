package human

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
	"testing/fstest"
)

// tinyOBJ is a body of two quads, one on the other, and a helper triangle.
const tinyOBJ = `# a tiny body
v -1 0 0
v 1 0 0
v 1 10 0
v -1 10 0
v -1 20 0
v 1 20 0
v 0 30 0
vt 0 0
vt 1 0
vt 1 1
vt 0 1
g body
f 1/1 2/2 3/3 4/4
f 4/1 3/2 6/3 5/4
g helper-hair
f 5 6 7
`

// tinySkeleton is a root bone from the feet up, and a spine above it.
const tinySkeleton = `{
 "bones": {
  "root": {"head": "root____head", "tail": "root____tail", "parent": null},
  "spine": {"head": "root____tail", "tail": "spine____tail", "parent": "root"}
 },
 "joints": {"root____head": [0, 1], "root____tail": [2, 3], "spine____tail": [4, 5]}
}`

const tinyWeights = `{"weights": {"root": [[0, 1], [1, 1], [2, 0.5], [3, 0.5]], "spine": [[2, 0.5], [3, 0.5], [4, 1], [5, 1], [6, 1]]}}`

// tinyKit is a kit of the tiny body, a target that grows it as a man,
// a skin and a shirt.
func tinyKit(t *testing.T) fstest.MapFS {
	t.Helper()
	return fstest.MapFS{
		baseMesh:     {Data: []byte(tinyOBJ)},
		baseSkeleton: {Data: []byte(tinySkeleton)},
		baseWeights:  {Data: []byte(tinyWeights)},
		// A man of any ancestry, young, is 2 units taller.
		"targets/macrodetails/universal-male-young-averagemuscle-averageweight.target": {Data: []byte("# taller\n4 0 2 0\n5 0 2 0\n6 0 2 0\n")},
		"skins/plain/plain.mhmat":   {Data: []byte("name plain\ndiffuseTexture plain.png\n")},
		"skins/plain/plain.png":     {Data: pngOf(t, 4, color.NRGBA{200, 150, 120, 255})},
		"clothes/shirt/shirt.mhclo": {Data: []byte("name shirt\nobj_file shirt.obj\nmaterial shirt.mhmat\nx_scale 0 1 2.0\ndelete_verts\n5\nverts 0\n2 3 2 0.5 0.5 0 0 0 1\n3\n2\n")},
		"clothes/shirt/shirt.obj":   {Data: []byte("v 0 0 0\nv 0 0 0\nv 0 0 0\nf 1 2 3\n")},
		"clothes/shirt/shirt.mhmat": {Data: []byte("name shirt\ndiffuseColor 0.2 0.3 0.8\n")},
	}
}

// pngOf is a PNG of size by size pixels of one colour.
func pngOf(t *testing.T, size int, c color.NRGBA) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = c.R, c.G, c.B, c.A
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// loaded is the tiny kit, loaded.
func loaded(t *testing.T, fsys fstest.MapFS) *Kit {
	t.Helper()
	k, err := Load(fsys, DefaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	return k
}
