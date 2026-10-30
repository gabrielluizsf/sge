package fetch

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// ErrMismatch is what a file whose bytes are not those pinned is.
var ErrMismatch = errors.New("fetch: SHA-256 mismatch")

// Get is the path of f in the cache, fetched now unless it is there with
// the bytes pinned. A file that fails its check is not kept.
func (fe Fetcher) Get(f File) (string, error) {
	if !filepath.IsLocal(f.Path) {
		return "", fmt.Errorf("fetch: %q leaves the cache", f.Path)
	}
	if err := os.MkdirAll(fe.Cache, 0o755); err != nil {
		return "", err
	}
	root, err := os.OpenRoot(fe.Cache)
	if err != nil {
		return "", err
	}
	defer root.Close()
	full := filepath.Join(fe.Cache, f.Path)
	if sum, err := hashOf(root, f.Path); err == nil && sum == f.SHA256 {
		return full, nil
	}
	data, err := fe.download(f.URL)
	if err != nil {
		return "", err
	}
	if err := keep(root, f, data); err != nil {
		return "", err
	}
	return full, nil
}

// Keep puts data in the cache as f, if it is the bytes f pins: for a file
// fetched some other way, as a tool pinning it has.
func (fe Fetcher) Keep(f File, data []byte) error {
	if !filepath.IsLocal(f.Path) {
		return fmt.Errorf("fetch: %q leaves the cache", f.Path)
	}
	if err := os.MkdirAll(fe.Cache, 0o755); err != nil {
		return err
	}
	root, err := os.OpenRoot(fe.Cache)
	if err != nil {
		return err
	}
	defer root.Close()
	return keep(root, f, data)
}

// keep checks data against f's pin and writes it in root.
func keep(root *os.Root, f File, data []byte) error {
	if sum := sha256.Sum256(data); hex.EncodeToString(sum[:]) != f.SHA256 {
		return fmt.Errorf("%w: %s is %x", ErrMismatch, f.URL, sum)
	}
	if dir := filepath.Dir(f.Path); dir != "." {
		if err := root.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	// Written beside it, then moved in, so no half file is ever there.
	tmp := f.Path + ".part"
	if err := root.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return root.Rename(tmp, f.Path)
}

// hashOf is the hex SHA-256 of name in root.
func hashOf(root *os.Root, name string) (string, error) {
	r, err := root.Open(name)
	if err != nil {
		return "", err
	}
	defer r.Close()
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// client is the fetcher's client, refusing to be redirected anywhere it
// may not fetch from.
func (fe Fetcher) client() *http.Client {
	c := http.Client{}
	if fe.Client != nil {
		c = *fe.Client
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("fetch: too many redirects")
		}
		return fe.Allowed(req.URL.String())
	}
	return &c
}
