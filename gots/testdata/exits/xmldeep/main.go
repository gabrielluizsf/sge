// XML nested deeper than a JavaScript stack can follow: decoding it into
// a recursive struct or skipping it is refused with an error, never a
// stack overflow (the pattern of CVE-2026-56859 in encoding/xml).
package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Node struct {
	Kids []Node `xml:"n"`
}

func main() {
	for _, n := range []int{100, 300, 4000, 20000} {
		doc := []byte(strings.Repeat("<n>", n) + strings.Repeat("</n>", n))
		var v Node
		err := xml.Unmarshal(doc, &v)
		var any struct{}
		err2 := xml.Unmarshal(doc, &any)
		fmt.Println(n, err == nil, err2 == nil)
	}
}
