package render3d

import "image"

// null is a device that holds and draws nothing, a screen w by h, for
// testing what runs on the CPU; it counts the draws it is given.
type null struct {
	w, h  int
	n     uint32
	draws int
}

func (d *null) next() uint32                               { d.n++; return d.n }
func (d *null) NewBuffer(BufferDesc) (Buffer, error)       { return Buffer(d.next()), nil }
func (d *null) Write(Buffer, int, []float32)               {}
func (d *null) NewTexture(TextureDesc) (Texture, error)    { return Texture(d.next()), nil }
func (d *null) NewPipeline(PipelineDesc) (Pipeline, error) { return Pipeline(d.next()), nil }
func (d *null) Release(Handle)                             {}
func (d *null) Begin(PassDesc) Pass                        { return nullPass{d} }
func (d *null) Read(Texture) (*image.NRGBA, error) {
	return image.NewNRGBA(image.Rect(0, 0, d.w, d.h)), nil
}
func (d *null) Size() (w, h int) { return d.w, d.h }
func (d *null) Resize(w, h int)  { d.w, d.h = w, h }

type nullPass struct{ d *null }

func (nullPass) SetPipeline(Pipeline) {}
func (nullPass) SetVertices(Buffer)   {}
func (nullPass) SetIndices(Buffer)    {}
func (nullPass) SetBindings(Bindings) {}
func (p nullPass) Draw(Draw)          { p.d.draws++ }
func (nullPass) End()                 {}
