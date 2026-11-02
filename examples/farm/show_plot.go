package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/tilemap"
)

// show sets p's tiles to what it is: wet or dry soil, and its crop's
// stage on it, or a withered plant.
func (f *farm) show(p *plot) {
	soil := "soil_dry"
	if p.watered {
		soil = "soil_wet"
	}
	f.setGround(p.col, p.row, tileID[soil])
	switch {
	case p.crop < 0:
		f.setCrop(p.col, p.row, tilemap.Empty)
	case p.dead:
		f.setCrop(p.col, p.row, tileID[witheredTile])
	default:
		f.setCrop(p.col, p.row, tileID[fmt.Sprint(crops[p.crop].name, p.stage)])
	}
}

// witheredTile is the tile of a plant that went a day dry.
const witheredTile = "withered"
