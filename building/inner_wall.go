package building

import "fmt"

// The back walls inside: how thick they are, from the cell's edge in, how
// tall the skirting board along their foot, and the window's opening in
// a north wall, across and up.
const (
	innerThick       = 0.14
	skirtH, skirtOut = 0.1, 0.02
	paneX            = 0.3
	paneLo, paneHi   = 0.95, 1.85
)

// innerWalls are the back walls of house style k's inside, h, a storey
// tall, in its walls' paint: along a cell's north edge, the same with a
// window, along its west, and the corner where both meet; a
// map's rows run along -Z and its columns along +X.
func innerWalls(k int, h House) []Building {
	tall := h.Storey
	if tall <= 0 {
		tall = 2.8
	}
	wall := func(name string, north, west, window bool) Building {
		return Building{Name: fmt.Sprint(name, k), Width: 1, Depth: 1, build: func(b *builder) {
			s := innerSide{b: b, h: h, tall: tall}
			if north {
				s.north(window)
			}
			if west {
				s.west()
			}
		}}
	}
	return []Building{wall("wall_n", true, false, false), wall("wall_nwin", true, false, true), wall("wall_w", false, true, false), wall("wall_nw", true, true, false)}
}

// innerSide is a back wall being built: its builder, its house's style,
// and how tall it stands.
type innerSide struct {
	b    *builder
	h    House
	tall float64
}

func (s innerSide) paint() paint { return paint{of: s.h.Wall, tint: s.h.Tint} }

// north is the wall along the cell's north edge, its back, which a map
// baked as the ground is lays at +Z: its face to -Z, a window in it when
// window is set.
func (s innerSide) north(window bool) {
	z0, z1 := 0.5-innerThick, 0.5
	p := s.paint()
	if !window {
		s.b.block(box{lo: vec{-0.5, 0, z0}, hi: vec{0.5, s.tall, z1}, paint: p})
	} else {
		s.b.block(box{lo: vec{-0.5, 0, z0}, hi: vec{0.5, paneLo, z1}, paint: p})
		s.b.block(box{lo: vec{-0.5, paneHi, z0}, hi: vec{0.5, s.tall, z1}, paint: p})
		s.b.block(box{lo: vec{-0.5, paneLo, z0}, hi: vec{-paneX, paneHi, z1}, paint: p})
		s.b.block(box{lo: vec{paneX, paneLo, z0}, hi: vec{0.5, paneHi, z1}, paint: p})
		s.window(z0, z1)
	}
	s.b.block(box{lo: vec{-0.5, 0, z0 - skirtOut}, hi: vec{0.5, skirtH, z0}, paint: paint{of: Timber}, grain: vec{1, 0, 0}})
}

// west is the wall along the cell's west edge, its face to +X.
func (s innerSide) west() {
	x0, x1 := -0.5, -0.5+innerThick
	s.b.block(box{lo: vec{x0, 0, -0.5}, hi: vec{x1, s.tall, 0.5}, paint: s.paint()})
	s.b.block(box{lo: vec{x1, 0, -0.5}, hi: vec{x1 + skirtOut, skirtH, 0.5}, paint: paint{of: Timber}, grain: vec{0, 0, 1}})
}

// window is the pane in a north wall from z0, its face, to z1, its frame
// round it in the house's trim and its sill under it, reaching into the
// room.
func (s innerSide) window(z0, z1 float64) {
	mid := (z0 + z1) / 2
	s.b.block(box{lo: vec{-paneX, paneLo, mid - 0.005}, hi: vec{paneX, paneHi, mid + 0.005}, paint: paint{of: Glass}})
	trim := coat(Timber, s.h.Trim)
	const f = 0.05 // the frame's width
	in, out := z0-0.01, z0+0.04
	s.b.block(box{lo: vec{-paneX, paneLo, in}, hi: vec{paneX, paneLo + f, out}, paint: trim})
	s.b.block(box{lo: vec{-paneX, paneHi - f, in}, hi: vec{paneX, paneHi, out}, paint: trim})
	s.b.block(box{lo: vec{-paneX, paneLo + f, in}, hi: vec{-paneX + f, paneHi - f, out}, paint: trim})
	s.b.block(box{lo: vec{paneX - f, paneLo + f, in}, hi: vec{paneX, paneHi - f, out}, paint: trim})
	s.b.block(box{lo: vec{-0.02, paneLo + f, z0 + 0.03}, hi: vec{0.02, paneHi - f, mid}, paint: trim})
	s.b.block(box{lo: vec{-paneX - 0.06, paneLo - 0.04, z0 - 0.08}, hi: vec{paneX + 0.06, paneLo, z0}, paint: trim, grain: vec{1, 0, 0}})
}
