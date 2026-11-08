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

// rects is where the pack's parts lie: the screen's panel, with the grid or,
// on the Craft tab, its rows; the tabs on its top; the store's panel beside
// it with its slots; the hotbar.
type rects struct {
	view  View
	panel Rect
	grid  []Rect
	bar   []Rect
	tabs  []Rect
	rows  []Rect
	side  Rect
	chest []Rect
}

// tabH is the height of a tab, in pixels of art.
const tabH = 12

// sideGap is the room between the store's panel and the pack's.
const sideGap = 4

// Place lays the pack out on v: the hotbar in a row at the bottom, the
// screen's panel with its grid above it, in the middle, and the store's
// panel at its left.
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
	if p.Crafting() {
		rows = max(p.Rows, 1)
	}
	gridW := p.Cols*step - slotGap*s
	w := max(gridW+2*padding*s, v.MinW)
	h := titleH*s + rows*step - slotGap*s + padding*s + hintH*s
	sideW, sideH, siderows := 0, 0, 0
	if p.Side != nil {
		siderows = (len(p.Side.Slots()) + sideCols - 1) / sideCols
		sideW = sideCols*step - slotGap*s + 2*padding*s
		sideH = titleH*s + siderows*step - slotGap*s + padding*s
	}
	total := w
	if sideW > 0 {
		total += sideW + sideGap*s
	}
	x0 := (v.W - total) / 2
	y0 := max(barMargin, (v.H-SlotSize*s-barMargin-h)/2)
	r.panel = Rect{x0 + total - w, y0, w, h}
	if sideW > 0 {
		r.side = Rect{x0, y0, sideW, sideH}
		for i := range p.Side.Slots() {
			r.chest = append(r.chest, Rect{r.side.X + padding*s + (i%sideCols)*step, r.side.Y + titleH*s + (i/sideCols)*step, SlotSize * s, SlotSize * s})
		}
	}
	if p.Crafting() {
		for i := range rows {
			r.rows = append(r.rows, Rect{r.panel.X + (w-gridW)/2, r.panel.Y + titleH*s + i*step, gridW, SlotSize * s})
		}
	} else {
		for i := range cells {
			r.grid = append(r.grid, Rect{r.panel.X + (w-gridW)/2 + (i%p.Cols)*step, r.panel.Y + titleH*s + (i/p.Cols)*step, SlotSize * s, SlotSize * s})
		}
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
	for i, r := range p.rects.chest {
		if r.In(at) {
			return Spot{Chest, i}, true
		}
	}
	return Spot{}, false
}

// rectOf is where slot s lies.
func (p *Pack) rectOf(s Spot) Rect {
	list := p.rects.grid
	switch s.Area {
	case Bar:
		list = p.rects.bar
	case Chest:
		list = p.rects.chest
	}
	if s.Index < 0 || s.Index >= len(list) {
		return Rect{}
	}
	return list[s.Index]
}

// Where is the middle of slot s, in pixels, as the pack was last placed.
func (p *Pack) Where(s Spot) Point {
	r := p.rectOf(s)
	return Point{r.X + r.W/2, r.Y + r.H/2}
}
