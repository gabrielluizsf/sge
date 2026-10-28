package main

import "fmt"

type Number interface {
	~int | ~int32 | ~float64 | ~uint8
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func Index[T comparable](xs []T, v T) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return -1
}

func Div[T Number](a, b T) T { return a / b }

func Zero[T any]() T {
	var z T
	return z
}

type Pair[K comparable, V any] struct {
	Key K
	Val V
}

func (p Pair[K, V]) String() string { return fmt.Sprintf("%v=%v", p.Key, p.Val) }

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T)           { s[v] = struct{}{} }
func (s Set[T]) Has(v T) bool      { _, ok := s[v]; return ok }
func NewSet[T comparable]() Set[T] { return Set[T]{} }

type Point struct{ X, Y int }

func Max[T int | float64 | string](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Keys[K comparable, V any](m map[K]V) []K {
	var ks []K
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func Reduce[T, A any](xs []T, acc A, f func(A, T) A) A {
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

type List[T any] struct {
	head *node[T]
	size int
}

type node[T any] struct {
	val  T
	next *node[T]
}

func (l *List[T]) Prepend(v T) {
	l.head = &node[T]{v, l.head}
	l.size++
}

func (l *List[T]) Each(f func(T)) {
	for n := l.head; n != nil; n = n.next {
		f(n.val)
	}
}

func main() {
	fmt.Println(Sum([]int{1, 2, 3}), Sum([]float64{1.5, 2.25}), Sum([]uint8{200, 100}))
	fmt.Println(Div(7, 2), Div(7.0, 2), Div[int32](-7, 2))
	strs := Map([]int{1, 2, 3}, func(i int) string { return fmt.Sprint("n", i) })
	fmt.Println(strs, len(strs))
	fmt.Println(Index([]string{"a", "b"}, "b"), Index([]Point{{1, 2}, {3, 4}}, Point{3, 4}), Index([]float64{1}, 2))
	fmt.Println(Zero[int](), Zero[string]() == "", Zero[*Point]() == nil, Zero[Point]())

	p := Pair[string, int]{"age", 30}
	fmt.Println(p, p.String())
	var st fmt.Stringer = Pair[int, bool]{1, true}
	fmt.Println(st)

	var s Stack[Point]
	s.Push(Point{1, 1})
	s.Push(Point{2, 2})
	top, ok := s.Pop()
	fmt.Println(top, ok, s.Len())
	ints := &Stack[int]{}
	ints.Push(5)
	v, _ := ints.Pop()
	_, empty := ints.Pop()
	fmt.Println(v, empty)

	set := NewSet[string]()
	set.Add("x")
	fmt.Println(set.Has("x"), set.Has("y"), len(set))

	fmt.Println(Max(3, 7), Max(2.5, 1.5), Max("go", "deno"))
	fmt.Println(len(Keys(map[string]int{"a": 1, "b": 2})))
	total := Reduce([]int{1, 2, 3, 4}, "", func(acc string, x int) string { return acc + fmt.Sprint(x) })
	fmt.Println(total)

	var l List[string]
	l.Prepend("c")
	l.Prepend("b")
	l.Prepend("a")
	l.Each(func(s string) { fmt.Print(s, " ") })
	fmt.Println(l.size)
	fmt.Printf("%T %T %v\n", p, &s, Pair[Point, []int]{Point{1, 2}, []int{3}})

	var anyPair any = Pair[string, int]{"k", 1}
	if pp, ok := anyPair.(Pair[string, int]); ok {
		fmt.Println("asserted", pp.Key)
	}
	_, wrong := anyPair.(Pair[int, int])
	fmt.Println(wrong)
}

// Container is a generic interface whose methods mention its own type
// parameter, used from inside generic code.
type Container[T any] interface {
	Put(T)
	All() []T
}

type bag[T any] struct{ items []T }

func (b *bag[T]) Put(v T)  { b.items = append(b.items, v) }
func (b *bag[T]) All() []T { return b.items }

func fill[T any](c Container[T], vs ...T) Container[T] {
	for _, v := range vs {
		c.Put(v)
	}
	var anon interface{ All() []T } = c
	_ = anon
	return c
}

func init() {
	fmt.Println(fill[int](&bag[int]{}, 1, 2, 3).All())
}
