// Package render3d draws glTF models in 3D, live, frame after frame: a
// model is uploaded once, and then its instances, each where its game
// puts it, posed by its animation, are culled, given the detail their
// size on screen calls for, sorted and drawn, with the sun's shadow.
//
// Drawing goes through a Device, a small interface shaped after WebGPU,
// which an OpenGL, a WebGPU or a software backend implements: package
// render3d/soft is the last, slow but exact, for tests on any machine.
// What runs on the CPU every frame, and what a backend must do with it,
// is fixed here: the layouts of vertices and uniforms in layout.go, and
// what each pipeline Kind computes in descs.go.
//
// The world is the game's: Z up, a glTF model's Y turned onto it as it
// is uploaded. Clip space is WebGPU's: X right and Y up from -1 to 1, and
// depth from 0, nearest, to 1.
package render3d

import "image"

// Device is a graphics card, or what stands for one: it holds buffers,
// textures and pipelines, and draws with them in passes.
type Device interface {
	// NewBuffer is a buffer of d's kind and data.
	NewBuffer(d BufferDesc) (Buffer, error)
	// Write puts data in a buffer of floats from its float at on.
	Write(b Buffer, at int, data []float32)
	// NewTexture is a texture of d's size and pixels.
	NewTexture(d TextureDesc) (Texture, error)
	// NewPipeline is the pipeline of d's built-in shaders.
	NewPipeline(d PipelineDesc) (Pipeline, error)
	// Release lets go of a buffer, texture or pipeline.
	Release(h Handle)
	// Begin starts a pass drawing into d's targets, cleared as d says.
	Begin(d PassDesc) Pass
	// Read is a target's pixels, top row first, sRGB: the screen's
	// (Screen) or a texture's; for tests and screenshots, not each frame.
	Read(t Texture) (*image.NRGBA, error)
	// Size is the screen's size in pixels; Resize changes it, as a
	// window's does.
	Size() (w, h int)
	Resize(w, h int)
}

// Pass is a pass drawing into its targets: the state set stays until set
// again, and each Draw draws with what is set.
type Pass interface {
	SetPipeline(p Pipeline)
	// SetVertices is the vertex buffer, in its pipeline's layout;
	// SetIndices the index buffer, three a triangle.
	SetVertices(b Buffer)
	SetIndices(b Buffer)
	SetBindings(b Bindings)
	// Draw draws d's triangles, as many times as d's instances.
	Draw(d Draw)
	// End ends the pass: its drawing done before the next pass begins.
	End()
}
