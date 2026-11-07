package farm

import (
	"fmt"
	"math"
	"os"
	"sort"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/pathfind"
)

// roomyArt is the stylized art with the houses and the stall as large as
// the realistic look's, so the towns are laid out with the wider gaps.
type roomyArt struct{ Art }

func (roomyArt) Footprint(name string) (cols, rows int) {
	switch name {
	case "stall":
		return 3, 2
	case "well":
		return 2, 2
	}
	return 8, 7
}

// share is the frames a trade's villagers in sight spent in each state.
type share struct{ walk, work, still, resting, n int }

// measureDay plays the hours from..to with the player standing in the
// east town and tallies what its adults in sight do, by trade.
func measureDay(a Art, from, to float64) map[string]*share {
	g, f := build(a, still{})
	t := eastTown(f)
	f.player.Transform.Position = f.center(pathfind.Cell(t.Center))
	cv, _ := canvas.NewCanvas(1280, 720)
	g.Draw(cv)
	f.clock.Hour = from
	out := map[string]*share{}
	last := map[*villager]struct{ x, y float64 }{}
	for fr := 0; f.clock.Hour < to && f.clock.Day == 0; fr++ {
		g.Step(frame)
		if fr%6 != 0 {
			continue
		}
		for _, v := range f.awake() {
			if v.ours || v.obj.Hidden || v.stage() < 3 || !f.inSight(v.obj.Transform.Position) {
				continue
			}
			p := v.obj.Transform.Position
			l, seen := last[v]
			last[v] = struct{ x, y float64 }{p.X, p.Y}
			if !seen {
				continue
			}
			s := out[v.role]
			if s == nil {
				s = &share{}
				out[v.role] = s
			}
			s.n++
			switch {
			case math.Hypot(p.X-l.x, p.Y-l.y) > 0.5:
				s.walk++
			case v.hands.acting > 0:
				s.work++
			default:
				s.still++
				if v.mind.rest > 0 {
					s.resting++
				}
			}
		}
	}
	return out
}

func (s share) String() string {
	if s.n == 0 {
		return "no sample"
	}
	n := float64(s.n)
	return fmt.Sprintf("walk %2.0f%%  work %2.0f%%  still %2.0f%% (resting %2.0f%%)  (%d samples)", 100*float64(s.walk)/n, 100*float64(s.work)/n, 100*float64(s.still)/n, 100*float64(s.resting)/n, s.n)
}

// TestMeasureWhatVillagersDo reports, for the stylized layout and the
// realistic one, how the adults in sight spend the working day. It is
// opt-in: FARM_MEASURE=1.
func TestMeasureWhatVillagersDo(t *testing.T) {
	if os.Getenv("FARM_MEASURE") == "" {
		t.Skip("FARM_MEASURE not set")
	}
	for name, a := range map[string]Art{"stylized": SVGArt(), "realistic layout": roomyArt{SVGArt()}} {
		got := measureDay(a, 6, 21)
		roles := make([]string, 0, len(got))
		for r := range got {
			roles = append(roles, r)
		}
		sort.Strings(roles)
		for _, r := range roles {
			t.Logf("%-16s %-10s %v", name, r, *got[r])
		}
	}
}
