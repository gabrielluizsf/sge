package farm

import "fmt"

// openFight sets a fight up on s, to look at it: the player at the smith of
// the east town with an axe in hand, hungry and a little hurt; the smith
// struck and bleeding, going after them with an axe, whose blow lands on
// the player a moment before the last frame. It is an error if there is no one.
func (s *Session) openFight(d Demo) error {
	f := s.f
	t := f.towns.Of(1, 0)
	f.settle(t)
	for _, v := range f.people {
		if v.town != t || !v.smith() {
			continue
		}
		f.player.Transform.Position = v.pos
		if _, err := s.Play(Demo{Frames: 60, Width: d.Width, Height: d.Height}); err != nil {
			return err
		}
		if v.obj == nil {
			break
		}
		f.hands.ctl = still{}
		f.axe = 1
		f.stock(itemAxe)
		f.hold(f.pack.Slot(itemAxe))
		f.vit.Health, f.vit.Food = 62, 20
		f.player.Transform.Position = v.obj.Transform.Position
		f.player.Transform.Position.X -= 0.9 * tile
		f.hands.facing.X, f.hands.facing.Y = 1, 0
		f.wound(v, 15, f.player.Transform.Position)
		v.swing = float64(d.Frames-30) / 60
		return nil
	}
	return fmt.Errorf("farm: no one to fight")
}
