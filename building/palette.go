package building

import (
	"image"
	"image/color"
)

// Material is what a part of a building is made of.
type Material int

// The materials buildings are made of.
const (
	Plaster   Material = iota // rendered walls, painted
	Brick                     // brick walls and chimneys
	Stone                     // stone walls, plinths, wells and forges
	Planks                    // boarded walls, doors, shutters, counters
	Timber                    // posts, beams and rails, rough sawn
	Tiles                     // clay roof tiles
	Thatch                    // reed thatch
	Slates                    // slate roof
	Iron                      // wrought iron: anvils, hinges, hoops
	Glass                     // window panes
	Water                     // a well's or trough's water
	Soot                      // a chimney's flue, a forge's coals
	Canvas                    // a stall's striped awning
	Flat                      // a plain colour, as a paint's tint says
	materials                 // how many there are
)

// Surface is how a material looks: a PBR texture set, as Poly Haven's
// are, its colour, its glTF metal-rough map (occlusion in red, roughness
// in green, metal in blue), and its normal map, nil where it has none;
// how many metres one repeat of it covers, and whether its boards or
// grain run across the image, not down it; its colour, in sRGB, times
// the texture's; and its roughness and metalness, times the map's.
type Surface struct {
	Albedo, ARM, Normal image.Image
	Span                float64
	Across              bool
	Color               color.NRGBA
	Roughness, Metallic float64
}

// Palette is the surface of each material; a material missing from it is
// plain, as Plain paints it.
type Palette map[Material]Surface

// surface is m's surface in the palette, or its plain one.
func (p Palette) surface(m Material) Surface {
	if s, ok := p[m]; ok {
		if s.Span <= 0 {
			s.Span = 1
		}
		return s
	}
	if m == Canvas {
		return canvas()
	}
	return plain[m]
}

// Plain is every material as a colour of its own, untextured, but for
// the canvas's stripes, which are painted: a palette that needs nothing
// fetched.
func Plain() Palette {
	p := Palette{Canvas: canvas()}
	for m, s := range plain {
		p[m] = s
	}
	return p
}

// canvas is the awning's canvas, its stripes painted.
func canvas() Surface {
	return Surface{Albedo: stripes(), Span: 1.2, Color: color.NRGBA{255, 255, 255, 255}, Roughness: 0.9}
}

// plain are the materials' plain colours, and how rough each is.
var plain = map[Material]Surface{
	Plaster: {Span: 1, Color: color.NRGBA{226, 218, 200, 255}, Roughness: 0.9},
	Brick:   {Span: 1, Color: color.NRGBA{150, 72, 52, 255}, Roughness: 0.85},
	Stone:   {Span: 1, Color: color.NRGBA{138, 132, 120, 255}, Roughness: 0.9},
	Planks:  {Span: 1, Color: color.NRGBA{120, 92, 64, 255}, Roughness: 0.8},
	Timber:  {Span: 1, Color: color.NRGBA{92, 68, 46, 255}, Roughness: 0.85},
	Tiles:   {Span: 1, Color: color.NRGBA{170, 82, 50, 255}, Roughness: 0.7},
	Thatch:  {Span: 1, Color: color.NRGBA{170, 140, 88, 255}, Roughness: 0.95},
	Slates:  {Span: 1, Color: color.NRGBA{78, 82, 90, 255}, Roughness: 0.6},
	Iron:    {Span: 1, Color: color.NRGBA{58, 58, 62, 255}, Roughness: 0.45, Metallic: 1},
	Glass:   {Span: 1, Color: color.NRGBA{40, 52, 60, 255}, Roughness: 0.08},
	Water:   {Span: 1, Color: color.NRGBA{22, 34, 36, 255}, Roughness: 0.05},
	Soot:    {Span: 1, Color: color.NRGBA{30, 27, 25, 255}, Roughness: 1},
	Flat:    {Span: 1, Color: color.NRGBA{255, 255, 255, 255}, Roughness: 0.7},
}
