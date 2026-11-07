package bakes

import (
	"math"
	"os"
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/kit"
	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/pose"
)

// A person's ankles, posed on the real skeleton, are a natural step apart
// while they walk, and do not splay as MakeHuman's rest stands: a step's
// width, ankle to ankle, is 10 to 20 centimetres in a healthy adult's
// stride, never the 40 of the rest pose. Each stands as far forward as the
// stride is long, some 70 centimetres from the other's.
func TestWalkKeepsTheFeetNearTheMiddle(t *testing.T) {
	if testing.Short() || os.Getenv("CI") != "" {
		t.Skip("makes a person from MakeHuman's data, downloaded the first time")
	}
	h, err := kit.ReadHumans("../../humans.json")
	if err != nil {
		t.Fatal(err)
	}
	makers, closeAll, err := h.Kit(kit.Fetcher())
	if err != nil {
		t.Skip("no MakeHuman data to make a person from: ", err)
	}
	defer closeAll()
	v := kit.Villagers[0]
	m, err := makers.Make(v.At(30), v.Look)
	if err != nil {
		t.Fatal(err)
	}
	if err := human.Animate(m, human.Walk); err != nil {
		t.Fatal(err)
	}
	nodes := map[string]int{}
	for i, n := range m.Nodes {
		nodes[n.Name] = i
	}
	p := pose.New(m)
	walk := p.Animation("walk")
	still := pose.TRS([3]float64{}, [4]float64{0, 0, 0, 1}, [3]float64{1, 1, 1})
	longest := 0.0
	for k := range 16 {
		w := p.At(walk, human.Walk.Length*float64(k)/16, still)
		l, r := w[nodes["foot.L"]].Point([3]float64{}), w[nodes["foot.R"]].Point([3]float64{})
		if width := math.Abs(l[0] - r[0]); width < 0.10 || width > 0.20 {
			t.Errorf("at %d of 16 the ankles are %.3f metres apart sideways", k, width)
		}
		longest = max(longest, math.Abs(l[2]-r[2]))
	}
	if longest < 0.6 || longest > 0.9 {
		t.Errorf("the longest step is %.2f metres", longest)
	}
}
