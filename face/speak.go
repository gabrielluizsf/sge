package face

// hold is the least time, in seconds, a mouth keeps a shape: about the
// pace sounds are made at, whatever the pace the letters are shown at.
const hold = 0.08

// Speaker is a mouth saying letters as they are shown.
type Speaker struct {
	shape Shape
	held  float64
}

// Update lets dt seconds go by, the letter r being said, if saying: the
// mouth takes r's shape once it has held its own long enough, or at once
// from shut, and shuts when nothing is said.
func (s *Speaker) Update(dt float64, r rune, saying bool) {
	s.held += dt
	want := Shut
	if saying {
		want = ShapeOf(r)
	}
	if want != s.shape && (s.held >= hold || s.shape == Shut || !saying) {
		s.shape, s.held = want, 0
	}
}

// Shape is the mouth's shape now.
func (s *Speaker) Shape() Shape { return s.shape }

// ShapeOf is the mouth's shape saying r, a letter of Portuguese or
// English: shut for the lips' letters and anything not a letter.
func ShapeOf(r rune) Shape {
	r = base(lower(r))
	switch r {
	case 'a', 'e', 'i':
		return Open
	case 'o', 'u', 'w':
		return Round
	case 'm', 'b', 'p':
		return Shut
	}
	if r >= 'a' && r <= 'z' || r == 'ç' || r == 'ñ' {
		return Ajar
	}
	return Shut
}

// lower is r in lower case, for the letters of Portuguese and English.
func lower(r rune) rune {
	switch {
	case r >= 'A' && r <= 'Z':
		return r + 'a' - 'A'
	case r >= 'À' && r <= 'Þ' && r != '×':
		return r + 'à' - 'À'
	}
	return r
}

// base is r without its accent, for the accented vowels of Portuguese.
func base(r rune) rune {
	switch r {
	case 'á', 'à', 'â', 'ã':
		return 'a'
	case 'é', 'ê':
		return 'e'
	case 'í':
		return 'i'
	case 'ó', 'ô', 'õ':
		return 'o'
	case 'ú', 'ü':
		return 'u'
	}
	return r
}
