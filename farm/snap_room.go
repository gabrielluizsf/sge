package farm

import "github.com/gabrielluizsf/sge/vector"

// roomPlace is where in a house the player is: the cell before its door, and
// the spot on its floor.
type roomPlace struct {
	Door cellState `json:"door"`
	At   point     `json:"at"`
}

// takeRoom is the house the player is inside, or nil out of doors.
func takeRoom(f *farm) *roomPlace {
	r := f.room
	if r == nil {
		return nil
	}
	p := f.player.Transform.Position
	return &roomPlace{Door: cellState{Col: r.door.Col, Row: r.door.Row}, At: point{X: p.X, Y: p.Y}}
}

// checkRoom refuses a room with its numbers out of range.
func checkRoom(s *snapshot) error {
	if s.Room == nil {
		return nil
	}
	return firstError(checkCell("room door", s.Room.Door.Col, s.Room.Door.Row), checkPoint("room spot", s.Room.At))
}

// applyRoom puts the player back in the house they were saved in, on the floor,
// and leaves them on the doorstep if no such house is there any more.
func applyRoom(f *farm, s *snapshot) {
	if s.Room == nil || f.room != nil {
		return
	}
	d, ok := f.doorAt(s.Room.Door.Col, s.Room.Door.Row)
	if !ok {
		return
	}
	f.goInside(d)
	at := f.onFloor(s.Room.At)
	f.player.Transform.Position.X, f.player.Transform.Position.Y = at.X, at.Y
}

// onFloor is p, or the nearest spot to it inside a room's walls.
func (f *farm) onFloor(p point) vector.Vector3 {
	low := f.ground.CellCenter(1, 1)
	high := f.ground.CellCenter(roomCols-2, roomRows-2)
	return vector.Vector3{X: min(max(p.X, low.X), high.X), Y: min(max(p.Y, low.Y), high.Y)}
}
