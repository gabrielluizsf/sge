package soft

// clipped draws tri cut by the near plane, where depth is 0: what is
// behind it left out, the rest a fan of triangles.
func (p *pass) clipped(s *state, tri [3]vout) {
	in := 0
	for _, v := range tri {
		if v.clip[2] >= 0 {
			in++
		}
	}
	if in == 3 {
		p.raster(s, tri)
		return
	}
	if in == 0 {
		return
	}
	var poly [4]vout
	n := 0
	for i := range 3 {
		a, b := tri[i], tri[(i+1)%3]
		if a.clip[2] >= 0 {
			poly[n] = a
			n++
		}
		if (a.clip[2] >= 0) != (b.clip[2] >= 0) {
			poly[n] = mix(a, b, a.clip[2]/(a.clip[2]-b.clip[2]))
			n++
		}
	}
	for i := 1; i+1 < n; i++ {
		p.raster(s, [3]vout{poly[0], poly[i], poly[i+1]})
	}
}

// mix is the vertex f of the way from a to b.
func mix(a, b vout, f float64) vout {
	var o vout
	for k := range 4 {
		o.clip[k] = a.clip[k] + (b.clip[k]-a.clip[k])*f
	}
	for k := range 3 {
		o.world[k] = a.world[k] + (b.world[k]-a.world[k])*f
		o.normal[k] = a.normal[k] + (b.normal[k]-a.normal[k])*f
	}
	for k := range 2 {
		o.uv[k] = a.uv[k] + (b.uv[k]-a.uv[k])*f
	}
	return o
}
