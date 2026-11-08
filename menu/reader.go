package menu

// Reader is a popup of a long text, such as a licence, scrolled by the
// arrows, the page keys and the wheel; Enter or Escape closes it.
type Reader struct {
	Title string
	Text  string
	// Top is the first line of the wrapped text in sight.
	Top int
	// Scale is how many pixels a pixel of the game's art takes.
	Scale int

	wrapped []string
	wrapW   int // the width the text was wrapped to, in pixels
	shown   int // how many lines fit
	words   map[string]int
}

// Open shows text under title from its start.
func (r *Reader) Open(title, text string) {
	r.Title, r.Text, r.Top, r.wrapW = title, text, 0, 0
}

// Handle does what in says: the arrows and the wheel scroll a line or a
// notch, the page keys a page; Enter or Escape closes it.
func (r *Reader) Handle(in Input) Result {
	switch {
	case in.Cancel || in.Select:
		return Cancelled
	case in.Dy != 0 || in.Page != 0 || in.Wheel != 0:
		r.Top = r.clamp(r.Top + in.Dy + in.Wheel*3 + in.Page*max(1, r.shown-1))
		return Moved
	}
	return Nothing
}

// clamp keeps top between the first line and the last page.
func (r *Reader) clamp(top int) int {
	return max(0, min(top, len(r.wrapped)-r.shown))
}

// wrap breaks the text into lines no wider than w pixels, once for a width.
func (r *Reader) wrap(w int, wide func(string) int) {
	if r.wrapW != w || r.wrapped == nil {
		r.wrapW, r.wrapped = w, Wrap(r.Text, w, wide)
	}
}
