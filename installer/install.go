package installer

// Install installs s, each job in turn; on an error, the new folder is
// taken away and an app already installed is left as it was.
func Install(s Setup) error {
	w, err := plan(s)
	if err != nil {
		return err
	}
	return w.run(nil)
}

// run runs the jobs in order, telling each one's start to began (when
// set); on an error, it takes away the new folder half filled.
func (w *work) run(began func(i int)) error {
	for i, j := range w.jobs {
		if began != nil {
			began(i)
		}
		if err := j.run(); err != nil {
			w.abandon()
			return err
		}
	}
	return nil
}

// abandon lets go of the new folder of an installing that failed, as
// Windows removes no folder still open, and takes it away.
func (w *work) abandon() {
	if w.root != nil {
		w.root.Close()
		w.root = nil
	}
	clear(w.l.stage(), w.l, w.s.App)
}
