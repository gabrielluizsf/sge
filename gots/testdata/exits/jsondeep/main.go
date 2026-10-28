// JSON nested deeper than a JavaScript stack can follow, though not past
// Go's own limit: on Deno json refuses it with an error rather than the
// program dying of a stack overflow on input it cannot control.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	deep := strings.Repeat("[", 5000) + strings.Repeat("]", 5000)
	var v any
	fmt.Println("arrays refused:", json.Unmarshal([]byte(deep), &v) != nil)
	var raw json.RawMessage
	fmt.Println("raw refused:", json.Unmarshal([]byte(deep), &raw) != nil)
	obj := strings.Repeat(`{"a":`, 5000) + "1" + strings.Repeat("}", 5000)
	var m map[string]any
	fmt.Println("objects refused:", json.Unmarshal([]byte(obj), &m) != nil)
	fmt.Println("valid:", json.Valid([]byte(deep)))

	// A value built in Go that deep is refused on the way out too.
	var nest any = 1
	for range 5000 {
		nest = []any{nest}
	}
	_, err := json.Marshal(nest)
	fmt.Println("marshal refused:", err != nil)
}
