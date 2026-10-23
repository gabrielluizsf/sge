package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	s := "héllo, 世界!"
	fmt.Println(len(s), utf8.RuneCountInString(s), s[1], s[:2] == "h\xc3")
	for i, r := range s {
		fmt.Print(i, ":", string(r), " ")
	}
	fmt.Println()
	runes := []rune(s)
	fmt.Println(len(runes), string(runes[7:9]))
	fmt.Println(strings.ToUpper(s), strings.Contains(s, "世"), strings.Index(s, "l"))
	fmt.Println(strings.Split("a,b,,c", ","), strings.Join([]string{"x", "y"}, "-"))
	fmt.Println(strings.Repeat("ab", 3), strings.TrimSpace("  pad  "), strings.Fields(" a  b c "))
	fmt.Println(strings.HasPrefix(s, "hé"), strings.HasSuffix(s, "!"), strings.Replace("aaaa", "a", "b", 2))
	var sb strings.Builder
	for i := 0; i < 3; i++ {
		sb.WriteString("go")
		sb.WriteByte('!')
	}
	fmt.Println(sb.String(), sb.Len())
	b := []byte(s)
	fmt.Println(len(b), b[0], string(b[1:3]))
	fmt.Printf("%q %x %X %s\n", "tab\tquote\"", "hi", []byte{1, 171}, []byte("bytes"))
	fmt.Println(strings.Count("cheese", "e"), strings.LastIndex("go gopher", "go"), strings.EqualFold("Go", "GO"))
	invalid := "\xff\xfe"
	for _, r := range invalid {
		fmt.Print(r, " ")
	}
	fmt.Println(utf8.ValidString(invalid), utf8.RuneLen('世'))
	fmt.Println("a" < "b", "abc" > "abd", "é" > "z")
}
