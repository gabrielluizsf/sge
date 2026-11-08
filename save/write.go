package save

import (
	"encoding/json"
	"os"
)

// Write saves v as the file called name, whole or not at all: it goes to a
// temporary file in the same folder, is synced, then renamed over the old
// save, so a failure on the way leaves the old save as it was. The folder
// is 0700 and the file 0600.
func (s *Store) Write(name string, v any) error {
	p, err := s.path(name)
	if err != nil {
		return err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	b, err := json.Marshal(envelope{Version: s.version, Data: data})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return err
	}
	return replace(p, b)
}
