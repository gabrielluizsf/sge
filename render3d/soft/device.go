// Package soft is a render3d.Device drawn by the CPU, a pixel at a time:
// slow, but the same on every machine, Go's or Deno's, for tests and
// golden pictures of what the GPU backends must draw. It rasterises
// triangles with a depth buffer, shades them by Lambert's law and the
// ambient light, shadowed by the sun's map, and samples textures at
// their nearest pixel, at their full size whatever mipmaps are asked.
package soft

import (
	"fmt"

	"github.com/gabrielluizsf/sge/render3d"
)

// Device is a software device, its screen a picture in memory.
type Device struct {
	screen    *texture
	buffers   []*buffer // by number; 0 is none
	textures  []*texture
	pipelines []*render3d.PipelineDesc
}

// buffer is a buffer's floats, or its indices.
type buffer struct {
	floats  []float32
	indices []uint32
}

var _ render3d.Device = (*Device)(nil)

// New is a device whose screen is w by h pixels.
func New(w, h int) *Device {
	d := &Device{buffers: []*buffer{nil}, textures: []*texture{nil}, pipelines: []*render3d.PipelineDesc{nil}}
	d.Resize(w, h)
	return d
}

// Size is the screen's size.
func (d *Device) Size() (w, h int) { return d.screen.w, d.screen.h }

// Resize makes the screen w by h pixels, cleared.
func (d *Device) Resize(w, h int) {
	d.screen = &texture{w: max(w, 1), h: max(h, 1), srgb: true}
	d.screen.target()
}

// Release lets go of a buffer, texture or pipeline: its number is not
// given again, and using it is an error.
func (d *Device) Release(h render3d.Handle) {
	switch h := h.(type) {
	case render3d.Buffer:
		if int(h) < len(d.buffers) {
			d.buffers[h] = nil
		}
	case render3d.Texture:
		if int(h) < len(d.textures) && h != render3d.Screen {
			d.textures[h] = nil
		}
	case render3d.Pipeline:
		if int(h) < len(d.pipelines) {
			d.pipelines[h] = nil
		}
	}
}

// NewBuffer is a buffer holding a copy of b's data.
func (d *Device) NewBuffer(b render3d.BufferDesc) (render3d.Buffer, error) {
	buf := &buffer{}
	switch {
	case b.Kind == render3d.Indices:
		buf.indices = append([]uint32(nil), b.Indices...)
	case b.Floats != nil:
		buf.floats = append([]float32(nil), b.Floats...)
	default:
		if b.Size < 0 {
			return 0, fmt.Errorf("soft: a buffer of %d floats", b.Size)
		}
		buf.floats = make([]float32, b.Size)
	}
	d.buffers = append(d.buffers, buf)
	return render3d.Buffer(len(d.buffers) - 1), nil
}

// Write copies data into b's floats from at on, as many as fit.
func (d *Device) Write(b render3d.Buffer, at int, data []float32) {
	buf := d.buffer(b)
	if buf == nil || at < 0 || at > len(buf.floats) {
		return
	}
	copy(buf.floats[at:], data)
}

// NewPipeline is a pipeline of p's kind.
func (d *Device) NewPipeline(p render3d.PipelineDesc) (render3d.Pipeline, error) {
	if p.Kind > render3d.Unlit {
		return 0, fmt.Errorf("soft: no pipeline kind %d", p.Kind)
	}
	d.pipelines = append(d.pipelines, &p)
	return render3d.Pipeline(len(d.pipelines) - 1), nil
}

// buffer is buffer b, or nil.
func (d *Device) buffer(b render3d.Buffer) *buffer {
	if int(b) >= len(d.buffers) {
		return nil
	}
	return d.buffers[b]
}

// texture is texture t, the screen for Screen, or nil.
func (d *Device) texture(t render3d.Texture) *texture {
	if t == render3d.Screen {
		return d.screen
	}
	if int(t) >= len(d.textures) {
		return nil
	}
	return d.textures[t]
}
