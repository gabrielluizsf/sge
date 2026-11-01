package main

import "github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"

// Of everything the manifest pins, this is what the farm draws: the
// grounds it walks on, the plants and rocks that stand on them, the
// houses and the things of the town.
var (
	grounds   = []string{"sparse_grass", "farm_soil", "grass_path_2", "coast_sand_01"}
	models    = []string{"dandelion_01", "rock_07"}
	buildings = []string{"stall", "anvil"}
)

// needed reports whether the farm draws what t bakes. The houses are all
// needed, and the crops, the water and the people.
func needed(t bakes.Task) bool {
	switch t.Kind {
	case "ground":
		return contains(grounds, t.ID)
	case "model":
		return contains(models, t.ID)
	case "tree":
		return t.Core
	case "building":
		return contains(buildings, t.ID) || bakes.IsHouse(t.ID)
	}
	return true
}

// contains reports whether list has s.
func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
