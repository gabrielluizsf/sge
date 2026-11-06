package menu

// Line is a line of a List: its text, a smaller line under it, and the
// label of a button at its end, each possibly empty. Key, when set, names
// the line for the game, so its text may change, as a language does,
// without the line losing its meaning.
type Line struct {
	Key, Text, Sub, Button string
}

// List is a titled panel of lines, one chosen, that scrolls to keep it in
// sight.
type List struct {
	Title string
	Hint  string
	Lines []Line
	// Choice is the line chosen, Top the first one in sight.
	Choice, Top int
	// Scale is how many pixels a pixel of the game's art takes; MinW the
	// panel's least width, in pixels of art.
	Scale, MinW int

	rows  []Rect // where each line in sight is
	shown int    // how many lines fit
	at    Point  // where the mouse was
}

// Handle does what in says: up and down choose, Enter or space picks the
// line chosen, Escape leaves, the wheel scrolls, and the mouse chooses the
// line it is over and picks it with a press.
func (l *List) Handle(in Input) Result {
	switch {
	case in.Cancel:
		return Cancelled
	case in.Select:
		return Chosen
	case in.Press:
		if i, ok := l.lineAt(in.Mouse); ok {
			l.Choice = i
			return Chosen
		}
	case in.Dy != 0 || in.Wheel != 0:
		return l.move(in)
	}
	moved := in.Mouse != l.at
	l.at = in.Mouse
	if i, ok := l.lineAt(in.Mouse); ok && moved && i != l.Choice {
		l.Choice = i
		return Moved
	}
	return Nothing
}

// move steps the choice by the arrows, or the top by the wheel.
func (l *List) move(in Input) Result {
	n := len(l.Lines)
	if n == 0 {
		return Nothing
	}
	if in.Dy != 0 {
		l.Choice = (l.Choice + in.Dy + n) % n
		l.keep()
		return Moved
	}
	l.Top = max(0, min(l.Top+in.Wheel, n-l.shown))
	l.Choice = max(l.Top, min(l.Choice, l.Top+l.shown-1))
	return Moved
}

// keep scrolls so that the line chosen is in sight and the list does not
// end with blank rows.
func (l *List) keep() {
	if l.shown <= 0 {
		return
	}
	if l.Choice < l.Top {
		l.Top = l.Choice
	}
	if l.Choice >= l.Top+l.shown {
		l.Top = l.Choice - l.shown + 1
	}
	l.Top = max(0, min(l.Top, len(l.Lines)-l.shown))
}

// lineAt is the line shown under the pixel at.
func (l *List) lineAt(at Point) (int, bool) {
	for i, r := range l.rows {
		if r.In(at) {
			return l.Top + i, true
		}
	}
	return 0, false
}
