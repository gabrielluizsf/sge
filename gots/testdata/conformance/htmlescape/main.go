// html/template escaping must be Go's to the byte, or it is an XSS: the
// attack patterns of CVE-2023-24538 (actions in JS template literals),
// CVE-2026-32289 (branches in them) and go.dev/issue/62398 (<script,
// <!-- and </script in script strings), and every context in between.
package main

import (
	"fmt"
	"html/template"
	"os"
)

const evil = `</script><script>alert(1)</script><!--<script ' " ` + "`${x}`" + ` & javascript:alert(2)  `

var pages = []string{
	`<p>{{.}}</p>`,
	`<p title="{{.}}">x</p>`,
	`<p title={{.}}>x</p>`,
	`<a href="{{.}}">x</a>`,
	`<a href="/s?q={{.}}&r=1">x</a>`,
	`<a href="javascript:{{.}}">x</a>`,
	`<img srcset="{{.}} 2x">`,
	`<a onclick="f({{.}})">x</a>`,
	`<script>var s = "{{.}}";</script>`,
	`<script>var s = '{{.}}';</script>`,
	`<script>var v = {{.}};</script>`,
	`<script>var re = /{{.}}/;</script>`,
	"<script>var t = `{{.}}`;</script>",
	"<script>var t = `a${ {{.}} }b`;</script>",
	"<script>var t = `${ {{if .}}1{{else}}2{{end}} }`;</script>",
	`<script>// {{.}}` + "\n" + `</script>`,
	`<style>p { color: {{.}} }</style>`,
	`<p style="color: {{.}}">x</p>`,
	`<!-- {{.}} --><p>ok</p>`,
	`<textarea>{{.}}</textarea>`,
	`<title>{{.}}</title>`,
}

func main() {
	for i, p := range pages {
		t, err := template.New(fmt.Sprint(i)).Parse(p)
		if err != nil {
			fmt.Println(i, "parse error:", err)
			continue
		}
		if err := t.Execute(os.Stdout, evil); err != nil {
			fmt.Println(i, "exec error:", err)
			continue
		}
		fmt.Println()
	}
	safe := template.Must(template.New("safe").Parse(`{{.H}} {{.U}} {{.J}}`))
	safe.Execute(os.Stdout, map[string]any{
		"H": template.HTML("<b>trusted</b>"),
		"U": template.URL("javascript:void(0)"),
		"J": template.JS("1+1"),
	})
	fmt.Println()
}
