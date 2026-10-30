package gots

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// shared is a set of packages transpiled once for every test that runs
// them: where it went, and whether it could be.
type shared struct {
	once sync.Once
	dir  string
	err  error
}

var (
	sharedMu   sync.Mutex
	sharedSets = map[string]*shared{}
)

// transpiled is the directory the set name, transpiled with opts, is in,
// transpiling it the first time a test asks; it is removed when the tests
// end. Tests that share a set must not change what is in it.
func transpiled(t *testing.T, name string, opts func() Options) string {
	t.Helper()
	sharedMu.Lock()
	s := sharedSets[name]
	if s == nil {
		s = &shared{}
		sharedSets[name] = s
	}
	sharedMu.Unlock()
	s.once.Do(func() {
		s.dir, s.err = os.MkdirTemp("", "gots-"+name)
		if s.err != nil {
			return
		}
		o := opts()
		o.Dir, o.Out = ".", s.dir
		_, s.err = Transpile(o)
	})
	if s.err != nil {
		t.Fatalf("transpile %s: %v", name, s.err)
	}
	return s.dir
}

// engine is sge's packages transpiled with their tests, antui's canvas
// and an index.
func engine(t *testing.T) string {
	return transpiled(t, "engine", func() Options {
		// The engine first: the first package's module is laid out at the
		// top of the directory, the others' under _gots/mod.
		var patterns []string
		for _, p := range enginePackages {
			patterns = append(patterns, "../"+p)
		}
		patterns = append(patterns, "github.com/gabrielluizsf/antui/canvas")
		return Options{Patterns: patterns, Tests: true, Index: true}
	})
}

// programs is the programs of testdata/exits and testdata/nosockets,
// transpiled together, each at gots/testdata/<set>/<name>/mod.ts.
func programs(t *testing.T) string {
	return transpiled(t, "programs", func() Options {
		dirs, _ := filepath.Glob("testdata/exits/*")
		dirs = append(dirs, "testdata/nosockets/tlsconfig")
		var patterns []string
		for _, d := range dirs {
			patterns = append(patterns, "./"+d)
		}
		return Options{Patterns: patterns}
	})
}

func TestMain(m *testing.M) {
	code := m.Run()
	for _, s := range sharedSets {
		if s.dir != "" {
			os.RemoveAll(s.dir)
		}
	}
	os.Exit(code)
}
