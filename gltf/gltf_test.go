package gltf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// triangle is a GLB of one textured triangle on a node moved up by 2: its
// positions, normals, UVs, indices and a PNG, in its binary chunk.
func triangle() []byte {
	bin := floats(0, 0, 0, 1, 0, 0, 0, 1, 0) // 36 bytes of positions
	bin = append(bin, floats(0, 0, 1, 0, 0, 1, 0, 0, 1)...)
	bin = append(bin, floats(0, 0, 1, 0, 0, 1)...) // 24 bytes of UVs
	idx := []byte{0, 0, 1, 0, 2, 0, 0, 0}          // 3 uint16, padded
	bin = append(bin, idx...)
	img := pngOf(2, 2)
	bin = append(bin, img...)
	js := fmt.Sprintf(`{"asset":{"version":"2.0"},"scene":0,"scenes":[{"nodes":[0]}],
"nodes":[{"name":"tri","mesh":0,"translation":[0,2,0]}],
"meshes":[{"primitives":[{"attributes":{"POSITION":0,"NORMAL":1,"TEXCOORD_0":2},"indices":3,"material":0}]}],
"materials":[{"pbrMetallicRoughness":{"baseColorFactor":[1,0.5,0.5,1],"baseColorTexture":{"index":0},"roughnessFactor":0.8}}],
"textures":[{"source":0}],"images":[{"bufferView":4,"mimeType":"image/png"}],
"accessors":[{"bufferView":0,"componentType":5126,"count":3,"type":"VEC3"},{"bufferView":1,"componentType":5126,"count":3,"type":"VEC3"},
{"bufferView":2,"componentType":5126,"count":3,"type":"VEC2"},{"bufferView":3,"componentType":5123,"count":3,"type":"SCALAR"}],
"bufferViews":[{"buffer":0,"byteOffset":0,"byteLength":36},{"buffer":0,"byteOffset":36,"byteLength":36},{"buffer":0,"byteOffset":72,"byteLength":24},
{"buffer":0,"byteOffset":96,"byteLength":6},{"buffer":0,"byteOffset":104,"byteLength":%d}],
"buffers":[{"byteLength":%d}]}`, len(img), len(bin))
	return glb(js, bin)
}

func TestATexturedTriangleIsRead(t *testing.T) {
	m, err := Decode(triangle(), DefaultLimits)
	if err != nil {
		t.Fatal(err)
	}
	p := m.Meshes[0].Primitives[0]
	if len(p.Positions) != 3 || p.Positions[1] != [3]float32{1, 0, 0} || p.Normals[2] != [3]float32{0, 0, 1} || p.UVs[2] != [2]float32{0, 1} {
		t.Errorf("vertices %v %v %v", p.Positions, p.Normals, p.UVs)
	}
	if fmt.Sprint(p.Indices) != "[0 1 2]" || p.Material != 0 {
		t.Errorf("indices %v material %d", p.Indices, p.Material)
	}
	mat := m.Materials[0]
	if mat.BaseColor != [4]float64{1, 0.5, 0.5, 1} || mat.Texture != 0 || mat.Roughness != 0.8 || mat.AlphaMode != "OPAQUE" {
		t.Errorf("material %+v", mat)
	}
	if b := m.Images[0].Bounds(); b.Dx() != 2 || b.Dy() != 2 {
		t.Errorf("image %v", b)
	}
	if len(m.Roots) != 1 || m.Nodes[0].Translation != [3]float64{0, 2, 0} || m.Nodes[0].Rotation != [4]float64{0, 0, 0, 1} {
		t.Errorf("nodes %+v roots %v", m.Nodes, m.Roots)
	}
}

// edit is the triangle's GLB with its JSON changed by replacing old with
// new, the binary chunk as it was.
func edit(t *testing.T, old, new string) []byte {
	t.Helper()
	data := triangle()
	js, bin, err := splitGLB(data)
	if err != nil {
		t.Fatal(err)
	}
	s := string(js)
	if !strings.Contains(s, old) {
		t.Fatalf("no %q in the model", old)
	}
	return glb(strings.Replace(s, old, new, 1), bin)
}

// Every offset, count and index a model gives is checked: an accessor or
// view past its data, sizes that overflow, overlapping elements, an index
// past the vertices, sparse accessors and a broken scene graph are
// refused, not read out of bounds.
func TestMalformedModelsAreRefused(t *testing.T) {
	for _, c := range []struct{ name, old, new string }{
		{"accessor past its view", `"count":3,"type":"VEC3"},{"bufferView":1`, `"count":4,"type":"VEC3"},{"bufferView":1`},
		{"count overflowing", `"count":3,"type":"VEC3"},{"bufferView":1`, `"count":9223372036854775807,"type":"VEC3"},{"bufferView":1`},
		{"negative count", `"count":3,"type":"VEC2"`, `"count":-3,"type":"VEC2"`},
		{"offset past its view", `{"bufferView":0,"componentType":5126`, `{"bufferView":0,"byteOffset":40,"componentType":5126`},
		{"view past its buffer", `"byteOffset":36,"byteLength":36}`, `"byteOffset":36,"byteLength":3600}`},
		{"view length overflowing", `"byteOffset":36,"byteLength":36}`, `"byteOffset":9223372036854775807,"byteLength":36}`},
		{"stride under an element", `{"buffer":0,"byteOffset":0,"byteLength":36}`, `{"buffer":0,"byteOffset":0,"byteLength":36,"byteStride":4}`},
		{"index past the vertices", `"componentType":5123,"count":3`, `"componentType":5123,"count":3,"byteOffset":0,"normalized":false,"min":[0],"max":[9]`},
		{"sparse accessor", `{"bufferView":2,"componentType":5126,"count":3,"type":"VEC2"}`, `{"bufferView":2,"componentType":5126,"count":3,"type":"VEC2","sparse":{"count":1}}`},
		{"missing accessor", `"POSITION":0`, `"POSITION":99`},
		{"missing material", `"material":0`, `"material":5`},
		{"missing texture image", `"textures":[{"source":0}]`, `"textures":[{"source":7}]`},
		{"node its own child", `"name":"tri",`, `"name":"tri","children":[0],`},
		{"scene root missing", `"scenes":[{"nodes":[0]}]`, `"scenes":[{"nodes":[3]}]`},
	} {
		data := edit(t, c.old, c.new)
		if c.name == "index past the vertices" {
			// Point the indices at bytes holding 9: past three vertices.
			data = indexNine(t, data)
		}
		if _, err := Decode(data, DefaultLimits); !errors.Is(err, errFormat) {
			t.Errorf("%s: %v", c.name, err)
		}
	}
}

// indexNine sets the triangle's second index to 9.
func indexNine(t *testing.T, data []byte) []byte {
	t.Helper()
	js, bin, _ := splitGLB(data)
	bin = append([]byte(nil), bin...)
	binary.LittleEndian.PutUint16(bin[98:], 9)
	return glb(string(js), bin)
}

// A broken container is refused: a length that lies, a chunk past the
// end, no JSON.
func TestBrokenGLBsAreRefused(t *testing.T) {
	good := triangle()
	lying := append([]byte(nil), good...)
	binary.LittleEndian.PutUint32(lying[8:], uint32(len(good)+100))
	long := append([]byte(nil), good...)
	binary.LittleEndian.PutUint32(long[12:], 1<<30)
	for name, data := range map[string][]byte{"length": lying, "chunk": long, "truncated": good[:10], "empty": glb("", nil)[:12]} {
		if _, err := Decode(data, DefaultLimits); err == nil {
			t.Errorf("%s: read", name)
		}
	}
}

// An image whose header claims more pixels than allowed is refused before
// its pixels are allocated, as a decompression bomb would make them.
func TestAnImageBombIsRefused(t *testing.T) {
	img := pngOf(64, 64)
	if _, err := decodeImage(img, 64*64); err != nil {
		t.Fatalf("a small image: %v", err)
	}
	// The IHDR claims 1,000,000 by 1,000,000 pixels; the CRC is not
	// checked by DecodeConfig's first look, but the size is refused first.
	bomb := append([]byte(nil), img...)
	binary.BigEndian.PutUint32(bomb[16:], 1000000)
	binary.BigEndian.PutUint32(bomb[20:], 1000000)
	if _, err := decodeImage(bomb, DefaultLimits.Pixels); err == nil {
		t.Error("a bomb was decoded")
	}
	if _, err := decodeImage(img, 64*63); err == nil {
		t.Error("an image over the limit was decoded")
	}
}

// The files a model names are read only inside its directory: no "..",
// absolute path, scheme or link leads out, and a model read from bytes
// may name none.
func TestModelsReadOnlyTheirOwnFiles(t *testing.T) {
	dir := t.TempDir()
	models := filepath.Join(dir, "models")
	os.MkdirAll(models, 0o755)
	secret := filepath.Join(dir, "secret.bin")
	os.WriteFile(secret, make([]byte, 36), 0o644)
	os.Symlink(secret, filepath.Join(models, "link.bin"))
	os.WriteFile(filepath.Join(models, "ok.bin"), floats(0, 0, 0, 1, 0, 0, 0, 1, 0), 0o644)
	gltfWith := func(uri string) string {
		return `{"asset":{"version":"2.0"},"meshes":[{"primitives":[{"attributes":{"POSITION":0}}]}],
"accessors":[{"bufferView":0,"componentType":5126,"count":3,"type":"VEC3"}],
"bufferViews":[{"buffer":0,"byteLength":36}],"buffers":[{"uri":"` + uri + `","byteLength":36}]}`
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	for i, uri := range []string{"../secret.bin", secret, "file:///" + secret, "link.bin", `..\secret.bin`, "%2e%2e/secret.bin"} {
		name := fmt.Sprintf("models/bad%d.gltf", i)
		os.WriteFile(filepath.Join(dir, name), []byte(gltfWith(uri)), 0o644)
		if _, err := Open(root, name, DefaultLimits); err == nil {
			t.Errorf("%q was read from outside the model's directory", uri)
		}
	}
	os.WriteFile(filepath.Join(models, "good.gltf"), []byte(gltfWith("ok.bin")), 0o644)
	if m, err := Open(root, "models/good.gltf", DefaultLimits); err != nil || len(m.Meshes[0].Primitives[0].Positions) != 3 {
		t.Errorf("a model's own file: %v", err)
	}
	if _, err := Decode([]byte(gltfWith("ok.bin")), DefaultLimits); err == nil {
		t.Error("a model read from bytes read a file")
	}
}

// Nesting deeper than the limit is refused.
func TestDeepNodesAreRefused(t *testing.T) {
	var nodes []string
	for i := range 10 {
		nodes = append(nodes, fmt.Sprintf(`{"children":[%d]}`, i+1))
	}
	nodes = append(nodes, `{}`)
	js := `{"asset":{"version":"2.0"},"nodes":[` + strings.Join(nodes, ",") + `]}`
	if _, err := Decode([]byte(js), Limits{Bytes: 1 << 20, Pixels: 1, Depth: 5}); !errors.Is(err, errFormat) {
		t.Errorf("10 levels under a limit of 5: %v", err)
	}
	if _, err := Decode([]byte(js), DefaultLimits); err != nil {
		t.Errorf("10 levels: %v", err)
	}
}

// Decoding any bytes returns a model or an error, and never panics.
func FuzzDecode(f *testing.F) {
	f.Add(triangle())
	f.Add([]byte(`{"asset":{"version":"2.0"}}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		Decode(data, Limits{Bytes: 1 << 20, Pixels: 1 << 16, Depth: 16})
	})
}
