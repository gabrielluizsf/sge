// Command pinhumans pins what the realistic farm's people are made of,
// into humans.json: MakeHuman's data files, from a copy of its repository
// at a commit, each to be fetched from GitHub at that commit, and its CC0
// asset packs, from copies of them. It keeps each copy in the asset
// cache, so they are not fetched again; every later fetch is checked
// against its pin.
//
//	go run ./cmd/pinhumans -data ~/makehuman/makehuman/data -commit a8bc2d5… -packs ~/packs
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/fetch"
)

// packs are the asset packs the villagers wear.
var packs = []string{"makehuman_system_assets", "shirts01", "pants01", "shoes01", "skirts01", "dress01"}

// data are the data files a kit reads: the base mesh, skeleton and
// weights, and the macro targets.
var data = []string{"3dobjs/base.obj", "rigs/default.mhskel", "rigs/default_weights.mhw", "targets/macrodetails"}

func main() {
	dir := flag.String("data", "", "MakeHuman's data folder, of its repository at the commit")
	commit := flag.String("commit", "", "the repository's commit")
	packDir := flag.String("packs", "", "the folder of the asset packs' zips")
	out := flag.String("out", "humans.json", "the pins to write")
	flag.Parse()
	if len(*commit) != 40 || strings.Trim(*commit, "0123456789abcdef") != "" {
		fail(fmt.Errorf("-commit %q is not a commit", *commit))
	}
	fe := kit.Fetcher()
	h := kit.Humans{Source: "https://github.com/makehumancommunity/makehuman", License: "CC0-1.0"}
	for _, name := range files(*dir) {
		url := "https://raw.githubusercontent.com/makehumancommunity/makehuman/" + *commit + "/makehuman/data/" + name
		h.Data = append(h.Data, pin(fe, filepath.Join(*dir, filepath.FromSlash(name)), fetch.File{URL: url, Path: path.Join(kit.HumansData, name)}))
	}
	for _, p := range packs {
		url := "https://files.makehumancommunity.org/asset_packs/" + p + "/" + p + "_cc0.zip"
		h.Packs = append(h.Packs, pin(fe, filepath.Join(*packDir, p+"_cc0.zip"), fetch.File{URL: url, Path: "makehuman/packs/" + p + "_cc0.zip"}))
	}
	b, _ := json.MarshalIndent(h, "", "  ")
	if err := os.WriteFile(*out, append(b, '\n'), 0o644); err != nil {
		fail(err)
	}
	fmt.Printf("%d data files and %d packs pinned\n", len(h.Data), len(h.Packs))
}

// files are the data files under dir to pin, their names from it.
func files(dir string) []string {
	var out []string
	for _, d := range data {
		err := fs.WalkDir(os.DirFS(dir), d, func(name string, e fs.DirEntry, err error) error {
			if err == nil && e.Type().IsRegular() && (d != "targets/macrodetails" || strings.HasSuffix(name, ".target")) {
				out = append(out, name)
			}
			return err
		})
		if err != nil {
			fail(err)
		}
	}
	return out
}

// pin is f pinned to the copy at local's bytes, kept in the cache.
func pin(fe fetch.Fetcher, local string, f fetch.File) fetch.File {
	b, err := os.ReadFile(local)
	if err != nil {
		fail(err)
	}
	sum := sha256.Sum256(b)
	f.SHA256 = hex.EncodeToString(sum[:])
	if err := fe.Keep(f, b); err != nil {
		fail(err)
	}
	return f
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "pinhumans:", err)
	os.Exit(1)
}
