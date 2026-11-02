package backpack

// Point is a pixel of the screen.
type Point struct{ X, Y int }

// Rect is a part of the screen, in pixels.
type Rect struct{ X, Y, W, H int }

// In reports whether p is inside r.
func (r Rect) In(p Point) bool {
	return p.X >= r.X && p.X < r.X+r.W && p.Y >= r.Y && p.Y < r.Y+r.H
}

// View is the screen the pack is laid out on: its size in pixels, how many
// pixels a pixel of the game's art takes, and the least width of the
// panel, for its words.
type View struct{ W, H, Scale, MinW int }

// The pack's measures, in pixels of art.
const (
	SlotSize = 20
	slotGap  = 2
	padding  = 6
	titleH   = 14
	hintH    = 12
	// barMargin is the room between the hotbar and the screen's bottom.
	barMargin = 8
)

// rects is where the pack's parts lie.
type rects struct {
	view  View
	panel Rect
	grid  []Rect
	bar   []Rect
}

// Place lays the pack out on v: the hotbar in a row at the bottom, the
// screen's panel with its grid above it, in the middle.
func (p *Pack) Place(v View) {
	s := v.Scale
	step := (SlotSize + slotGap) * s
	r := rects{view: v}
	barW := len(p.Hotbar)*step - slotGap*s
	for i := range p.Hotbar {
		r.bar = append(r.bar, Rect{(v.W-barW)/2 + i*step, v.H - SlotSize*s - barMargin, SlotSize * s, SlotSize * s})
	}
	cells := p.cells()
	rows := cells / p.Cols
	gridW := p.Cols*step - slotGap*s
	w := max(gridW+2*padding*s, v.MinW)
	h := titleH*s + rows*step - slotGap*s + padding*s + hintH*s
	r.panel = Rect{(v.W - w) / 2, max(barMargin, (v.H-SlotSize*s-barMargin-h)/2), w, h}
	for i := range cells {
		r.grid = append(r.grid, Rect{r.panel.X + (w-gridW)/2 + (i%p.Cols)*step, r.panel.Y + titleH*s + (i/p.Cols)*step, SlotSize * s, SlotSize * s})
	}
	p.rects = r
}

// spotAt is the slot under the pixel at, if any.
func (p *Pack) spotAt(at Point) (Spot, bool) {
	for i, r := range p.rects.bar {
		if r.In(at) {
			return Spot{Bar, i}, true
		}
	}
	for i, r := range p.rects.grid {
		if r.In(at) {
			return Spot{Grid, i}, true
		}
	}
	return Spot{}, false
}

// rectOf is where slot s lies.
func (p *Pack) rectOf(s Spot) Rect {
	list := p.rects.grid
	if s.Area == Bar {
		list = p.rects.bar
	}
	if s.Index < 0 || s.Index >= len(list) {
		return Rect{}
	}
	return list[s.Index]
}
