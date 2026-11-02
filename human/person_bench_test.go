package human

import "testing"

func BenchmarkMake(b *testing.B) {
	k := realKit(&testing.T{})
	c := baselines[0]
	for b.Loop() {
		if _, err := k.Make(c.body, c.look); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRecolour(b *testing.B) {
	t := &testing.T{}
	k := realKit(t)
	c := baselines[0]
	p, err := k.MakePerson(c.body, c.look)
	if err != nil {
		b.Fatal(err)
	}
	l := c.look
	i := 0
	for b.Loop() {
		i++
		l.SkinTone = 0.1 + float64(i%8)/10
		p.Recolour(l)
	}
}

func BenchmarkRecolourClothes(b *testing.B) {
	k := realKit(&testing.T{})
	c := baselines[0]
	p, _ := k.MakePerson(c.body, c.look)
	l := c.look
	l.ClothesColours = make([]Colour, len(l.Wear))
	i := 0
	for b.Loop() {
		i++
		l.ClothesColours[4] = RGB(float64(i%9)/9, 0.3, 0.2)
		p.Recolour(l)
	}
}
