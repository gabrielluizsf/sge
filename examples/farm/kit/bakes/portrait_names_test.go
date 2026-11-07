package bakes

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/kit"
)

func taskNames(b *Baker) map[string]bool {
	out := map[string]bool{}
	for _, t := range append(b.portraitTasks(), b.portraitFaceTasks()...) {
		out[t.Name] = true
	}
	return out
}

// A villager whose eyes changed is baked again under a new name, one whose
// look did not keeps the name the shelf already knows.
func TestPortraitTasksNoticeANewLook(t *testing.T) {
	b := &Baker{}
	before := taskNames(b)
	if !before["portrait ana adult"] || !before["portrait face ana adult v2"] {
		t.Fatalf("unchanged look lost its name: %v", before)
	}
	old := kit.Villagers[1].Look.Eyes
	kit.Villagers[1].Look.Eyes = "eyes/materials/blue.mhmat"
	defer func() { kit.Villagers[1].Look.Eyes = old }()
	after := taskNames(b)
	if after["portrait ana adult"] || !after["portrait ana adult eyes blue"] {
		t.Errorf("a new look kept its old name")
	}
	if !after["portrait beto adult"] {
		t.Error("another villager's name changed")
	}
	for n := range after {
		if strings.Contains(n, "  ") {
			t.Errorf("odd name %q", n)
		}
	}
}
