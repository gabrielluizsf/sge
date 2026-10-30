// Package kit is what the realistic farm is made of: the same game as the farm, its
// pictures baked from CC0 models of real things rather than drawn.
package kit

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/sge/fetch"
)

// Asset is a model or texture the game is made from: its source, who
// made it and under what licence, and its files, pinned.
type Asset struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Biome   string   `json:"biome"`
	Source  string   `json:"source"`
	Authors []string `json:"authors"`
	License string   `json:"license"`
	// Kind is "model", a thing standing on the ground; "ground", a
	// tiling material the ground is made of; or "material", a tiling
	// material things are built of: walls, roofs, posts.
	Kind  string `json:"kind"`
	Model string `json:"model"` // the model's file, of Files
	// Span is how many metres one repeat of a tiling material covers,
	// as Poly Haven measured it; 0 for a model.
	Span float64 `json:"span,omitempty"`
	// Whole is whether a model's roots are the parts of one thing, a
	// barrel's staves, hoops and lid, baked together, rather than
	// variants laid side by side, each baked alone.
	Whole bool         `json:"whole,omitempty"`
	Files []fetch.File `json:"files"`
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

// Hosts are where assets are fetched from: Poly Haven's API and files,
// MakeHuman's repository on GitHub and its asset packs.
var Hosts = []string{"api.polyhaven.com", "dl.polyhaven.org", "raw.githubusercontent.com", "files.makehumancommunity.org"}

// Fetcher fetches the assets into the user's cache, sge-assets, outside
// the repository, from the Hosts only, no file over 512 MiB.
func Fetcher() fetch.Fetcher {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return fetch.Fetcher{Cache: filepath.Join(dir, "sge-assets"), Hosts: Hosts, Max: 512 << 20}
}
