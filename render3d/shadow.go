package render3d

import "math"

// sunCamera is the sun's view of its shadow's area: orthographic, from
// far enough toward the sun to see what stands up to three radii away
// from the area cast its shadow into it.
func sunCamera(s *Sun) Camera {
	r := max(s.Area.Radius, 1e-3)
	t := unit(s.Toward)
	if t == [3]float64{} {
		t = [3]float64{0, 0, 1}
	}
	c := s.Area.Center
	eye := [3]float64{c[0] + t[0]*4*r, c[1] + t[1]*4*r, c[2] + t[2]*4*r}
	upward := [3]float64{0, 0, 1}
	if math.Abs(t[2]) > 0.99 {
		upward = [3]float64{0, 1, 0}
	}
	return Orthographic(OrthographicDesc{Eye: eye, Target: c, Up: upward, Width: 2 * r, Height: 2 * r, Near: 0, Far: 5 * r})
}

// shadowMap is the depth texture of the sun's shadow, made, or made
// again, at size pixels square; 0 when size is 0.
func (r *Renderer) shadowMap(size int) (Texture, error) {
	if size == r.shadowAt {
		return r.shadow, nil
	}
	if r.shadow != 0 {
		r.dev.Release(r.shadow)
		r.shadow, r.shadowAt = 0, 0
	}
	if size <= 0 {
		return 0, nil
	}
	t, err := r.dev.NewTexture(TextureDesc{Width: size, Height: size, Format: Depth})
	if err != nil {
		return 0, err
	}
	r.shadow, r.shadowAt = t, size
	return t, nil
}
