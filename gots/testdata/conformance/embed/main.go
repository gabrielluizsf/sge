package main

import "fmt"

type Base struct {
	ID   int
	name string
}

func (b Base) Describe() string { return fmt.Sprintf("#%d %s", b.ID, b.name) }
func (b *Base) Rename(n string) { b.name = n }

type Level int

func (l Level) Label() string { return fmt.Sprint("L", int(l)) }
func (l *Level) Up()          { *l++ }

type Player struct {
	Base
	Level
	Score int
}

type Boss struct {
	*Player
	Rage float64
}

type Describer interface{ Describe() string }

func main() {
	p := Player{Base: Base{1, "ana"}, Level: 3}
	fmt.Println(p.Describe(), p.ID, p.Label())
	p.Rename("bia")
	p.Up()
	fmt.Println(p.Describe(), p.Label(), p.Level)

	var d Describer = p
	p.Rename("changed")
	fmt.Println(d.Describe(), p.Describe())
	d = &p
	p.Rename("live")
	fmt.Println(d.Describe())

	b := Boss{Player: &p, Rage: 9.5}
	b.Rename("boss")
	b.Up()
	b.Score = 100
	fmt.Println(b.Describe(), p.name, p.Score, b.Label())

	desc := p.Describe
	p.Rename("later")
	fmt.Println(desc())
	rename := p.Rename
	rename("via method value")
	fmt.Println(p.name)
	f := Base.Describe
	fmt.Println(f(Base{2, "expr"}))
	g := (*Base).Rename
	g(&p.Base, "expr2")
	fmt.Println(p.name)
	fmt.Printf("%v\n%+v\n", p, b.Player.Base)
}

type inner struct{ W, H int }
type outer struct {
	inner
	Name string
}

func init() {
	// Go 1.27 allows a struct literal to set promoted fields.
	o := outer{W: 3, Name: "box", H: 4}
	fmt.Println(o.W, o.H, o.Name, o.inner)
}
