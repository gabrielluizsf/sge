// encoding/gob: structs, maps, slices, pointers and interfaces there and
// back, and the bytes Go writes.
package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
)

type Point struct{ X, Y int }

type Shape interface{ Area() float64 }

type Rect struct{ W, H float64 }

func (r Rect) Area() float64 { return r.W * r.H }

type Save struct {
	Name   string
	Score  int64
	Ratio  float64
	Points []Point
	Tags   map[string]bool
	Best   *Point
	Shapes []Shape
	Blob   []byte
	skip   int
}

func main() {
	gob.Register(Rect{})
	in := Save{
		Name: "farm", Score: -42, Ratio: 0.25,
		Points: []Point{{1, 2}, {3, 4}},
		Tags:   map[string]bool{"seeds": true},
		Best:   &Point{9, 9},
		Shapes: []Shape{Rect{2, 3}},
		Blob:   []byte{0, 1, 255},
		skip:   7,
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		fmt.Println("encode:", err)
		return
	}
	fmt.Println(len(buf.Bytes()), hex.EncodeToString(buf.Bytes()[:24]))
	var out Save
	if err := gob.NewDecoder(&buf).Decode(&out); err != nil {
		fmt.Println("decode:", err)
		return
	}
	fmt.Println(out.Name, out.Score, out.Ratio, out.Points, out.Tags, *out.Best, out.Shapes[0].Area(), out.Blob, out.skip)
	var bad Save
	fmt.Println(gob.NewDecoder(bytes.NewReader([]byte{1, 2, 3})).Decode(&bad) != nil)
}
