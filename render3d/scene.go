package render3d

// Sun is the light of a scene: Toward, the direction to the sun; its
// Color and the Ambient light's, linear; and its shadow, a map Shadow
// pixels square (0 for none) covering Area, offset by Bias in depth.
type Sun struct {
	Toward         [3]float64
	Color, Ambient [3]float64
	Shadow         int
	Area           Sphere
	Bias           float64
}

// Scene is what a frame draws: its instances, seen by its camera, lit by
// its sun, over its Clear colour, linear.
type Scene struct {
	Camera    Camera
	Sun       Sun
	Instances []*Instance
	Clear     [4]float32
}

// Stats are what the last frame drew: the instances shown and culled,
// the draws, and their triangles.
type Stats struct {
	Shown, Culled, Draws, Triangles int
}
