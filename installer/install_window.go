package installer

import "github.com/gabrielluizsf/antui"

// installWindow installs s in a window of its own, closed when done.
func installWindow(s Setup, look Window) error {
	w, err := plan(s)
	if err != nil {
		return err
	}
	title := look.Title
	if title == "" {
		title = s.App.Name
	}
	win, err := antui.Open(title, max(look.Width, 320), max(look.Height, 180))
	if err != nil {
		return err
	}
	defer win.Close()
	return show(win, look, w)
}
