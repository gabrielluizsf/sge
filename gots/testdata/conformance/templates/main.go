// text/template and html/template, which call methods and functions
// through reflect.
package main

import (
	"fmt"
	htmltemplate "html/template"
	"os"
	"strings"
	"text/template"
)

type Item struct {
	Name  string
	Price float64
	Tags  []string
}

func (i Item) Label() string               { return strings.ToUpper(i.Name) }
func (i *Item) Discount(p float64) float64 { return i.Price * (1 - p/100) }

type Page struct {
	Title string
	Items []*Item
	Meta  map[string]int
	Note  string
}

const text = `{{define "row"}}- {{.Label}}: {{printf "%.2f" .Price}} ({{.Discount 10 | printf "%.1f"}}){{range $i, $t := .Tags}}{{if $i}},{{end}} {{$t}}{{end}}
{{end}}{{.Title | upper}}
{{range .Items}}{{template "row" .}}{{else}}none{{end}}{{with .Meta}}{{range $k, $v := .}}{{$k}}={{$v}};{{end}}{{end}}
{{len .Items}} {{index .Items 0 | label}} {{if and .Title (gt (len .Items) 1)}}many{{end}} {{range 3}}{{.}}{{end}}
{{block "foot" .}}{{.Note | html}}{{end}}
`

func main() {
	funcs := template.FuncMap{
		"upper": strings.ToUpper,
		"label": func(i *Item) string { return "[" + i.Label() + "]" },
	}
	t := template.Must(template.New("page").Funcs(funcs).Parse(text))
	p := Page{
		Title: "loja",
		Items: []*Item{{Name: "lápis", Price: 2, Tags: []string{"a", "b"}}, {Name: "caneta", Price: 3.5}},
		Meta:  map[string]int{"z": 1, "a": 2},
		Note:  "<b>&</b>",
	}
	if err := t.Execute(os.Stdout, p); err != nil {
		fmt.Println("error:", err)
	}
	err := template.Must(template.New("x").Parse("{{.Nope}}")).Execute(os.Stdout, p)
	fmt.Println("\nerror:", err)

	h := htmltemplate.Must(htmltemplate.New("h").Parse(`<a href="/s?q={{.}}" title="{{.}}">{{.}}</a><script>var x = {{.}};</script>` + "\n"))
	h.Execute(os.Stdout, `O'Reilly <x> & "y"`)
}
