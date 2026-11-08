package gltf

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
)

// Open reads the model name, a .glb or .gltf file in root, with lim
// bounding what it may allocate: the files it names are read from root
// too, so none outside it.
func Open(root *os.Root, name string, lim Limits) (*Model, error) {
	f, err := root.Open(name)
	if err != nil {
		return nil, err
	}
	data, err := readLimited(f, lim.Bytes)
	f.Close()
	if err != nil {
		return nil, err
	}
	dir := path.Dir(name)
	sub := root
	if dir != "." {
		if sub, err = root.OpenRoot(dir); err != nil {
			return nil, err
		}
		defer sub.Close()
	}
	return decode(data, source{root: sub, limits: lim})
}

// Decode reads a model from data, a GLB or glTF JSON, which may name no
// file, only data: URIs.
func Decode(data []byte, lim Limits) (*Model, error) {
	if int64(len(data)) > lim.Bytes {
		return nil, malformed("model over %d bytes", lim.Bytes)
	}
	return decode(data, source{limits: lim})
}

func decode(data []byte, src source) (*Model, error) {
	js, bin := data, []byte(nil)
	if bytes.HasPrefix(data, []byte("glTF")) {
		var err error
		if js, bin, err = splitGLB(data); err != nil {
			return nil, err
		}
	}
	var doc document
	if err := json.Unmarshal(js, &doc); err != nil {
		return nil, malformed("JSON: %v", err)
	}
	v := views{doc: &doc}
	for i, b := range doc.Buffers {
		var data []byte
		switch {
		case b.URI == "" && i == 0 && bin != nil:
			data = bin
		case b.URI == "":
			return nil, malformed("buffer %d has no data", i)
		default:
			var err error
			if data, err = src.read(b.URI); err != nil {
				return nil, err
			}
		}
		if b.ByteLength < 0 || b.ByteLength > int64(len(data)) {
			return nil, malformed("buffer %d is %d bytes, not %d", i, len(data), b.ByteLength)
		}
		v.buffers = append(v.buffers, data[:b.ByteLength])
	}
	return build(&doc, v, src)
}

// build reads everything doc describes, through its views and source.
func build(doc *document, v views, src source) (*Model, error) {
	m := &Model{}
	steps := []func() error{
		func() (err error) { m.Images, err = images(doc, v, src); return },
		func() (err error) { m.Materials, err = materials(doc, len(m.Images)); return },
		func() (err error) { m.Meshes, err = meshes(doc, v, len(m.Materials)); return },
		func() (err error) { m.Nodes, m.Roots, err = nodes(doc, src.limits.Depth); return },
		func() (err error) { m.Skins, err = skins(doc, v); return },
		func() (err error) { m.Animations, err = animations(doc, v); return },
	}
	for _, s := range steps {
		if err := s(); err != nil {
			return nil, err
		}
	}
	return m, nil
}
