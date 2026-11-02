package main

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/gabrielluizsf/sge/examples/farm"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
)

//go:embed licenses/CC0-1.0.txt
var cc0License string

//go:embed fonts/OFL.txt
var serifLicense string

// Credits are everything third-party the realistic farm is made of: its
// font, MakeHuman's people, and Poly Haven's models and materials, each
// under a licence that lets it ship with the repository.
func (a *realArt) Credits() []farm.Credit {
	cc0 := func() string { return cc0License }
	list := []farm.Credit{
		{Name: "Fonte Source Serif 4", Author: "Frank Grießhammer, Adobe", URL: "https://github.com/adobe-fonts/source-serif",
			License: func() string { return serifLicense }},
		{Name: "MakeHuman (corpos, roupas e cabelos)", Author: "The MakeHuman Community", URL: humansSource(), License: cc0},
	}
	return append(list, polyHaven(cc0)...)
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
