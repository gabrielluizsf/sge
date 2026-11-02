package main

import (
	"fmt"
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
)

// Tile is the tile name drawn from baked sprites where there are some,
// and the stylized art's where there are none.
func (a *realArt) Tile(name string) *canvas.Canvas {
	if cv := a.tile(name); cv != nil {
		return cv
	}
	return a.Art.Tile(name)
}

// tile is the baked picture of the tile name, nil when it is not baked or
// is not a tile that is.
func (a *realArt) tile(name string) *canvas.Canvas {
	var mask, f int
	if n, _ := fmt.Sscanf(name, "shore%d_%d", &mask, &f); n == 2 {
		return a.shoreTile(mask, f)
	}
	base, k := splitNumber(name)
	switch base {
	case "grass":
		return a.groundTile(grassGround, k)
	case "flowers":
		return a.flowersTile(k)
	case "tree":
		return a.treeTile(k)
	}
	if contains(art.CropNames, base) {
		return a.cropTile(base, k)
	}
	switch name {
	case "withered":
		return a.witheredTile()
	case "path":
		return a.groundTile(pathGround, 0)
	case "soil_dry":
		return a.groundTile(soilGround, 0)
	case "soil_wet":
		return a.wetSoilTile()
	case "block":
		return a.thingTile("rock_07")
	case "anvil":
		return a.buildingTile("anvil")
	}
	return nil
}

// splitNumber is name without the digits it ends in, and those digits as
// a number: "grass3" is "grass", 3.
func splitNumber(name string) (string, int) {
	base := strings.TrimRight(name, "0123456789")
	var k int
	fmt.Sscanf(name[len(base):], "%d", &k)
	return base, k
}
