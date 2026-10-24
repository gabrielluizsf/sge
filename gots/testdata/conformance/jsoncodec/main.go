// encoding/json, reflect's main user: structs with tags, embedded
// fields, pointers, maps, slices, interfaces and custom marshalers.
package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Base struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
}

type Address struct {
	City string `json:"city"`
	Zip  string `json:"zip,omitempty"`
}

type Level int

func (l Level) MarshalText() ([]byte, error) { return []byte(strings.Repeat("*", int(l))), nil }

func (l *Level) UnmarshalText(b []byte) error {
	*l = Level(len(b))
	return nil
}

type User struct {
	Base
	Name     string             `json:"name"`
	Email    string             `json:"email,omitempty"`
	Age      int                `json:"age,string"`
	Admin    bool               `json:"-"`
	Address  *Address           `json:"address"`
	Tags     []string           `json:"tags"`
	Scores   map[string]float64 `json:"scores"`
	Level    Level              `json:"level"`
	Extra    any                `json:"extra"`
	Raw      json.RawMessage    `json:"raw"`
	internal int
}

func main() {
	u := User{
		Base:    Base{ID: 7, Created: time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)},
		Name:    "Ana \"Go\" <ts>",
		Age:     30,
		Admin:   true,
		Address: &Address{City: "São Paulo"},
		Tags:    []string{"a", "b"},
		Scores:  map[string]float64{"z": 1.5, "a": 2},
		Level:   3,
		Extra:   map[string]any{"n": 1, "list": []any{true, nil, "x"}},
		Raw:     json.RawMessage(`{"pre":"made"}`),
	}
	b, err := json.Marshal(u)
	fmt.Println(string(b), err)
	b, _ = json.MarshalIndent(Address{City: "Rio", Zip: "20000"}, "", "  ")
	fmt.Println(string(b))

	var back User
	err = json.Unmarshal([]byte(`{"id":9,"created":"2020-01-02T03:04:05Z","name":"Bia","age":"41","address":{"city":"Recife"},"tags":["x"],"scores":{"k":3},"level":"**","extra":[1,"two",{"three":3}],"raw":[1, 2],"unknown":true}`), &back)
	fmt.Println(err)
	fmt.Printf("%d %s %q %d %v %v %v %d %v %s\n", back.ID, back.Created.Format(time.RFC3339), back.Name, back.Age, *back.Address, back.Tags, back.Scores, back.Level, back.Extra, back.Raw)

	var anything any
	json.Unmarshal([]byte(`{"a":[1,2.5,"s",null,false],"b":{"c":{}}}`), &anything)
	fmt.Println(anything)

	var nums []int
	err = json.Unmarshal([]byte(`[1, 2, "three"]`), &nums)
	fmt.Println(nums, err)
	err = json.Unmarshal([]byte(`{"id": "seven"}`), &back)
	fmt.Println(err)
	fmt.Println(json.Valid([]byte(`{"ok":true}`)), json.Valid([]byte(`{nope}`)))

	dec := json.NewDecoder(strings.NewReader(`{"city":"A"} {"city":"B"}`))
	for {
		var a Address
		if err := dec.Decode(&a); err != nil {
			break
		}
		fmt.Print(a.City, " ")
	}
	fmt.Println()
}
