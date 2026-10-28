package main

import (
	"errors"
	"fmt"
)

type Shape interface {
	Area() float64
	Name() string
}

type Square struct{ Side float64 }

func (s Square) Area() float64 { return s.Side * s.Side }
func (s Square) Name() string  { return "square" }

type Circle struct{ R float64 }

func (c *Circle) Area() float64 { return 3 * c.R * c.R }
func (c *Circle) Name() string  { return "circle" }

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1f°C", float64(c)) }

type Counter int

func (c *Counter) Inc() { *c++ }

type MyErr struct{ Code int }

func (e *MyErr) Error() string { return fmt.Sprintf("code %d", e.Code) }

func find(code int) error {
	if code == 0 {
		return nil
	}
	return &MyErr{code}
}

func describe(x any) string {
	switch v := x.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int %d", v+1)
	case string, []byte:
		return fmt.Sprintf("text %v", v)
	case Square:
		return fmt.Sprintf("square %v", v.Side)
	case *Circle:
		return "circle ptr " + v.Name()
	case fmt.Stringer:
		return "stringer " + v.String()
	case error:
		return "error " + v.Error()
	default:
		return fmt.Sprintf("other %T", v)
	}
}

func main() {
	shapes := []Shape{Square{2}, &Circle{1}}
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
		fmt.Println(s.Name(), s.Area())
	}
	fmt.Println(total)

	sq := Square{3}
	var s Shape = sq
	sq.Side = 100
	fmt.Println(s.Area(), sq.Area())

	if c, ok := s.(Square); ok {
		fmt.Println("is square", c.Side)
	}
	_, isCircle := s.(*Circle)
	fmt.Println(isCircle)

	fmt.Println(describe(nil), describe(41), describe("hi"), describe([]byte("b")))
	fmt.Println(describe(Square{1}), describe(&Circle{2}), describe(Celsius(21.5)))
	fmt.Println(describe(errors.New("boom")), describe(3.5), describe([]int{1}))

	var c Counter
	c.Inc()
	c.Inc()
	cp := &c
	cp.Inc()
	fmt.Println(c, *cp)

	err := find(7)
	var me *MyErr
	if errors.As(err, &me) {
		fmt.Println("as", me.Code)
	}
	fmt.Println(find(0) == nil, err != nil, err)
	wrapped := fmt.Errorf("wrap: %w", err)
	fmt.Println(wrapped, errors.Is(wrapped, err), errors.Unwrap(wrapped) == err)

	var a, b any = 1, 1
	var d any = "1"
	fmt.Println(a == b, a == d, a == 1, Square{1} == Square{1})
	var sh1, sh2 Shape = Square{1}, Square{1}
	fmt.Println(sh1 == sh2)

	fmt.Printf("%v %s %d %T %T\n", Celsius(10), Celsius(-3), Counter(4), Celsius(1), shapes)
	var st fmt.Stringer = Celsius(0)
	fmt.Println(st)
}
