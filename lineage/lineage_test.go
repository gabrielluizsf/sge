package lineage

import (
	"errors"
	"testing"
)

// family is grandparents ana and joao; their children bia, born first,
// and caio; bia's child with rui, duda; and caio's with mei, edu.
func family(t *testing.T) *Tree {
	t.Helper()
	var tr Tree
	for _, p := range []Person{
		{ID: "ana"}, {ID: "joao"}, {ID: "rui"}, {ID: "mei"}, {ID: "zeca"},
		{ID: "bia", Parents: [2]string{"ana", "joao"}, Born: 10},
		{ID: "caio", Parents: [2]string{"ana", "joao"}, Born: 20},
		{ID: "duda", Parents: [2]string{"bia", "rui"}, Born: 40},
		{ID: "edu", Parents: [2]string{"caio", "mei"}, Born: 45},
	} {
		if err := tr.Add(p); err != nil {
			t.Fatal(err)
		}
	}
	return &tr
}

func TestKinNeverCourt(t *testing.T) {
	tr := family(t)
	for _, pair := range [][2]string{{"ana", "bia"}, {"bia", "caio"}, {"ana", "duda"}, {"caio", "duda"}, {"duda", "edu"}, {"rui", "duda"}, {"bia", "bia"}} {
		if !tr.Kin(pair[0], pair[1]) || !tr.Kin(pair[1], pair[0]) {
			t.Errorf("%v not kin", pair)
		}
	}
	for _, pair := range [][2]string{{"ana", "joao"}, {"rui", "mei"}, {"zeca", "bia"}, {"rui", "caio"}} {
		if tr.Kin(pair[0], pair[1]) {
			t.Errorf("%v kin", pair)
		}
	}
}

// The eldest living child is heir; with no child living, a grandchild;
// with none, no one.
func TestTheEldestChildTakesOver(t *testing.T) {
	tr := family(t)
	if h, ok := tr.Heir("ana"); !ok || h != "bia" {
		t.Errorf("ana's heir %q", h)
	}
	tr.Die("bia")
	if h, _ := tr.Heir("ana"); h != "caio" {
		t.Errorf("bia dead, ana's heir %q", h)
	}
	tr.Die("caio")
	if h, _ := tr.Heir("ana"); h != "duda" {
		t.Errorf("both dead, ana's heir %q", h)
	}
	if _, ok := tr.Heir("zeca"); ok {
		t.Error("zeca, childless, has an heir")
	}
}

// No one is added twice, before their parents, or of one parent twice.
func TestTheTreeHasNoLoops(t *testing.T) {
	tr := family(t)
	for _, p := range []Person{{ID: "ana"}, {ID: "x", Parents: [2]string{"nobody"}}, {ID: "y", Parents: [2]string{"rui", "rui"}}, {ID: ""}} {
		if err := tr.Add(p); err == nil {
			t.Errorf("added %+v", p)
		}
	}
	if err := tr.Die("nobody"); !errors.Is(err, ErrUnknown) {
		t.Errorf("an unknown died: %v", err)
	}
	if c := tr.Children("ana"); len(c) != 2 || c[0] != "bia" {
		t.Errorf("ana's children %v", c)
	}
}

func TestAChildIsBornWhenDue(t *testing.T) {
	e := Expect([2]string{"bia", "rui"}, 5)
	if e.Born(5+Gestation-1) || !e.Born(5+Gestation) {
		t.Errorf("due %d", e.Due)
	}
}
