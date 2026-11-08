package toast

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func TestLifeFollowsLength(t *testing.T) {
	short, long := Life("Jogo salvo"), Life("uma frase bem mais longa que conta muitas palavras ao jogador agora")
	if short < 2 || short > 2.4 || long != 3.5 {
		t.Errorf("short %v long %v", short, long)
	}
	if mid := Life("um dois tres quatro cinco seis"); mid <= short || mid >= long {
		t.Errorf("mid %v", mid)
	}
}

func TestAShortOneIsGoneAfterAboutTwoSecondsAndTheExit(t *testing.T) {
	l := Add(nil, "Jogo salvo")
	for range 60 * 2 {
		l = Step(l, 1.0/60)
	}
	if len(l) != 1 || !l[0].Leaving() && l[0].Age < 1.9 {
		t.Fatalf("at 2 s: %v", l)
	}
	for range 60 / 2 {
		l = Step(l, 1.0/60)
	}
	if len(l) != 0 {
		t.Errorf("still up at 2.5 s")
	}
}

func TestOpacityAndPositionThroughEnterHoldAndExit(t *testing.T) {
	x := &Toast{Life: 2}
	if x.Opacity() != 0 || x.Shift() != -1 {
		t.Errorf("start %v %v", x.Opacity(), x.Shift())
	}
	x.Age = EnterFor / 2
	if o := x.Opacity(); o <= 0.5 || o >= 1 { // eased out: past halfway at half time
		t.Errorf("mid enter %v", o)
	}
	if s := x.Shift(); s >= 0 || s <= -0.5 {
		t.Errorf("mid enter shift %v", s)
	}
	x.Age = 1
	if x.Opacity() != 1 || x.Shift() != 0 {
		t.Errorf("hold %v %v", x.Opacity(), x.Shift())
	}
	x.Age = 2 + ExitFor/2
	if o := x.Opacity(); math.Abs(o-0.5) > 1e-9 || x.Shift() >= 0 {
		t.Errorf("mid exit %v %v", o, x.Shift())
	}
	x.Age = 2 + ExitFor
	if x.Opacity() != 0 || !x.Gone() {
		t.Errorf("end")
	}
}

func TestStackedToastsCapAndCloseTheGap(t *testing.T) {
	var l []*Toast
	for _, s := range []string{"a", "b", "c", "d"} {
		l = Add(l, s)
	}
	if !l[0].Leaving() || l[1].Leaving() || len(l) != 4 {
		t.Fatalf("cap: first leaving %v", l[0].Leaving())
	}
	for range 60 {
		l = Step(l, 1.0/60)
	}
	if len(l) != 3 || l[0].Text != "b" {
		t.Fatalf("after: %d", len(l))
	}
	// the gap closes smoothly: b moved up from row 1 to 0 without a jump
	before := l[0].Row
	l[0].Age = l[0].Life // b leaves; c and d should rise gradually
	l = Step(l, 1.0/60)
	c := l[1]
	if !(c.Row < 2 && c.Row > c.Target) || before < 0 {
		t.Errorf("c row %v target %v", c.Row, c.Target)
	}
}

func TestBlitAlphaScalesOpacity(t *testing.T) {
	src, _ := canvas.NewCanvas(2, 2)
	src.FillRect(canvas.Area{Width: 2, Height: 2}, canvas.Color(0xFFFFFFFF))
	dst, _ := canvas.NewCanvas(2, 2)
	dst.FillRect(canvas.Area{Width: 2, Height: 2}, canvas.Color(0xFF000000))
	BlitAlpha(dst, 0, 0, src, 0.5)
	if r := dst.At(0, 0).R(); r < 126 || r > 129 {
		t.Errorf("red %d", r)
	}
}
