// Decodes the gob stream in the file GOB_INPUT names into a type it does
// not match, so the decoder skips what it holds: a stream nested deeper
// than a JavaScript stack takes must be refused, not end the program.
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

func main() {
	b, err := os.ReadFile(os.Getenv("GOB_INPUT"))
	if err != nil {
		fmt.Println("read:", err)
		return
	}
	var out struct{ Hello int }
	err = gob.NewDecoder(bytes.NewReader(b)).Decode(&out)
	fmt.Println("refused:", err != nil)
}
