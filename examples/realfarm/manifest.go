// Package realfarm is the realistic farm: the same game as the farm, its
// pictures baked from CC0 models of real things rather than drawn.
package realfarm

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/sge/fetch"
)

// Asset is a model or texture the game is made from: its source, who
// made it and under what licence, and its files, pinned.
type Asset struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Biome   string       `json:"biome"`
	Source  string       `json:"source"`
	Authors []string     `json:"authors"`
	License string       `json:"license"`
	Model   string       `json:"model"` // the model's file, of Files
	Files   []fetch.File `json:"files"`
}

// Manifest is every asset, as assets.json lists them.
type Manifest struct {
	Assets []Asset `json:"assets"`
}

// ReadManifest reads the manifest at path.
func ReadManifest(path string) (Manifest, error) {
	var m Manifest
	data, err := os.ReadFile(path)
	if err != nil {
		return m, err
	}
	return m, json.Unmarshal(data, &m)
}

// Hosts are where assets are fetched from: Poly Haven's API and files.
var Hosts = []string{"api.polyhaven.com", "dl.polyhaven.org"}

// Fetcher fetches the assets into the user's cache, sge-assets, outside
// the repository, from Poly Haven only, no file over 512 MiB.
func Fetcher() fetch.Fetcher {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return fetch.Fetcher{Cache: filepath.Join(dir, "sge-assets"), Hosts: Hosts, Max: 512 << 20}
}
