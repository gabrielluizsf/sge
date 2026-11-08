package installer

// Install installs s, each job in turn; on an error, the new folder is
// taken away and an app already installed is left as it was.
func Install(s Setup) error {
	jobs, err := plan(s)
	if err != nil {
		return err
	}
	return runJobs(s, jobs, nil)
}

// runJobs runs jobs in order, telling each one's start to began (when
// set); on an error, it takes away the new folder half filled.
func runJobs(s Setup, jobs []job, began func(i int)) error {
	for i, j := range jobs {
		if began != nil {
			began(i)
		}
		if err := j.run(); err != nil {
			abandon(s)
			return err
		}
	}
	return nil
}

// abandon takes away the new folder of an installing that failed.
func abandon(s Setup) {
	l, err := layoutOf(s.App, s.Place, nil)
	if err == nil {
		clear(l.stage(), l, s.App)
	}
}
