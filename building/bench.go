package building

// Bench is a plain bench of boards, 1.6 m long, its seat at the height
// of a knee, on two slab legs braced by a stretcher.
var Bench = Building{Name: "bench", Width: 1.6, Depth: 0.34, build: func(b *builder) {
	boards := paint{of: Planks}
	b.block(box{lo: vec{-0.8, 0.42, -0.17}, hi: vec{0.8, 0.47, 0.17}, paint: boards, grain: vec{1, 0, 0}})
	for _, x := range []float64{-0.6, 0.6} {
		b.block(box{lo: vec{x - 0.025, 0, -0.15}, hi: vec{x + 0.025, 0.42, 0.15}, paint: boards})
	}
	b.lay(beam{a: vec{-0.6, 0.18, 0}, b: vec{0.6, 0.18, 0}, w: 0.04, h: 0.08, paint: paint{of: Timber}})
}}
