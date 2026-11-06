package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gabrielluizsf/sge/examples/i18n"
	"github.com/gabrielluizsf/sge/human"
)

// eyeNames are the eye colours as English keys, by the file's name.
var eyeNames = map[string]string{
	"brown": "Brown", "brownlight": "Light brown", "bluegreen": "Blue-green",
	"blue": "Blue", "deepblue": "Dark blue", "lightblue": "Light blue",
	"green": "Green", "grey": "Grey", "ice": "Ice",
}

// eyeName is the name of the eyes in material file item.
func eyeName(item string) string {
	n := strings.TrimSuffix(strings.TrimPrefix(item, "eyes/materials/"), ".mhmat")
	if s, ok := eyeNames[n]; ok {
		return i18n.T(s)
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

// sexName is the sex.
func sexName(s human.Sex) string {
	if s == human.Man {
		return i18n.T("Man")
	}
	return i18n.T("Woman")
}
