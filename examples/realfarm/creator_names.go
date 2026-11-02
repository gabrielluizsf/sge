package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gabrielluizsf/sge/human"
)

// eyeNames are the eye colours in Portuguese, by the file's name.
var eyeNames = map[string]string{
	"brown": "Castanhos", "brownlight": "Castanho-claros", "bluegreen": "Azul-esverdeados",
	"blue": "Azuis", "deepblue": "Azul-escuros", "lightblue": "Azul-claros",
	"green": "Verdes", "grey": "Cinzentos", "ice": "Gelo",
}

// eyeName is the Portuguese name of the eyes in material file item.
func eyeName(item string) string {
	n := strings.TrimSuffix(strings.TrimPrefix(item, "eyes/materials/"), ".mhmat")
	if s, ok := eyeNames[n]; ok {
		return s
	}
	return n
}

// ordinal is "i / n": where item is in list, or "-" when it is not.
func ordinal(list []string, item string) string {
	if i := slices.Index(list, item); i >= 0 {
		return fmt.Sprintf("%d / %d", i+1, len(list))
	}
	return "-"
}

// sexName is the sex in Portuguese.
func sexName(s human.Sex) string {
	if s == human.Man {
		return "Homem"
	}
	return "Mulher"
}
