package gots

import "go/types"

// derefPtr is js, a pointer of type ptr, dereferenced: a pointer to a
// struct or an array is the object itself; any other is a cell holding
// the value.
func derefPtr(js string, ptr types.Type) string {
	if isObjectPtr(ptr) {
		return js
	}
	return js + ".v"
}
