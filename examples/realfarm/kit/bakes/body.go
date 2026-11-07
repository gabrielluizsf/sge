package bakes

import (
	"fmt"

	"github.com/gabrielluizsf/sge/examples/realfarm/kit"
	"github.com/gabrielluizsf/sge/human"
)

// bodyOf is villager v made as they are at stage g, animated in all of
// the stage's motions. The last one made is kept, as a villager's motions
// are baked one after another.
func (b *Baker) bodyOf(v kit.Villager, g kit.Grown) (*stageBody, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	key := stageKey{v.ID, g.Stage, kit.SpriteTag(v.ID)}
	if b.body != nil && b.last == key {
		return b.body, nil
	}
	makers, err := b.humanKit()
	if err != nil {
		return nil, err
	}
	if err := onlyCC0(makers, v); err != nil {
		return nil, err
	}
	m, err := makers.Make(v.At(g.Years), v.Look)
	if err == nil {
		err = human.Animate(m, g.Motions...)
		if err == nil {
			err = b.holdAxe(m, g.Motions)
		}
		if err == nil {
			err = b.holdForge(m, g.Motions)
		}
		if err == nil {
			err = b.holdHoe(m, g.Motions)
		}
		if err == nil {
			err = b.holdCan(m, g.Motions)
		}
		if err == nil {
			err = b.holdPickaxe(m, g.Motions)
		}
		if err == nil {
			err = b.holdMorsel(m, g.Motions)
		}
		if err == nil {
			err = b.holdFillingCan(m, g.Motions)
		}
		if err == nil {
			err = b.seatOnChair(m, g.Motions)
		}
		if err == nil {
			err = b.layInBed(m, g.Motions)
		}
	}
	if err != nil {
		return nil, err
	}
	b.last, b.body = key, &stageBody{model: m}
	return b.body, nil
}

// humanKit is MakeHuman's data and packs, fetched and opened once. It is
// kept open for as long as the baker lives.
func (b *Baker) humanKit() (*human.Kit, error) {
	b.makeOnce.Do(func() { b.makers, _, b.makeErr = b.Humans.Kit(kit.Fetcher()) })
	return b.makers, b.makeErr
}

// onlyCC0 refuses a villager made of any file that does not say it is
// CC0, so a game sold with them owes no one: a pack's name is no promise.
func onlyCC0(makers *human.Kit, v kit.Villager) error {
	lics, err := makers.Licences(v.Body, v.Look)
	if err != nil {
		return fmt.Errorf("%s: %w", v.ID, err)
	}
	for name, lic := range lics {
		if lic != "CC0" {
			return fmt.Errorf("%s: %s is under %q, not CC0", v.ID, name, lic)
		}
	}
	return nil
}
