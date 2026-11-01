//go:build linux && amd64

package gl

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/gabrielluizsf/sge/render3d"
)

// defines are what kind k defines in the shaders.
var defines = [...]string{
	render3d.Static:        "",
	render3d.Skinned:       "#define SKINNED\n",
	render3d.Shadow:        "#define SHADOW\n",
	render3d.ShadowSkinned: "#define SKINNED\n#define SHADOW\n",
	render3d.Unlit:         "#define UNLIT\n",
}

// program is kind k's program, compiled and linked the first time.
func (d *Device) program(k render3d.Kind) (uint32, error) {
	if p := d.programs[k]; p != 0 {
		return p, nil
	}
	g := &d.gl
	head := "#version 420 core\n" + defines[k]
	vs, err := d.shader(glVertexShader, head+vertexShader)
	if err != nil {
		return 0, err
	}
	defer call(g.deleteShader, n(vs))
	fs, err := d.shader(glFragmentShader, head+fragmentShader)
	if err != nil {
		return 0, err
	}
	defer call(g.deleteShader, n(fs))
	p := uint32(call(g.createProgram))
	call(g.attachShader, n(p), n(vs))
	call(g.attachShader, n(p), n(fs))
	call(g.linkProgram, n(p))
	if !d.status(status{g.getProgramiv, g.getProgramInfoLog, p, glLinkStatus}) {
		err := fmt.Errorf("gl: linking kind %d: %s", k, d.infoLog(status{g.getProgramiv, g.getProgramInfoLog, p, glLinkStatus}))
		call(g.deleteProgram, n(p))
		return 0, err
	}
	d.programs[k] = p
	return p, nil
}

// shader is a shader of kind compiled from src.
func (d *Device) shader(kind uint32, src string) (uint32, error) {
	g := &d.gl
	s := uint32(call(g.createShader, n(kind)))
	text := [1]unsafe.Pointer{cstring(src)}
	var pin runtime.Pinner
	pin.Pin(text[0])
	call(g.shaderSource, n(s), n(1), ptr(unsafe.Pointer(&text)), n(0))
	pin.Unpin()
	call(g.compileShader, n(s))
	st := status{g.getShaderiv, g.getShaderInfoLog, s, glCompileStatus}
	if !d.status(st) {
		err := fmt.Errorf("gl: compiling a shader: %s", d.infoLog(st))
		call(g.deleteShader, n(s))
		return 0, err
	}
	return s, nil
}
