package farm

import (
	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// seat is the chair the player sits on, and where they stood before, to
// stand back there.
type seat struct {
	col, row int
	from     vector.Vector3
}

// seatFacing is the way a person sits on a chair turned each of its four
// headings: facing where the chair's front looks.
var seatFacing = []string{"down", "left", "up", "right"}

// sitDown seats the player on the chair at index i of the house's pieces,
// facing the way it is turned. Time goes on; the player can not be pushed
// off, and stands only by moving.
func (f *farm) sitDown(i int) {
	p := f.home[i]
	f.seat = &seat{col: p.col, row: p.row, from: f.player.Transform.Position}
	f.player.Transform.Position = f.ground.CellCenter(p.col, p.row)
	f.face(seatFacing[p.heading%headings])
	if seatBaked(f.art) {
		f.room.things.Set(p.col, p.row, tilemap.Empty)
	}
	f.sound("ui_tick", nil)
}

// standUp stands the player back where they were before they sat.
func (f *farm) standUp() {
	if f.seat == nil {
		return
	}
	f.player.Transform.Position = f.seat.from
	f.seat = nil
	if seatBaked(f.art) && f.room != nil {
		f.relay()
	}
	f.sound("ui_tick", nil)
}

// seatKeys is what the player pressed while seated: any step, or a swing at someone, stands them up.
func (f *farm) seatKeys(win *antui.Window) {
	a := arrows(win)
	f.seatStep(a.X != 0 || a.Y != 0 || win.KeyPressed(antui.KeyE) || win.KeyPressed(antui.KeyF))
}

// seatStep stands the player up if they mean to move.
func (f *farm) seatStep(moving bool) {
	if moving {
		f.standUp()
	}
}

// chairNear is the interaction of a chair within reach: sit on it, or, when
// seated, stand.
func (f *farm) chairNear() (words string, do func(), ok bool) {
	if f.seat != nil {
		return i18n.T("Move to stand up"), f.standUp, true
	}
	if i := f.nearPiece(f.player, "chair"); i >= 0 {
		return i18n.T("E: sit"), func() { f.sitDown(i) }, true
	}
	return "", nil, false
}
