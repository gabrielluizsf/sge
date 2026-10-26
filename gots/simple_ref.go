package gots

import "regexp"

var plainRef = regexp.MustCompile(`^[A-Za-z_$][\w$]*(\.[A-Za-z_$][\w$]*)*$`)

// simpleJS reports whether js is a name or a chain of property reads,
// which can be written twice without doing anything twice.
func (e *emitter) simpleJS(js string) bool {
	return plainRef.MatchString(js)
}
