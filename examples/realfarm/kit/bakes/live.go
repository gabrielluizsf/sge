package bakes

import (
	"fmt"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/human"
)

// Person is character c standing, made for a screen to draw live and to
// recolour as the player changes the colours: a stance, no bake. The
// human kit is the baker's, so what was read for a bake is not read again.
func (b *Baker) Person(c kit.Character) (*human.Person, error) {
	if err := c.Check(); err != nil {
		return nil, err
	}
	makers, err := b.humanKit()
	if err != nil {
		return nil, err
	}
	lics, err := makers.Licences(c.Body, c.Look)
	if err != nil {
		return nil, err
	}
	for name, lic := range lics {
		if lic != "CC0" {
			return nil, fmt.Errorf("person: %s is under %q, not CC0", name, lic)
		}
	}
	p, err := makers.MakePerson(c.Body, c.Look)
	if err == nil {
		err = human.Animate(p.Model, human.Stand)
	}
	return p, err
}
