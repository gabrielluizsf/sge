package kit

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/human"
)

func TestLookTag(t *testing.T) {
	v := Villager{ID: "ana", Look: human.Look{Eyes: eyes("brown")}}
	if v.LookTag() != "" {
		t.Errorf("unchanged eyes tag %q", v.LookTag())
	}
	v.Look.Eyes = eyes("blue")
	if !strings.Contains(v.LookTag(), "blue") {
		t.Errorf("changed eyes tag %q", v.LookTag())
	}
}

func TestVillagersEyesAreInThePack(t *testing.T) {
	have := map[string]bool{}
	zr, err := zip.OpenReader(filepath.Join(os.Getenv("HOME"), ".cache/sge-assets/makehuman/packs/makehuman_system_assets_cc0.zip"))
	if err == nil {
		defer zr.Close()
		for _, f := range zr.File {
			have[f.Name] = true
		}
	}
	colours := map[string]bool{}
	for _, v := range Villagers {
		colours[v.Look.Eyes] = true
		if len(have) > 0 && !have[v.Look.Eyes] {
			t.Errorf("%s: eyes %s are not in the pack", v.ID, v.Look.Eyes)
		}
		if !strings.HasPrefix(v.Look.Eyes, "eyes/materials/") || !strings.HasSuffix(v.Look.Eyes, ".mhmat") {
			t.Errorf("%s: odd eyes %q", v.ID, v.Look.Eyes)
		}
	}
	if len(colours) < 5 {
		t.Errorf("only %d eye colours in the cast", len(colours))
	}
}
