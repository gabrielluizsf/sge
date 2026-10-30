package loader

// Request asks for j to be made, at priority, lower sooner; asking for an
// asset made or queued already only raises its priority when lower.
func (l *Loader) Request(j Job, priority int) {
	if l.assets == nil {
		l.assets, l.groups = map[string]*asset{}, map[string]*progress{}
	}
	if _, ok := l.assets[j.Name]; ok {
		return
	}
	for i := range l.queue {
		if l.queue[i].job.Name == j.Name {
			l.queue[i].priority = min(l.queue[i].priority, priority)
			return
		}
	}
	g := l.groups[j.Group]
	if g == nil {
		g = &progress{}
		l.groups[j.Group] = g
	}
	g.asked++
	l.seq++
	l.queue = append(l.queue, queued{job: j, priority: priority, seq: l.seq})
}

// Get is the asset name, if it is made, and the error making it gave; ok
// is false while it is not made yet.
func (l *Loader) Get(name string) (value any, err error, ok bool) {
	a, found := l.assets[name]
	if !found {
		return nil, nil, false
	}
	l.clock++
	a.used = l.clock
	return a.value, a.err, true
}

// Progress is how many of group's jobs are done, of how many asked.
func (l *Loader) Progress(group string) (done, asked int) {
	if g := l.groups[group]; g != nil {
		return g.done, g.asked
	}
	return 0, 0
}

// Pending is how many jobs wait.
func (l *Loader) Pending() int { return len(l.queue) }
