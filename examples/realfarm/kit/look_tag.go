package kit

import "strings"

// bakedEyes are the eyes each villager was first baked with. A villager
// whose eyes are still these keeps their baked pictures; one whose eyes
// changed is baked again, under a task name of their own.
var bakedEyes = map[string]string{
	"joao": eyes("brown"), "ana": eyes("brown"), "beto": eyes("brown"),
	"rui": eyes("bluegreen"), "mei": eyes("brown"), "lucia": eyes("green"),
	"kenji": eyes("brown"), "iara": eyes("brownlight"), "player": eyes("brown"),
}

// LookTag is what a task name of villager v adds to tell their look from
// the one baked before: "" while their eyes are the first ones.
func (v Villager) LookTag() string {
	if t := PlayerTag(v); t != "" {
		return t
	}
	if e := v.Look.Eyes; e != bakedEyes[v.ID] {
		return " eyes " + strings.TrimSuffix(strings.TrimPrefix(e, "eyes/materials/"), ".mhmat")
	}
	return ""
}
