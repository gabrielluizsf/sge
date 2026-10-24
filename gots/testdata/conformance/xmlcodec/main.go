// encoding/xml through reflect: attributes, nested elements, character
// data, omitempty, and decoding back.
package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	ID      int      `xml:"id,attr"`
	Name    string   `xml:"name"`
	Tags    []string `xml:"tags>tag"`
	Note    string   `xml:",comment"`
	Price   float64  `xml:"price,omitempty"`
}

type Catalog struct {
	XMLName xml.Name `xml:"catalog"`
	Owner   string   `xml:"owner,attr"`
	Items   []Item   `xml:"item"`
}

func main() {
	c := Catalog{Owner: "Zé & cia", Items: []Item{{ID: 1, Name: "lápis <2B>", Tags: []string{"a", "b"}, Note: " ok ", Price: 1.5}, {ID: 2, Name: "caneta"}}}
	b, err := xml.MarshalIndent(c, "", " ")
	fmt.Println(string(b), err)
	var back Catalog
	err = xml.Unmarshal(b, &back)
	fmt.Println(err, back.Owner, len(back.Items), back.Items[0].Name, back.Items[0].Tags, back.Items[1].Price)
	d := xml.NewDecoder(strings.NewReader(`<a x="1"><b>t</b><!-- c --></a>`))
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		fmt.Printf("%T ", tok)
	}
	fmt.Println()
	var bad Catalog
	fmt.Println(xml.Unmarshal([]byte(`<catalog><item id="x"/></catalog>`), &bad) != nil)
}
