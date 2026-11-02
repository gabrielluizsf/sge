package dialogue

// Saying is the letter last shown on the page, while it is still being
// said, for a speaker's mouth to shape; false when the page is all shown
// or not yet laid out.
func (b *Box) Saying() (rune, bool) {
	if !b.laid || b.Done() || b.page >= len(b.pages) {
		return 0, false
	}
	n := max(0, int(b.shown)-1)
	for _, l := range b.pages[b.page] {
		rs := []rune(l)
		if n < len(rs) {
			return rs[n], true
		}
		n -= len(rs)
	}
	return 0, false
}
