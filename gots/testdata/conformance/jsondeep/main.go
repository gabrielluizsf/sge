// Deeply nested JSON from an untrusted source: json refuses it past its
// depth limit with an error, and short of it decodes it. On Deno the limit
// is lower, so these are the depths Go and gots agree on; exits_test.go
// checks the ones in between are refused, not fatal.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func nested(n int) []byte {
	return []byte(strings.Repeat("[", n) + strings.Repeat("]", n))
}

func main() {
	for _, n := range []int{100, 900, 10001, 100000} {
		var v any
		err := json.Unmarshal(nested(n), &v)
		fmt.Println(n, err != nil, json.Valid(nested(n)))
		var raw json.RawMessage
		err = json.Unmarshal(nested(n), &raw)
		fmt.Println(n, "raw", err != nil)
	}
	obj := strings.Repeat(`{"a":`, 900) + "1" + strings.Repeat("}", 900)
	var m map[string]any
	fmt.Println("objects", json.Unmarshal([]byte(obj), &m) != nil)
	b, err := json.Marshal(m)
	fmt.Println("marshal back", len(b) > 0, err)
}
