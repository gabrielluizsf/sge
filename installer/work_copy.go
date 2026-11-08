package installer

import "fmt"

// copies are the jobs copying list, a batch of about batchBytes each.
func (w *work) copies(list []entry) []job {
	var jobs []job
	for start := 0; start < len(list); {
		end, size := start, int64(0)
		for end < len(list) && (end == start || size+list[end].size <= batchBytes) {
			size += list[end].size
			end++
		}
		batch := list[start:end]
		name := fmt.Sprintf("copy %d/%d", end, len(list))
		jobs = append(jobs, job{name, func() error { return w.copyAll(batch) }})
		start = end
	}
	return jobs
}

func (w *work) copyAll(batch []entry) error {
	to := placed{root: w.root, dir: w.l.Files, exe: w.s.App.Exe, left: &w.left}
	if to.dir == "" {
		to.dir = "."
	}
	for _, e := range batch {
		if err := copyFile(w.s.Files, to, e); err != nil {
			return err
		}
	}
	return nil
}
