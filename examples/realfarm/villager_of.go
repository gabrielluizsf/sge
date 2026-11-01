package main

import "github.com/gabrielluizsf/sge/examples/farm"

// folk are the villagers that may be who, by gender and trade: the
// genes pick among them, so two farmers are not always the same person.
var folk = map[bool]map[string][]string{
	false: {
		"farmer":     {"joao", "player"},
		"smith":      {"beto"},
		"lumberjack": {"rui"},
		"wanderer":   {"kenji", "joao", "rui", "beto"},
		"player":     {"player", "joao"},
	},
	true: {
		"farmer":     {"ana", "iara"},
		"smith":      {"mei", "iara"},
		"lumberjack": {"iara", "ana"},
		"wanderer":   {"mei", "lucia", "iara", "ana"},
		"player":     {"iara", "ana"},
	},
}

// villagerOf is the id of the villager that who is drawn as. The founding
// farmer, whose numbers are all 0, is the player.
func villagerOf(who farm.Who) string {
	role := who.Role
	if who.Player {
		role = "player"
	}
	list := folk[who.Woman][role]
	if len(list) == 0 {
		list = folk[who.Woman]["wanderer"]
	}
	return list[who.Genes[0]%len(list)]
}
