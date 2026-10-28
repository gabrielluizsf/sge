package main

import "fmt"

type T struct {
	A int
	B string
	c []float64
	P *T
}

func main() {
	t := T{A: 1, B: "x", c: []float64{1.5, 2}}
	fmt.Printf("%v|%+v\n", t, t)
	fmt.Printf("%v\n", &t)
	fmt.Printf("[%5d][%-5d][%05d][%x][%X][%o][%b][%#x][%+d]\n", 42, 42, 42, 255, 255, 8, 5, 255, 3)
	fmt.Printf("[%8.3f][%-8.2f][%e][%g][%G][%.0f][%6.2e]\n", 3.14159, 2.5, 123456.789, 1e21, 1e-7, 2.5, 0.000123)
	fmt.Printf("[%s][%10s][%-10s][%.2s][%q][%v]\n", "go", "right", "left", "trunc", 'x', 'x')
	fmt.Printf("[%t][%v][%c][%U][%%]\n", true, false, 0x4E16, 0x4E16)
	fmt.Printf("%v %v %v %v\n", []string{"a", "b"}, map[string]int{"z": 1, "a": 2}, [2]bool{true}, struct{}{})
	fmt.Printf("%v %d %s\n", nil, []int{1, 2}, []string{"q"})
	fmt.Printf("%d %s\n", "str", 5)
	fmt.Printf("%d\n")
	fmt.Printf("%d %d\n", 1, 2, 3)
	fmt.Printf("%T %T %T %T %T\n", 1, "s", 2.0, []any{}, map[int]bool{})
	fmt.Printf("%#v %#v %#v\n", t.c, "q", 42)
	fmt.Printf("%x %X % x\n", "hello", []byte("hi"), []byte{1, 2})
	fmt.Printf("%6.2f%%\n", 99.5)
	fmt.Printf("%*d|%-*d|\n", 5, 1, 4, 2)
	fmt.Println(1e6, 1e20, 1e21, 0.000001, 123456789.0, float32(1)/3, -0.0)
	s := fmt.Sprint("a", 1, 2, "b", "c", 3.5, nil)
	fmt.Println(s)
	fmt.Print("x", 1, "\n")
	fmt.Println(fmt.Sprintln("end", 1) == "end 1\n")
	var ip *int
	var e error
	fmt.Println(ip, e, []int(nil), map[string]int(nil))
}

func init() {
	var x, d int
	var f float64
	var s string
	var u8 uint8
	var u64 uint64
	n, err := fmt.Sscanf("ff 42 2.5e1 word", "%x %d %g %s", &x, &d, &f, &s)
	fmt.Println(n, err, x, d, f, s)
	n, err = fmt.Sscanf("7,8", "%d,%d", &x, &d)
	fmt.Println(n, err, x, d)
	n, err = fmt.Sscan("  10\n 0x1F  -3.25 str", &u8, &u64, &f, &s)
	fmt.Println(n, err, u8, u64, f, s)
	n, err = fmt.Sscanf("300", "%d", &u8)
	fmt.Println(n, err != nil, u8)
	n, err = fmt.Sscanf("abc", "%d", &x)
	fmt.Println(n, err != nil)
}
