package render3d

import "github.com/gabrielluizsf/sge/pose"

// Renderer draws scenes on a device, frame after frame, reusing its
// memory: after the first frames, preparing one allocates nothing.
type Renderer struct {
	// Stats are what the last frame drew.
	Stats Stats

	dev       Device
	models    int
	materials int
	pipelines [kinds][2]Pipeline

	uniforms Buffer // the frame's uniforms, and its objects'
	capacity int    // how many floats uniforms holds
	staging  []float32
	shadow   Texture
	shadowAt int // the shadow map's size
	light    Camera

	visible []seen
	draws   []drawItem
	partAt  []int // scratch: where a part's objects are, for a group
	skinAt  []int // scratch: where a skin's palette is, for an instance
}

// NewRenderer is a renderer drawing on dev.
func NewRenderer(dev Device) *Renderer { return &Renderer{dev: dev} }

// Device is the device the renderer draws on.
func (r *Renderer) Device() Device { return r.dev }

// Draw draws s: Prepare, then Submit.
func (r *Renderer) Draw(s *Scene) error {
	r.Prepare(s)
	return r.Submit(s)
}

// pipeline is the pipeline of kind, two-sided or not, made the first time
// it is asked for.
func (r *Renderer) pipeline(kind Kind, double bool) (Pipeline, error) {
	d := 0
	if double {
		d = 1
	}
	if p := r.pipelines[kind][d]; p != 0 {
		return p, nil
	}
	p, err := r.dev.NewPipeline(PipelineDesc{Kind: kind, DoubleSided: double})
	r.pipelines[kind][d] = p
	return p, err
}

// reserve is where n floats of objects start, made room for at the end
// of the staging uniforms, on a multiple of UniformAlign.
func (r *Renderer) reserve(n int) int {
	at := len(r.staging)
	end := at + aligned(n)
	if end > cap(r.staging) {
		grown := make([]float32, at, max(2*cap(r.staging), end))
		copy(grown, r.staging)
		r.staging = grown
	}
	r.staging = r.staging[:end]
	return at
}

// put writes m into the staging uniforms at at.
func (r *Renderer) put(at int, m *pose.Mat) {
	s := r.staging[at : at+16]
	for i := range s {
		s[i] = float32(m[i])
	}
}
