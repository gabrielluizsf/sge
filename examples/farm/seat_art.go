package farm

// SeatArt is what an Art may also have whose person sitting is drawn with
// the chair, one picture, so the chair hides the body where it should on
// every side. While they sit, the chair's own tile is taken away.
type SeatArt interface {
	SeatBaked() bool
}

// seatBaked reports whether a draws the chair with the person on it.
func seatBaked(a Art) bool {
	s, ok := a.(SeatArt)
	return ok && s.SeatBaked()
}
