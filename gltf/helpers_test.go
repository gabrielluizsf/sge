package gltf

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"math"
)

// glb is a GLB of the JSON js and the binary chunk bin.
func glb(js string, bin []byte) []byte {
	pad := func(b []byte, with byte) []byte {
		for len(b)%4 != 0 {
			b = append(b, with)
		}
		return b
	}
	j := pad([]byte(js), ' ')
	bin = pad(append([]byte(nil), bin...), 0)
	var b bytes.Buffer
	w := func(v uint32) { binary.Write(&b, binary.LittleEndian, v) }
	total := 12 + 8 + len(j)
	if len(bin) > 0 {
		total += 8 + len(bin)
	}
	w(glbMagic)
	w(glbVersion)
	w(uint32(total))
	w(uint32(len(j)))
	w(chunkJSON)
	b.Write(j)
	if len(bin) > 0 {
		w(uint32(len(bin)))
		w(chunkBIN)
		b.Write(bin)
	}
	return b.Bytes()
}

// floats is the little-endian bytes of fs.
func floats(fs ...float32) []byte {
	var b []byte
	for _, f := range fs {
		b = binary.LittleEndian.AppendUint32(b, math.Float32bits(f))
	}
	return b
}

// pngOf is a w by h PNG of one colour.
func pngOf(w, h int) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, color.NRGBA{200, 100, 50, 255})
		}
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	return b.Bytes()
}
