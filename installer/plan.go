package installer

import (
	"errors"
	"os"
)

// job is one step of installing: a name for the window, and the work.
type job struct {
	name string
	run  func() error
}

// batchBytes is about how much one copying job copies, so the window's
// bar moves and is drawn between them.
const batchBytes = 8 << 20

// plan is the installing as jobs, run in order: everything is checked
// before the first runs, and the files are copied into a new folder that
// takes the old one's place only once all of them are in.
func plan(s Setup) ([]job, error) {
	if s.Files == nil {
		return nil, errors.New("installer: no files")
	}
	if s.Cache != "" {
		if err := cacheOK(s); err != nil {
			return nil, err
		}
	}
	icon, err := iconData(s)
	if err != nil {
		return nil, err
	}
	l, err := layoutOf(s.App, s.Place, icon)
	if err != nil {
		return nil, err
	}
	list, err := entries(s.Files, s.Limits.orDefault(), l.own())
	if err != nil {
		return nil, err
	}
	if !hasFile(list, s.App.Exe) {
		return nil, errors.New("installer: the files have no " + s.App.Exe)
	}
	w := &work{s: s, l: l, left: s.Limits.orDefault().Bytes}
	jobs := []job{{"prepare", w.prepare}}
	jobs = append(jobs, w.copies(list)...)
	jobs = append(jobs, job{"finish", w.finish}, job{"register", w.register})
	for _, st := range s.Steps {
		jobs = append(jobs, job{st.Name, func() error { return st.Run(l.join(l.Dir, l.Files)) }})
	}
	return append(jobs, job{"cache", w.markCache}), nil
}

func hasFile(list []entry, name string) bool {
	for _, e := range list {
		if e.name == name {
			return true
		}
	}
	return false
}

// own are the paths in the files the installer writes itself.
func (l layout) own() []string {
	var own []string
	for _, m := range l.Inside {
		if rel, ok := cutDir(m.Path, l.Files); ok {
			own = append(own, rel)
		}
	}
	if l.Uninstaller != "" && l.Files == "" {
		own = append(own, l.Uninstaller)
	}
	return own
}

// cutDir is p relative to dir, when it is in it.
func cutDir(p, dir string) (string, bool) {
	if dir == "" {
		return p, true
	}
	if len(p) > len(dir) && p[:len(dir)+1] == dir+"/" {
		return p[len(dir)+1:], true
	}
	return "", false
}

// work is an installing under way.
type work struct {
	s    Setup
	l    layout
	root *os.Root
	left int64
}

// stage is the new folder being filled, and old the one it replaces.
func (l layout) stage() string { return l.Dir + ".installing" }
func (l layout) old() string   { return l.Dir + ".old" }
