package human

import (
	"math"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/gabrielluizsf/sge/gltf"
)

// units is a head and its jaw, two poses: at rest, and the jaw dropped a
// quarter about X.
const unitsBVH = `HIERARCHY
ROOT root
{
	OFFSET 0 0 0
	CHANNELS 3 Xrotation Yrotation Zrotation
	JOINT head
	{
		OFFSET 0 0 1
		CHANNELS 3 Xrotation Yrotation Zrotation
		JOINT jaw
		{
			OFFSET 0 -0.1 0
			CHANNELS 3 Xrotation Yrotation Zrotation
			End Site
			{
				OFFSET 0 -0.1 0
			}
		}
	}
}
MOTION
Frames: 3
Frame Time: 1
0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 90 0 0
0 0 0 0 0 0 0 90 0
`

// faceKit is a kit of no person, but the face pose units, under lic.
func faceKit(lic string) *Kit {
	fsys := fstest.MapFS{
		faceUnitsBVH:  {Data: []byte(unitsBVH)},
		faceUnitsJSON: {Data: []byte(`{"license": "` + lic + `", "framemapping": ["Rest", "JawDrop", "JawTurn"]}`)},
	}
	return &Kit{fsys: fsys, lim: DefaultLimits}
}

// faceModel is a model of the bones the units turn.
func faceModel() *gltf.Model {
	n := func(name string) gltf.Node {
		return gltf.Node{Name: name, Mesh: -1, Skin: -1, Rotation: [4]float64{3: 1}, Scale: [3]float64{1, 1, 1}}
	}
	return &gltf.Model{Nodes: []gltf.Node{n("root"), n("head"), n("jaw")}, Roots: []int{0}}
}

// An expression turns the face's bones as its units do, in the model's
// axes, Y up, each part of the way its weight says; the bones of the body
// and the model it is made from are left as they were.
func TestAnExpressionTurnsTheFace(t *testing.T) {
	k, m := faceKit("CC0"), faceModel()
	e, err := k.Express(m, Expression{"JawDrop": 1})
	if err != nil {
		t.Fatal(err)
	}
	q := quat(e.Nodes[2].Rotation)
	if want := math.Sqrt(0.5); math.Abs(q[0]-want) > 1e-9 || math.Abs(q[3]-want) > 1e-9 {
		t.Errorf("the jaw dropped a quarter about X turns %v", q)
	}
	if e.Nodes[0].Rotation != m.Nodes[0].Rotation || m.Nodes[2].Rotation != [4]float64{3: 1} {
		t.Error("the body, or the model expressed, was turned")
	}
	half, _ := k.Express(m, Expression{"JawDrop": 0.5})
	if a := 2 * math.Acos(half.Nodes[2].Rotation[3]); math.Abs(a-math.Pi/4) > 1e-9 {
		t.Errorf("half the unit turns %v, not an eighth", a)
	}
	// About the file's Y, which points back into the head, Z up: about
	// the model's -Z.
	turned, _ := k.Express(m, Expression{"JawTurn": 1})
	if r := turned.Nodes[2].Rotation; r[2] > -0.7 || math.Abs(r[1]) > 1e-9 {
		t.Errorf("a turn about the file's Y is %v, not about the model's -Z", r)
	}
	if _, err := k.Express(m, Expression{"Wink": 1}); err == nil {
		t.Error("an expression of no unit")
	}
}

// Pose units under any licence but CC0 are not read.
func TestFaceUnitsMustBeCC0(t *testing.T) {
	if _, err := faceKit("AGPL3").Express(faceModel(), Lids(1)); err == nil || !strings.Contains(err.Error(), "CC0") {
		t.Errorf("read units not CC0: %v", err)
	}
}
