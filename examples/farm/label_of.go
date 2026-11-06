package farm

import (
	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/life"
)

// labelOf is who v is to the player, on screen: their spouse, their
// child, or of their family; otherwise a child or a youth of the town,
// or their trade.
func (f *farm) labelOf(v *villager) string {
	switch {
	case v == f.self.spouse:
		return i18n.TF(v.woman, "wife")
	case f.childOf(v):
		return i18n.TF(v.woman, "daughter")
	case v.ours:
		return i18n.T("family")
	case v.stage() < life.Teen:
		return i18n.TF(v.woman, "girl")
	case v.stage() < life.Adult:
		return i18n.T("youth")
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
