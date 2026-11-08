package dialogue

// Handle does what in says. Enter, space or a press on the box shows the
// page whole if it is being typed, and goes to the next page if it is
// shown; on the last page, with the words said, Enter chooses the button
// chosen and a press on a button chooses it. The arrows move among the
// buttons, and the mouse over a button chooses it. Escape leaves.
func (b *Box) Handle(in Input) Result {
	moved := in.Mouse != b.at
	b.at = in.Mouse
	over, onButton := b.buttonAt(in.Mouse)
	switch {
	case in.Cancel:
		return Left
	case in.Press && onButton && b.Ready():
		b.Choice = over
		return Chosen
	case in.Select && b.Ready() && len(b.Choices) > 0:
		return Chosen
	case in.Select || in.Press && b.frame.In(in.Mouse):
		return b.advance()
	case in.Dy != 0 && len(b.Choices) > 0:
		b.Choice = (b.Choice + in.Dy + len(b.Choices)) % len(b.Choices)
		b.keep()
		return Moved
	case onButton && moved && over != b.Choice:
		b.Choice = over
		return Moved
	}
	return Nothing
}

// advance shows the rest of the page, or begins the next one.
func (b *Box) advance() Result {
	switch {
	case !b.Done():
		b.shown = float64(b.letters())
		return Skipped
	case !b.Last():
		b.page++
		b.shown = 0
		return Advanced
	}
	return Nothing
}

// buttonAt is the button under the pixel at.
func (b *Box) buttonAt(at Point) (int, bool) {
	for i, r := range b.buttons {
		if r.In(at) {
			return b.first + i, true
		}
	}
	return 0, false
}

// keep scrolls the buttons to keep the one chosen in sight.
func (b *Box) keep() {
	n := len(b.buttons)
	if n == 0 {
		return
	}
	if b.Choice < b.first {
		b.first = b.Choice
	}
	if b.Choice >= b.first+n {
		b.first = b.Choice - n + 1
	}
}
