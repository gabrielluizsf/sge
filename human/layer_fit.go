package human

import "fmt"

// fitted is a proxy's mesh and where it lies on the body, and its own
// outward normals once it is worn.
type fitted struct {
	mesh *mesh
	at   [][3]float64
}

// fitAll reads every proxy's mesh and fits it to the body: eyebrows stand
// off the skin, and each layer of clothes lies out of the ones under it,
// so nothing sinks into the flesh or another garment.
func (d *dresser) fitAll(proxies []*proxy) error {
	d.fitted = make([]fitted, len(proxies))
	normals := d.kit.base.bodyNormals(d.pos)
	for i, p := range proxies {
		name := d.look.Wear[i]
		objName, err := beside(name, p.obj)
		if err != nil {
			return err
		}
		data, err := d.kit.lim.read(d.kit.fsys, objName)
		if err != nil {
			return err
		}
		m, err := parseOBJ(data)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		if len(m.pos) != len(p.refs) {
			return fmt.Errorf("%s: human: %d vertices placed for a mesh of %d", name, len(p.refs), len(m.pos))
		}
		at := p.fit(d.pos)
		if hasFolder(name, "eyebrows") {
			p.standOff(at, d.pos, normals)
		}
		d.fitted[i] = fitted{mesh: m, at: at}
	}
	d.layer(proxies, normals)
	return nil
}

// layer pushes each garment out of the garments worn under it.
func (d *dresser) layer(proxies []*proxy, normals [][3]float64) {
	for i := range proxies {
		if layerOf(d.look.Wear[i]) != 2 {
			continue
		}
		for j := range proxies {
			if layerOf(d.look.Wear[j]) == 1 && near(d.fitted[i].at, d.fitted[j].at) {
				newGrid(d.skinOf(j, proxies[j], normals)).pushOut(d.fitted[i].at)
				newGrid(d.skinOf(i, proxies[i], normals)).pullIn(d.fitted[j].at)
			}
		}
	}
}

// skinOf is the surface of the garment worn at j.
func (d *dresser) skinOf(j int, p *proxy, normals [][3]float64) *skinOf {
	s := &skinOf{at: d.fitted[j].at, normals: make([][3]float64, len(p.refs))}
	for v, r := range p.refs {
		var n [3]float64
		for k := range 3 {
			n = add(n, scale(normals[r.v[k]], r.w[k]))
		}
		s.normals[v] = normalize(n)
	}
	for _, f := range d.fitted[j].mesh.faces {
		for k := 1; k+1 < len(f.v); k++ {
			s.tris = append(s.tris, [3]int32{f.v[0], f.v[k], f.v[k+1]})
		}
	}
	return s
}
