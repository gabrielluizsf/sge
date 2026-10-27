package main

import (
	"go/format"
	"log"
	"os"
)

// write writes src, formatted, to name.
func write(name string, src []byte) {
	out, err := format.Source(src)
	if err != nil {
		log.Fatalf("%s: %v\n%s", name, err, src)
	}
	if err := os.WriteFile(name, out, 0o644); err != nil {
		log.Fatal(err)
	}
}
