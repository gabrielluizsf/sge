package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
)

//go:embed licenses/CC0-1.0.txt
var cc0License string

//go:embed fonts/OFL.txt
var serifLicense string

// Credits are everything third-party the realistic farm is made of: its
// font, MakeHuman's people, Poly Haven's models and materials, and, when
// the player has the folder of their own Sketchfab models, which the farm
// reads from there and never keeps, each of those.
func (a *realArt) Credits() []farm.Credit {
	cc0 := func() string { return cc0License }
	list := []farm.Credit{
		{Name: "Fonte Source Serif 4", Author: "Frank Grießhammer, Adobe", URL: "https://github.com/adobe-fonts/source-serif",
			License: func() string { return serifLicense }},
		{Name: "MakeHuman (corpos, roupas e cabelos)", Author: "The MakeHuman Community", URL: humansSource(), License: cc0},
	}
	list = append(list, polyHaven(cc0)...)
	return append(list, sketchfab()...)
}

// humansSource is where the people's data is from, as humans.json pins it.
func humansSource() string {
	var h struct{ Source string }
	if err := json.Unmarshal(humansJSON, &h); err != nil || h.Source == "" {
		return "https://www.makehumancommunity.org"
	}
	return h.Source
}

// polyHaven are the credits of each asset of the manifest, with whoever
// made it, if the manifest says.
func polyHaven(cc0 func() string) []farm.Credit {
	var m kit.Manifest
	if err := json.Unmarshal(assetsJSON, &m); err != nil {
		return nil
	}
	var out []farm.Credit
	for _, as := range m.Assets {
		author := "Poly Haven"
		if len(as.Authors) > 0 {
			author = strings.Join(as.Authors, ", ") + " (Poly Haven)"
		}
		out = append(out, farm.Credit{Name: "Poly Haven - " + as.Name, Author: author, URL: as.Source, License: cc0})
	}
	return out
}

// sketchfabDirs are the folders, under the player's imports, of the
// Sketchfab models the farm draws from, each with the licence its author
// gave beside it.
var sketchfabDirs = []string{"farm", "farm/buildings", "farm/buildings/waterwell", "water", "animals"}

// licenceFile is the name of the licence beside a Sketchfab model.
const licenceFile = "license.txt"

// licenceMax is the most of a licence file that is read.
const licenceMax = 1 << 20

// sketchfab are the credits of the models the player has in their imports
// folder, none if they have not: the title, author and link are read from
// each licence file there, and its text is read again when it is asked
// for. Nothing of it is copied.
func sketchfab() []farm.Credit {
	var out []farm.Credit
	for _, dir := range sketchfabDirs {
		name := dir + "/" + licenceFile
		text, err := readLicence(bakes.ImportsDir(), name)
		if err != nil {
			continue
		}
		info := licenceInfo(text)
		out = append(out, farm.Credit{Name: "Sketchfab - " + info["title"], Author: info["author"], URL: info["source"],
			License: func() string {
				t, err := readLicence(bakes.ImportsDir(), name)
				if err != nil {
					return "Licença não encontrada na pasta de modelos."
				}
				return t
			}})
	}
	return out
}

// readLicence is the licence file name under dir, at most licenceMax bytes
// of it, opened from inside dir.
func readLicence(dir, name string) (string, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return "", err
	}
	defer root.Close()
	f, err := root.Open(filepath.FromSlash(name))
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf, err := io.ReadAll(io.LimitReader(f, licenceMax))
	return string(buf), err
}

// licenceInfo is the title, source and author a Sketchfab licence file
// names, from its lines "* title:<tab>...".
func licenceInfo(text string) map[string]string {
	info := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
		if k, v, ok := strings.Cut(line, ":"); ok {
			if k = strings.TrimSpace(k); k == "title" || k == "source" || k == "author" {
				info[k] = strings.TrimSpace(v)
			}
		}
	}
	return info
}
