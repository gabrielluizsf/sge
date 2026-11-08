package save

import (
	"errors"
	"strings"
)

// ErrName is returned for a file name that is not a plain, fixed name.
var ErrName = errors.New("save: bad file name")

// maxName is the longest file name, in bytes.
const maxName = 64

// checkName reports whether name is a plain name of lowercase letters,
// digits and dashes: no separators, no dots, so never a path.
func checkName(name string) error {
	if name == "" || len(name) > maxName || strings.Trim(name, "abcdefghijklmnopqrstuvwxyz0123456789-") != "" {
		return ErrName
	}
	return nil
}
