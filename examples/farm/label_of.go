package farm

import "github.com/gabrielluizsf/sge/life"

// labelOf is who v is to the player, on screen: their spouse, their
// child, or of their family; otherwise a child or a youth of the town,
// or their trade.
func (f *farm) labelOf(v *villager) string {
	switch {
	case v == f.self.spouse:
		return pick(v, "esposa", "marido")
	case f.childOf(v):
		return pick(v, "filha", "filho")
	case v.ours:
		return "família"
	case v.stage() < life.Teen:
		return pick(v, "menina", "menino")
	case v.stage() < life.Adult:
		return "jovem"
	}
	return roleLabel(v)
}

// childOf reports whether v is the player's child.
func (f *farm) childOf(v *villager) bool {
	for _, id := range f.family.Children(f.self.id) {
		if id == v.id {
			return true
		}
	}
	return false
}

// pick is woman's word for v if v is a woman, man's if not.
func pick(v *villager, woman, man string) string {
	if v.woman {
		return woman
	}
	return man
}
