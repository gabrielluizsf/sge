package bakes

import (
	"encoding/json"
	"image/png"
	"os"
	"path/filepath"
	"sync"

	"github.com/gabrielluizsf/sge/bake"
)

// Sprite is where a baked picture stands: the pixel its model's foot is.
type Sprite struct {
	W, H   int
	PivotX int
	PivotY int
}

// Shelf is a folder of baked PNGs and the sprites.json that says where
// each stands; it is safe to put to from many goroutines. Tasks finished
// are kept in done.json beside it.
type Shelf struct {
	Dir string

	mu      sync.Mutex
	sprites map[string]Sprite
	done    map[string]bool
}

// Open opens the shelf in dir, made if it is not there, with what it
// already holds.
func Open(dir string) (*Shelf, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Shelf{Dir: dir, sprites: map[string]Sprite{}, done: map[string]bool{}}
	readJSON(filepath.Join(dir, "sprites.json"), &s.sprites)
	var names []string
	readJSON(filepath.Join(dir, "done.json"), &names)
	for _, n := range names {
		s.done[n] = true
	}
	return s, nil
}

// readJSON reads the file at name into v, if it is there and sound.
func readJSON(name string, v any) {
	if data, err := os.ReadFile(name); err == nil {
		json.Unmarshal(data, v)
	}
}

// Put writes sprite sp as name's PNG and records where it stands.
func (s *Shelf) Put(name string, sp bake.Sprite) error {
	f, err := os.Create(filepath.Join(s.Dir, name+".png"))
	if err != nil {
		return err
	}
	if err := png.Encode(f, sp.Image); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	b := sp.Image.Bounds()
	s.mu.Lock()
	s.sprites[name] = Sprite{W: b.Dx(), H: b.Dy(), PivotX: sp.Pivot.X, PivotY: sp.Pivot.Y}
	s.mu.Unlock()
	return nil
}

// Sprite is where the sprite name stands, if it is baked.
func (s *Shelf) Sprite(name string) (Sprite, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sp, ok := s.sprites[name]
	return sp, ok
}

// Has reports whether the sprite name is baked.
func (s *Shelf) Has(name string) bool {
	_, ok := s.Sprite(name)
	return ok
}

// Finished reports whether the task name ran to its end.
func (s *Shelf) Finished(task string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.done[task]
}

// Mark records the task name as run to its end.
func (s *Shelf) Mark(task string) {
	s.mu.Lock()
	s.done[task] = true
	s.mu.Unlock()
}

// Save writes sprites.json and done.json; encoding/json sorts the keys.
func (s *Shelf) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	names := make([]string, 0, len(s.done))
	for n := range s.done {
		names = append(names, n)
	}
	for file, v := range map[string]any{"sprites.json": s.sprites, "done.json": names} {
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(s.Dir, file), append(data, '\n'), 0o644); err != nil {
			return err
		}
	}
	return nil
}
