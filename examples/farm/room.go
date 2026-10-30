package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// room is the inside of a house being visited: its floor and what stands
// on it, the house's door outside, and whose it is.
type room struct {
	ground, things *tilemap.Map
	door           town.Cell
	faceOut        string
	own            bool
	kind           int
}

// The rooms' size, in cells: the back walls along the first row and
// column, the floor inside, and the doormat in the middle of the last
// row, the way out.
const (
	roomCols, roomRows = 8, 7
	matCol             = 3
)

// houseDoor is a house's door seen from outside: the cell before it,
// which way faces out, its kind, and a seed for its furniture.
type houseDoor struct {
	at      town.Cell
	faceOut string
	kind    int
	seed    uint32
	own     bool
}

// newRoom lays out the inside of the house at door, its furniture placed
// as its seed falls.
func newRoom(d houseDoor) *room {
	g, t := roomLayout(d.seed)
	r := &room{door: d.at, faceOut: d.faceOut, own: d.own, kind: d.kind}
	legend := map[rune]int{'.': tileID["floor"], 'm': tileID["mat"], 'r': tileID[fmt.Sprint("rug", d.seed%4)]}
	r.ground = parseRoom("room floor", g, legend)
	things := map[rune]int{
		'C': tileID[fmt.Sprint("wall_nw", d.kind)], 'n': tileID[fmt.Sprint("wall_n", d.kind)],
		'W': tileID[fmt.Sprint("wall_nwin", d.kind)], 'w': tileID[fmt.Sprint("wall_w", d.kind)],
		'x': tileID["block"], 'b': tileID["bed"], 'B': tileID["bed_blue"], 't': tileID["table"], 'c': tileID["chair"],
		's': tileID["shelf"], 'o': tileID["stove"], 'p': tileID["plant"], 'h': tileID["chest"],
	}
	r.things = parseRoom("room things", t, things)
	r.ground.Z, r.things.Z = zGround, zThings
	r.things.Sorted = true
	return r
}

func parseRoom(name, rows string, legend map[rune]int) *tilemap.Map {
	m, err := tilemap.Parse(name, roomTiles, rows, legend)
	if err != nil {
		panic(err)
	}
	return m
}

// roomTiles is the tile set the rooms are made of: the farm's.
var roomTiles *tilemap.TileSet
