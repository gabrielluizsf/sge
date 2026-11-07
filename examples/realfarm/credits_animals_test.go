package main

import (
	"os"
	"strings"
	"testing"
)

// CC-BY asks for a link to what was used: each animal's credit, and the
// list kept beside its sprites, links the model itself, not only its author.
func TestAnimalCreditsLinkTheirModels(t *testing.T) {
	source, err := os.ReadFile("animals/SOURCE.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range animalCredits() {
		if !strings.HasPrefix(c.URL, "https://sketchfab.com/3d-models/") {
			t.Errorf("%s links %s, not its model", c.Name, c.URL)
		}
		if !strings.Contains(string(source), c.URL) {
			t.Errorf("SOURCE.txt does not link %s", c.URL)
		}
	}
}
