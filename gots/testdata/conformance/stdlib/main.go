package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/bits"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type byLen []string

func (s byLen) Len() int           { return len(s) }
func (s byLen) Less(i, j int) bool { return len(s[i]) < len(s[j]) }
func (s byLen) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type person struct {
	name string
	age  int
}

func main() {
	fmt.Println(strconv.Itoa(-42), strconv.FormatInt(255, 16), strconv.FormatFloat(3.14159, 'f', 2, 64))
	n, err := strconv.Atoi("123")
	_, bad := strconv.Atoi("12a")
	fmt.Println(n, err, bad)
	f, _ := strconv.ParseFloat("2.5e3", 64)
	b, _ := strconv.ParseBool("true")
	u, _ := strconv.ParseUint("ff", 16, 64)
	fmt.Println(f, b, u, strconv.Quote("hi\n\"there\""), strconv.FormatFloat(1.0/3, 'g', -1, 64))
	q, _ := strconv.Unquote(`"a\tb"`)
	fmt.Println(q, strconv.AppendInt(nil, 7, 10), strconv.FormatFloat(1e21, 'e', 3, 64))

	xs := []int{5, 2, 8, 1, 9, 3}
	slices.Sort(xs)
	idx, found := slices.BinarySearch(xs, 8)
	fmt.Println(xs, idx, found, slices.Contains(xs, 4), slices.Index(xs, 9), slices.Max(xs), slices.Min(xs))
	ys := slices.Clone(xs)
	ys[0] = 100
	slices.Reverse(ys)
	fmt.Println(xs[0], ys, slices.Equal(xs, ys))
	people := []person{{"ana", 30}, {"bob", 25}, {"cid", 30}, {"dan", 20}}
	slices.SortStableFunc(people, func(a, b person) int { return a.age - b.age })
	fmt.Println(people)
	fmt.Println(slices.Compact([]int{1, 1, 2, 2, 2, 3, 1}), slices.Insert([]int{1, 4}, 1, 2, 3))

	words := []string{"banana", "kiwi", "apple", "fig"}
	sort.Sort(byLen(words))
	fmt.Println(words)
	sort.Strings(words)
	fmt.Println(words, sort.SearchStrings(words, "fig"))
	sort.Slice(people, func(i, j int) bool { return people[i].name > people[j].name })
	fmt.Println(people)
	fl := []float64{2.5, -1, 3}
	sort.Float64s(fl)
	fmt.Println(fl, sort.IsSorted(sort.Float64Slice(fl)))

	var buf bytes.Buffer
	buf.WriteString("hello ")
	buf.Write([]byte("world"))
	fmt.Fprintf(&buf, " %d!", 42)
	fmt.Println(buf.String(), buf.Len())
	fmt.Println(bytes.Equal([]byte("a"), []byte("a")), bytes.Contains([]byte("seafood"), []byte("foo")))
	fmt.Println(string(bytes.ToUpper([]byte("shout"))), bytes.Split([]byte("a,b,c"), []byte(",")))

	fmt.Println(bits.OnesCount(255), bits.LeadingZeros32(1), bits.TrailingZeros64(8), bits.Len(1023), bits.Reverse8(1))
	hi, lo := bits.Mul64(1<<63, 4)
	fmt.Println(hi, lo, bits.RotateLeft32(1, -1))

	fmt.Println(unicode.IsUpper('Ä'), unicode.IsDigit('٣'), unicode.ToLower('Σ'), string(unicode.ToTitle('ǆ')), unicode.IsSpace(' '))
	fmt.Println(strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) {
			return -1
		}
		return unicode.ToUpper(r)
	}, "héllo, wörld!"))

	base := errors.New("base")
	e1 := fmt.Errorf("layer1: %w", base)
	e2 := fmt.Errorf("layer2: %w", e1)
	fmt.Println(e2, errors.Is(e2, base), errors.Join(base, errors.New("other")))
	fmt.Println(math.Floor(-2.5), math.Round(2.5), math.Trunc(-2.7), math.Mod(7, 3), math.Inf(1), math.MaxInt32, math.Pi)
	fmt.Println(strings.NewReplacer("a", "1", "b", "2").Replace("abcab"), strings.Title("go gopher"))
	r := strings.NewReader("line one\nline two")
	data := make([]byte, 4)
	k, _ := r.Read(data)
	fmt.Println(k, string(data), r.Len())
	fmt.Println(strings.Cut("key=value", "="))
	fmt.Println(strings.SplitN("a.b.c.d", ".", 2), strings.TrimFunc("123abc456", unicode.IsDigit), strings.ContainsAny("hello", "xyz"))
}
