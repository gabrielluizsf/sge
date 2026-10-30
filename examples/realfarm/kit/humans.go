package kit

import (
	"archive/zip"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/sge/fetch"
	"github.com/gabrielluizsf/sge/human"
)

// Humans is where the people are made from, as humans.json pins it:
// MakeHuman's data, the base mesh, skeleton and macro targets, from its
// repository at a commit, and its CC0 asset packs, skins, hair and
// clothes, each pinned.
type Humans struct {
	Source  string       `json:"source"`
	License string       `json:"license"`
	Data    []fetch.File `json:"data"`
	Packs   []fetch.File `json:"packs"`
}

// HumansData is where MakeHuman's data files are kept in the cache.
const HumansData = "makehuman/data"

// ReadHumans reads the pins at path.
func ReadHumans(path string) (Humans, error) {
	var h Humans
	data, err := os.ReadFile(path)
	if err != nil {
		return h, err
	}
	return h, json.Unmarshal(data, &h)
}

// Kit fetches every pinned file, checked, and opens them as one kit: the
// data folder through os.Root, the packs read in their zips, never
// unpacked. Close lets the packs go.
func (h Humans) Kit(fe fetch.Fetcher) (*human.Kit, func(), error) {
	for _, f := range h.Data {
		if _, err := fe.Get(f); err != nil {
			return nil, nil, err
		}
	}
	root, err := os.OpenRoot(filepath.Join(fe.Cache, filepath.FromSlash(HumansData)))
	if err != nil {
		return nil, nil, err
	}
	fss := []fs.FS{root.FS()}
	closers := []func() error{root.Close}
	closeAll := func() {
		for _, c := range closers {
			c()
		}
	}
	for _, f := range h.Packs {
		path, err := fe.Get(f)
		if err != nil {
			closeAll()
			return nil, nil, err
		}
		z, err := zip.OpenReader(path)
		if err != nil {
			closeAll()
			return nil, nil, err
		}
		fss, closers = append(fss, z), append(closers, z.Close)
	}
	kit, err := human.Load(human.Layers(fss...), human.DefaultLimits)
	if err != nil {
		closeAll()
		return nil, nil, err
	}
	return kit, closeAll, nil
}
