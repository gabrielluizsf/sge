package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/tilemap"
)

// show sets p's tiles to what it is: wet or dry soil, and its crop's
// stage on it.
func (f *farm) show(p *plot) {
	soil := "soil_dry"
	if p.watered {
		soil = "soil_wet"
	}
	f.setGround(p.col, p.row, tileID[soil])
	if p.crop < 0 {
		f.setCrop(p.col, p.row, tilemap.Empty)
		return
	}
	f.setCrop(p.col, p.row, tileID[fmt.Sprint(crops[p.crop].name, p.stage)])
}
