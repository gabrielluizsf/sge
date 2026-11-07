package farm

import (
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/town"
)

// household is the family of a house: a man and a woman who live as a
// couple, and their child.
type household struct{ man, woman, child *villager }

// settler makes the people of one town from its seed.
type settler struct {
	f   *farm
	t   *town.Town
	src *random.Source
}

// post is where a villager is made: the house, the trade, their sex and
// the cell they stand on.
type post struct {
	home  *town.House
	role  string
	woman bool
}

// person is a villager at post, standing at their house's door.
func (s *settler) person(st post) *villager {
	names := menNames
	if st.woman {
		names = womenNames
	}
	h := st.home
	v := &villager{
		name: names[s.src.Intn(len(names))], woman: st.woman, role: st.role,
		who:  Who{Woman: st.woman, Role: st.role, Genes: [5]int{s.src.Intn(64), s.src.Intn(64), s.src.Intn(64), s.src.Intn(64), s.src.Intn(64)}},
		kit:  &kit{money: 15, axe: 0, slot: slotCan},
		town: s.t, home: h,
		pos:  s.f.ground.CellCenter(h.Door.Col, h.Door.Row),
		luck: random.New(uint32(s.src.Intn(1 << 30))),
		wait: float64(s.src.Intn(orderPause)),
	}
	v.kit.bag.Add(seedsOf(s.src.Intn(len(crops))), 2)
	v.mind = &mind{v: v}
	return v
}

// household makes the family of house h whose head has the trade role:
// the head is a man or a woman as the seed falls, their partner a
// farmer unless the trade is the farmer's or the wanderer's, and their
// child a wanderer, a son or a daughter.
func (s *settler) household(h *town.House, role string) household {
	headWoman := s.src.Float() < 0.5
	other := roleFarmer
	if role == roleFarmer || role == roleWanderer {
		other = role
	}
	head := s.person(post{home: h, role: role, woman: headWoman})
	mate := s.person(post{home: h, role: other, woman: !headWoman})
	kid := s.person(post{home: h, role: roleWanderer, woman: s.src.Float() < 0.5})
	hh := household{man: head, woman: mate, child: kid}
	if headWoman {
		hh.man, hh.woman = mate, head
	}
	if role == roleSmith {
		smith := head
		smith.kit.money, smith.fuel = 40, 2
		smith.pos = s.f.ground.CellCenter(s.t.Anvil.Col, s.t.Anvil.Row+1)
	}
	return hh
}
