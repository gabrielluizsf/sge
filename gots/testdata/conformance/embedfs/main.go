package main

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed assets/greeting.txt
var greeting string

//go:embed assets/greeting.txt
var raw []byte

//go:embed assets
var assets embed.FS

func main() {
	fmt.Println(greeting, len(raw), raw[0])
	data, err := assets.ReadFile("assets/sub/lines.txt")
	fmt.Printf("%q %v\n", data, err)
	_, err = assets.ReadFile("assets/missing")
	fmt.Println(err)
	entries, _ := assets.ReadDir("assets")
	for _, en := range entries {
		fmt.Println(en.Name(), en.IsDir())
	}
	var names []string
	fs.WalkDir(assets, ".", func(p string, d fs.DirEntry, err error) error {
		names = append(names, p)
		return err
	})
	fmt.Println(strings.Join(names, ","))
	b, err := fs.ReadFile(assets, "assets/greeting.txt")
	fmt.Println(string(b), err)
}
