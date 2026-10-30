package farm

// notice is news shown in the middle of the screen for a few seconds.
type notice struct {
	text string
	left float64
}

const noticeFor = 4

// notify shows text in the middle of the screen.
func (f *farm) notify(text string) {
	f.notices = append(f.notices, &notice{text: text, left: noticeFor})
}

// fade lets the notices go as their time runs out.
func (f *farm) fade(dt float64) {
	kept := f.notices[:0]
	for _, n := range f.notices {
		if n.left -= dt; n.left > 0 {
			kept = append(kept, n)
		}
	}
	f.notices = kept
}
