package save

import (
	"os"
	"path/filepath"
)

// Store is the folder of one game's saves, and the version its files have.
type Store struct {
	dir     string
	version int
}

// Open is the store of the game titled title, in the user's config
// directory, in a folder named by Id(title). The folder is made when first
// written to.
func Open(title string, version int) (*Store, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return In(filepath.Join(base, Id(title)), version), nil
}

// In is the store in folder dir, for tests and for games that choose their
// own place.
func In(dir string, version int) *Store { return &Store{dir: dir, version: version} }

// Dir is the store's folder.
func (s *Store) Dir() string { return s.dir }

// path is where the file called name is, once name is checked.
func (s *Store) path(name string) (string, error) {
	if err := checkName(name); err != nil {
		return "", err
	}
	return filepath.Join(s.dir, name+".json"), nil
}

// Has reports whether a file called name has been saved.
func (s *Store) Has(name string) bool {
	p, err := s.path(name)
	if err != nil {
		return false
	}
	st, err := os.Stat(p)
	return err == nil && st.Mode().IsRegular()
}
