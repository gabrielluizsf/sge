package gots

import "strings"

// asyncPrefix is the head of a function declaration or class method made
// async: export function f becomes export async function f, a method M
// async M.
func asyncPrefix(prefix string) string {
	if strings.Contains(prefix, "function ") {
		return strings.Replace(prefix, "function ", "async function ", 1)
	}
	return "async " + prefix
}
