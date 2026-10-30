// Command pin adds Poly Haven models to the realistic farm's manifest:
// for each id, it asks Poly Haven's API for the model's glTF files at a
// resolution, fetches each once, checks it against the MD5 the API gives,
// pins its SHA-256, so every later fetch is checked against that, and
// keeps it in the asset cache.
// Every asset on Poly Haven is CC0. A texture is pinned as the ground's
// unless -kind material says it is a material things are built of; a
// texture's span, the metres one repeat covers, is pinned with it.
//
//	go run ./cmd/pin -biome campo -res 1k jacaranda_tree boulder_01
//	go run ./cmd/pin -biome building -kind material brick_wall_006
//	go run ./cmd/pin -biome town -whole wine_barrel_01
package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/gabrielluizsf/sge/examples/realfarm"
	"github.com/gabrielluizsf/sge/fetch"
)

func main() {
	file := flag.String("manifest", "assets.json", "the manifest to add to")
	res := flag.String("res", "1k", "the textures' resolution")
	biome := flag.String("biome", "", "the biome the assets are of")
	kind := flag.String("kind", "", `"material" pins textures as materials to build with, not ground`)
	whole := flag.Bool("whole", false, "the models' roots are the parts of one thing, not variants")
	flag.Parse()
	if *kind != "" && *kind != "material" {
		fail(fmt.Errorf(`kind %q: only "material" may be asked for`, *kind))
	}
	m, err := realfarm.ReadManifest(*file)
	if err != nil && !os.IsNotExist(err) {
		fail(err)
	}
	for _, id := range flag.Args() {
		a, err := pin(request{id: id, res: *res, biome: *biome, kind: *kind})
		a.Whole = *whole
		if err != nil {
			fail(fmt.Errorf("%s: %w", id, err))
		}
		m.Assets = slices.DeleteFunc(m.Assets, func(x realfarm.Asset) bool { return x.ID == id })
		m.Assets = append(m.Assets, a)
		fmt.Printf("%s: %d files pinned\n", id, len(a.Files))
	}
	sort.Slice(m.Assets, func(i, j int) bool { return m.Assets[i].ID < m.Assets[j].ID })
	out, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(*file, append(out, '\n'), 0o644); err != nil {
		fail(err)
	}
}

// apiFile is a file as Poly Haven's API describes it.
type apiFile struct {
	URL     string             `json:"url"`
	MD5     string             `json:"md5"`
	Size    int64              `json:"size"`
	Include map[string]apiFile `json:"include"`
}

// request is an asset to pin: its id, the resolution of its textures,
// its biome, and the kind to pin a texture as, "" for the ground.
type request struct {
	id, res, biome, kind string
}

// pin is the asset r asks for, its files fetched and pinned.
func pin(r request) (realfarm.Asset, error) {
	id, res, biome := r.id, r.res, r.biome
	var info struct {
		Name    string            `json:"name"`
		Authors map[string]string `json:"authors"`
		Type    int               `json:"type"` // 1 a texture, 2 a model
		// Dimensions are a texture's width and height, in millimetres.
		Dimensions []float64 `json:"dimensions"`
	}
	if err := getJSON("https://api.polyhaven.com/info/"+id, &info); err != nil {
		return realfarm.Asset{}, err
	}
	var files struct {
		GLTF map[string]map[string]apiFile `json:"gltf"`
	}
	if err := getJSON("https://api.polyhaven.com/files/"+id, &files); err != nil {
		return realfarm.Asset{}, err
	}
	g, ok := files.GLTF[res]["gltf"]
	if !ok {
		return realfarm.Asset{}, fmt.Errorf("no glTF at %s", res)
	}
	a := realfarm.Asset{ID: id, Name: info.Name, Biome: biome, Source: "https://polyhaven.com/a/" + id, License: "CC0-1.0",
		Model: path.Join(id, path.Base(g.URL)), Kind: map[bool]string{true: "ground", false: "model"}[info.Type == 1]}
	if info.Type == 1 {
		span, err := spanOf(info.Dimensions)
		if err != nil {
			return a, err
		}
		a.Span = span
		if r.kind != "" {
			a.Kind = r.kind
		}
	}
	for who := range info.Authors {
		a.Authors = append(a.Authors, who)
	}
	sort.Strings(a.Authors)
	all := map[string]apiFile{path.Base(g.URL): g}
	for name, f := range g.Include {
		all[name] = f
	}
	names := make([]string, 0, len(all))
	for n := range all {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, name := range names {
		f := all[name]
		// The API names the files: each must stay in the asset's folder.
		if p := path.Join(id, name); !strings.HasPrefix(p, id+"/") || strings.Contains(name, "\\") {
			return a, fmt.Errorf("file %q leaves the asset's folder", name)
		}
		data, sum, err := pinned(f)
		if err != nil {
			return a, fmt.Errorf("%s: %w", name, err)
		}
		file := fetch.File{URL: f.URL, Path: path.Join(id, name), SHA256: sum}
		if err := realfarm.Fetcher().Keep(file, data); err != nil {
			return a, err
		}
		a.Files = append(a.Files, file)
	}
	return a, nil
}

// spanOf is the metres a texture of dimensions, in millimetres, spans,
// refused unless it is a sensible size: a centimetre to a hundred metres.
func spanOf(dimensions []float64) (float64, error) {
	if len(dimensions) == 0 {
		return 0, fmt.Errorf("the texture's size is not given")
	}
	span := dimensions[0] / 1000
	if !(span >= 0.01 && span <= 100) {
		return 0, fmt.Errorf("the texture spans %v metres", span)
	}
	return span, nil
}

// pinned is the file at f.URL, fetched from Poly Haven's hosts only and
// checked against the MD5 the API gave for it, and its SHA-256.
func pinned(f apiFile) ([]byte, string, error) {
	data, err := get(f.URL, max(f.Size, 1)+1)
	if err != nil {
		return nil, "", err
	}
	m := md5.Sum(data)
	if hex.EncodeToString(m[:]) != f.MD5 {
		return nil, "", fmt.Errorf("MD5 %x, not the %s Poly Haven gives", m, f.MD5)
	}
	s := sha256.Sum256(data)
	return data, hex.EncodeToString(s[:]), nil
}

func getJSON(url string, v any) error {
	data, err := get(url, 16<<20)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// get is the body at url, from Poly Haven's hosts over HTTPS only, at
// most limit bytes.
func get(url string, limit int64) ([]byte, error) {
	fe := fetch.Fetcher{Hosts: realfarm.Hosts}
	if err := fe.Allowed(url); err != nil {
		return nil, err
	}
	resp, err := (&http.Client{CheckRedirect: func(r *http.Request, _ []*http.Request) error { return fe.Allowed(r.URL.String()) }}).Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", url, resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("%s is over %d bytes", url, limit)
	}
	return data, err
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "pin:", err)
	os.Exit(1)
}
