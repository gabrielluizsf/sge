// reflect over the values gots represents: types, fields and tags,
// setting through pointers, slices, maps and interfaces.
package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1fC", float64(c)) }

type Point struct {
	X, Y   int    `json:"x" db:"px"`
	Label  string `json:"label,omitempty"`
	hidden bool
	Temp   Celsius
	Tags   []string
	Next   *Point
}

type Shape interface{ Area() float64 }

type Square struct{ Side float64 }

func (s Square) Area() float64 { return s.Side * s.Side }

func main() {
	for _, v := range []any{1, int8(2), uint64(3), 1.5, "s", true, []int{1}, [2]bool{}, map[string]int{}, &Point{}, Point{}, Celsius(3), Square{}} {
		t := reflect.TypeOf(v)
		fmt.Printf("%-18s kind=%-7s name=%q\n", t.String(), t.Kind(), t.Name())
	}

	t := reflect.TypeOf(Point{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f.Name, f.Type, f.IsExported(), f.Tag.Get("json"), f.Tag.Get("db"), f.Index)
	}
	if f, ok := t.FieldByName("Label"); ok {
		v, ok := f.Tag.Lookup("json")
		fmt.Println("lookup", v, ok)
	}

	// Setting fields through a pointer.
	p := &Point{X: 1}
	v := reflect.ValueOf(p).Elem()
	v.Field(0).SetInt(10)
	v.FieldByName("Label").SetString("origin")
	v.Field(4).SetFloat(36.6)
	v.Field(5).Set(reflect.ValueOf([]string{"a", "b"}))
	fmt.Println(*p, v.Field(3).CanSet(), v.Field(0).CanSet(), reflect.ValueOf(*p).Field(0).CanSet())
	next := reflect.New(t)
	next.Elem().Field(1).SetInt(7)
	v.Field(6).Set(next)
	fmt.Println(p.Next.Y, v.Field(6).Elem().Field(1).Int())

	// Slices and arrays.
	s := []int{1, 2, 3}
	sv := reflect.ValueOf(s)
	sv.Index(1).SetInt(20)
	fmt.Println(s, sv.Len(), sv.Cap())
	grow := reflect.ValueOf(&s).Elem()
	grow.Grow(5)
	grow.SetLen(5)
	grow.Index(4).SetInt(50)
	fmt.Println(s, len(s))
	arr := [3]string{"x", "y", "z"}
	av := reflect.ValueOf(&arr).Elem()
	av.Index(2).SetString("Z")
	fmt.Println(arr, av.Len())
	made := reflect.MakeSlice(reflect.TypeOf([]float64{}), 2, 4)
	made.Index(0).SetFloat(1.25)
	fmt.Println(made.Interface(), made.Len(), made.Cap())

	// Maps.
	m := map[string]int{"a": 1}
	mv := reflect.ValueOf(m)
	mv.SetMapIndex(reflect.ValueOf("b"), reflect.ValueOf(2))
	mv.SetMapIndex(reflect.ValueOf("a"), reflect.Value{})
	var keys []string
	for it := mv.MapRange(); it.Next(); {
		keys = append(keys, fmt.Sprint(it.Key(), "=", it.Value()))
	}
	sort.Strings(keys)
	fmt.Println(keys, mv.MapIndex(reflect.ValueOf("b")), mv.MapIndex(reflect.ValueOf("zz")).IsValid())
	made2 := reflect.MakeMap(reflect.TypeOf(map[int]bool{}))
	made2.SetMapIndex(reflect.ValueOf(3), reflect.ValueOf(true))
	fmt.Println(made2.Interface(), made2.Len())

	// Interfaces and methods.
	shapeT := reflect.TypeFor[Shape]()
	fmt.Println(reflect.TypeOf(Square{}).Implements(shapeT), reflect.TypeOf(Point{}).Implements(shapeT), shapeT.Kind())
	var sh Shape = Square{3}
	shv := reflect.ValueOf(&sh).Elem()
	fmt.Println(shv.Kind(), shv.Elem().Type(), shv.Elem().Field(0).Float())
	shv.Set(reflect.ValueOf(Square{4}))
	fmt.Println(sh.Area())
	stringer := reflect.TypeFor[fmt.Stringer]()
	fmt.Println(reflect.TypeOf(Celsius(0)).Implements(stringer), reflect.TypeOf(1).Implements(stringer))

	// Interface round trips, zero values and kinds.
	x := reflect.ValueOf(Celsius(21.5)).Interface()
	fmt.Println(x, x.(fmt.Stringer).String())
	fmt.Println(reflect.Zero(t).Interface(), reflect.Zero(reflect.TypeOf(0)).Int(), reflect.ValueOf(Point{}).IsZero(), reflect.ValueOf(p).IsZero())
	var nilp *Point
	// A nil *Point in an interface is nil itself to gots (see README), so
	// it is reached through a pointer to the variable.
	fmt.Println(reflect.ValueOf(&nilp).Elem().IsNil(), reflect.ValueOf(nil).IsValid(), reflect.Indirect(reflect.ValueOf(p)).Type())
	if got, ok := reflect.TypeAssert[Celsius](reflect.ValueOf(Celsius(2))); ok {
		fmt.Println("assert", got)
	}
	fmt.Println(reflect.ValueOf(int8(7)).Convert(reflect.TypeOf(1.0)).Float(), reflect.ValueOf(65).Convert(reflect.TypeOf("")).String())
	fmt.Println(reflect.ValueOf(uint16(9)).Uint(), reflect.ValueOf(int8(0)).OverflowInt(200), reflect.TypeOf(uint8(0)).Size(), reflect.TypeOf(Point{}).Size())
	fmt.Println(reflect.DeepEqual(map[string][]int{"a": {1}}, map[string][]int{"a": {1}}), reflect.DeepEqual(Point{X: 1}, Point{X: 2}))
	fmt.Println(strings.Repeat("-", 3), reflect.TypeOf(p) == reflect.PointerTo(t), reflect.TypeOf(s).Elem() == reflect.TypeOf(0))
}
