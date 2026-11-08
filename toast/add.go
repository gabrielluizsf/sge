package toast

// Add puts a new toast with text at the end of list, in its place at once,
// and makes the oldest ones leave if more than MaxShown would be shown.
func Add(list []*Toast, text string) []*Toast {
	list = append(list, &Toast{Text: text, Life: Life(text)})
	shown := 0
	for _, t := range list {
		if !t.Leaving() {
			shown++
		}
	}
	for _, t := range list {
		if shown <= MaxShown {
			break
		}
		if !t.Leaving() {
			t.Age = max(t.Age, t.Life)
			shown--
		}
	}
	place(list)
	last := list[len(list)-1]
	last.Row = last.Target
	return list
}

// place gives each toast still holding its row among those holding; one
// leaving keeps the row it has, and those after it close the gap.
func place(list []*Toast) {
	row := 0.0
	for _, t := range list {
		if t.Leaving() {
			t.Target = t.Row
			continue
		}
		t.Target = row
		row++
	}
}
