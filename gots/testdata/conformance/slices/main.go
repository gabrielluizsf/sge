package main

import "fmt"

func main() {
	var s []int
	fmt.Println(s == nil, len(s), cap(s))
	for i := 0; i < 10; i++ {
		s = append(s, i*i)
	}
	fmt.Println(s, len(s))

	t := s[2:5]
	t[0] = -1
	fmt.Println(s[2], t, len(t), cap(t))

	u := append(t, 100)
	fmt.Println(s[5], u)

	full := s[2:5:5]
	full = append(full, 200)
	full[0] = 7
	fmt.Println(s[2], s[5], full)

	m := make([]string, 3, 10)
	m[1] = "b"
	fmt.Printf("%q %d %d\n", m, len(m), cap(m))

	dst := make([]int, 3)
	n := copy(dst, s)
	fmt.Println(n, dst)
	copy(s[1:], s)
	fmt.Println(s)

	b := []byte("héllo")
	b[0] = 'H'
	fmt.Println(string(b), len(b))
	b = append(b, " world"...)
	fmt.Println(string(b))

	grid := make([][]int, 3)
	for i := range grid {
		grid[i] = make([]int, 3)
		for j := range grid[i] {
			grid[i][j] = i * j
		}
	}
	fmt.Println(grid)

	var arr [5]int
	sl := arr[1:3]
	sl[0] = 9
	fmt.Println(arr, sl)

	words := []string{"a", "b", "c", "d"}
	words = append(words[:1], words[2:]...)
	fmt.Println(words, len(words))

	clear(dst)
	fmt.Println(dst)

	var u8 []uint8
	u8 = append(u8, 250, 251)
	u8[0] += 10
	fmt.Println(u8)

	lit := []int{5: 1, 2, 0: 3}
	fmt.Println(lit, len(lit))
}
