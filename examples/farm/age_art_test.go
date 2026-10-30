package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/bond"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/life"
)

// The stylized art draws the old grey-haired, the young smaller, and an
// heir in the farmer's clothes with their own hair, a woman's long.
func TestTheArtDrawsAge(t *testing.T) {
	who := Who{Woman: true, Role: roleFarmer, Genes: [5]int{1, 2, 3, 4, 5}, Stage: life.Adult}
	if svgLook(who).Hair == grey {
		t.Error("a grown-up grey")
	}
	who.Stage = life.Elder
	if svgLook(who).Hair != grey {
		t.Error("an elder not grey")
	}
	a := SVGArt()
	who.Stage = life.Adult
	grown := height(a.Person(who))
	who.Stage = life.Child
	if child := height(a.Person(who)); child >= grown {
		t.Errorf("a child %d tall, a grown-up %d", child, grown)
	}
	heir := heirLook(Who{Player: true, Woman: true, Genes: [5]int{1, 2, 3, 4, 5}, Stage: life.Adult})
	if heir.Style != art.Long || !heir.Hat || heir.Hair != art.Villager(true, [5]int{1, 2, 3, 4, 5}).Hair {
		t.Errorf("the heir %+v", heir)
	}
}

// The HUD shows the player's age, and, talking, how old the one talked
// to is and how fond of the player in hearts.
func TestTheHUDShowsAgesAndHearts(t *testing.T) {
	g, f := newGame(idle{})
	v := single(t, f)
	v.bond.Affection = 2*perHeart + 10
	f.talk = &talk{v: v}
	cv, err := canvas.NewCanvas(1280, 720)
	if err != nil {
		t.Fatal(err)
	}
	g.Draw(cv)
	if heartsLabel(v) != "Conhecidos" || v.bond.Stage() != bond.Acquainted {
		t.Errorf("label %q", heartsLabel(v))
	}
	if got := f.labelOf(v); got != roleLabel(v) {
		t.Errorf("a grown villager labelled %q", got)
	}
}
