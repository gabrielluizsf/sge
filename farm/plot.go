package farm

import "github.com/gabrielluizsf/sge/town"

// plot is a cell of soil: what grows on it, how far, and whether it has
// been watered today.
type plot struct {
	col, row int
	crop     int // an index of crops, or -1
	stage    int
	watered  bool
	days     int        // days it has been watered
	dead     bool       // it went a day dry, and withered
	town     *town.Town // whose garden it is, nil for the player's field
	// owner is the villager who sowed this soil of the farm and keeps the
	// crop until it is picked, nil while it is the player's to use.
	owner *villager
}

// thirsty reports whether watering would keep the plot's crop alive and
// growing: it is growing, not ripe or dead, and dry today.
func (p *plot) thirsty() bool {
	return p.crop >= 0 && !p.dead && p.stage < ripe && !p.watered
}

// ripened reports whether the plot's crop is ripe, waiting to be picked.
func (p *plot) ripened() bool { return p.crop >= 0 && !p.dead && p.stage == ripe }

// clear leaves the plot bare: nothing grows on it.
func (p *plot) clear() {
	p.crop, p.stage, p.watered, p.days, p.dead = -1, 0, false, 0, false
	p.owner = nil
}
