package save

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Read loads the file called name into v, which must be a pointer. It reads
// at most limit bytes, an error if the file is larger; refuses fields v
// does not have, anything after the one JSON value, and a version other
// than the store's.
func (s *Store) Read(name string, v any, limit int64) error {
	p, err := s.path(name)
	if err != nil {
		return err
	}
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil {
		return err
	}
	if int64(len(b)) > limit {
		return fmt.Errorf("save: %s is larger than %d bytes", name, limit)
	}
	return s.decode(b, v)
}

// decode reads the envelope in b and its data into v.
func (s *Store) decode(b []byte, v any) error {
	var e envelope
	if err := strict(b, &e); err != nil {
		return err
	}
	if e.Version != s.version {
		return fmt.Errorf("save: version %d, not %d", e.Version, s.version)
	}
	return strict(e.Data, v)
}

// strict decodes b, one JSON value and nothing else, into v, refusing
// fields v does not have.
func strict(b []byte, v any) error {
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return err
	}
	if _, err := d.Token(); !errors.Is(err, io.EOF) {
		return errors.New("save: data after the end of the value")
	}
	return nil
}
