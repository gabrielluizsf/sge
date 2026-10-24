// Calling through reflect: functions, variadic ones, methods of structs,
// pointers and other named types, by index and by name.
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Greeter struct{ Name string }

func (g Greeter) Hello(to string) string { return "hello " + to + " from " + g.Name }
func (g *Greeter) Rename(n string)       { g.Name = n }
func (g Greeter) Pair() (int, error)     { return len(g.Name), errors.New("pair") }

type Celsius float64

func (c Celsius) F() float64 { return float64(c)*9/5 + 32 }

type Counter int

func (c *Counter) Inc(by ...int) int {
	for _, b := range by {
		*c += Counter(b)
	}
	return int(*c)
}

func sum(xs ...int) int {
	t := 0
	for _, x := range xs {
		t += x
	}
	return t
}

func main() {
	up := reflect.ValueOf(strings.ToUpper)
	fmt.Println(up.Call([]reflect.Value{reflect.ValueOf("gots")})[0], up.Type().NumIn(), up.Type().In(0), up.Type().Out(0))

	s := reflect.ValueOf(sum)
	fmt.Println(s.Call(nil)[0], s.Call([]reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)})[0], s.Type().IsVariadic())
	fmt.Println(s.CallSlice([]reflect.Value{reflect.ValueOf([]int{4, 5})})[0])

	g := &Greeter{Name: "Ana"}
	gv := reflect.ValueOf(g)
	fmt.Println(gv.NumMethod(), gv.Elem().NumMethod(), reflect.TypeOf(g).Method(0).Name, reflect.TypeOf(g).Method(0).Type)
	fmt.Println(gv.MethodByName("Hello").Call([]reflect.Value{reflect.ValueOf("Bia")})[0])
	gv.MethodByName("Rename").Call([]reflect.Value{reflect.ValueOf("Caio")})
	fmt.Println(g.Name, gv.Elem().MethodByName("Rename").IsValid())
	outs := gv.MethodByName("Pair").Call(nil)
	fmt.Println(outs[0].Int(), outs[1].Interface().(error), outs[1].Type())
	m, ok := reflect.TypeOf(Greeter{}).MethodByName("Hello")
	fmt.Println(m.Name, ok, m.Type.NumIn(), m.Index)

	fmt.Println(reflect.ValueOf(Celsius(100)).MethodByName("F").Call(nil)[0])
	var c Counter = 1
	inc := reflect.ValueOf(&c).MethodByName("Inc")
	fmt.Println(inc.Call([]reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)})[0], c, inc.Type())

	var i fmt.Stringer = reflect.TypeOf(0)
	fmt.Println(reflect.ValueOf(&i).Elem().MethodByName("String").Call(nil)[0])

	n := 0
	for v := range reflect.ValueOf(3).Seq() {
		n += int(v.Int())
	}
	for k, v := range reflect.ValueOf([]string{"a", "b"}).Seq2() {
		fmt.Print(k, v, " ")
	}
	fmt.Println(n, reflect.TypeOf(sum).CanSeq(), reflect.TypeOf([]int{}).CanSeq2())

	defer func() { fmt.Println("recovered:", recover() != nil) }()
	reflect.ValueOf(sum).Call([]reflect.Value{reflect.ValueOf("not an int")})
}
