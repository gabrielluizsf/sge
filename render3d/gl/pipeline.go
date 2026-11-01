//go:build linux && amd64

package gl

import (
	"fmt"

	"github.com/gabrielluizsf/sge/render3d"
)

// pipeline is a pipeline's program, its kind, and whether it draws both
// faces of a triangle.
type pipeline struct {
	program uint32
	kind    render3d.Kind
	double  bool
	stride  int // its vertices' floats
}

// NewPipeline is a pipeline of p's kind, its program compiled once for
// every pipeline of the kind.
func (d *Device) NewPipeline(p render3d.PipelineDesc) (render3d.Pipeline, error) {
	if p.Kind > render3d.Unlit {
		return 0, fmt.Errorf("gl: no pipeline kind %d", p.Kind)
	}
	var err error
	d.do(func() {
		var prog uint32
		if prog, err = d.program(p.Kind); err != nil {
			return
		}
		stride := render3d.StaticFloats
		if p.Kind == render3d.Skinned || p.Kind == render3d.ShadowSkinned {
			stride = render3d.SkinnedFloats
		}
		d.pipelines = append(d.pipelines, pipeline{program: prog, kind: p.Kind, double: p.DoubleSided, stride: stride})
	})
	if err != nil {
		return 0, err
	}
	return render3d.Pipeline(len(d.pipelines) - 1), nil
}

// pipeline is pipeline p, or nil.
func (d *Device) pipeline(p render3d.Pipeline) *pipeline {
	if int(p) >= len(d.pipelines) || d.pipelines[p].program == 0 {
		return nil
	}
	return &d.pipelines[p]
}
