package main

import "fmt"

// machine mirrors compress/flate's decoder: gotos into labels forward and
// backward, out of a switch, from inside nested blocks.
func machine(input []int, start int) []string {
	var out []string
	i := 0
	switch start {
	case 0:
		goto read
	case 1:
		goto copyBack
	}
read:
	{
		if i >= len(input) {
			return append(out, "eof")
		}
		v := input[i]
		i++
		switch {
		case v < 10:
			out = append(out, fmt.Sprint("lit", v))
			goto read
		case v == 10:
			return append(out, "end")
		}
		out = append(out, fmt.Sprint("ref", v))
		goto copyBack
	}
copyBack:
	out = append(out, "copy")
	if i < len(input) {
		goto read
	}
	return out
}

func forward(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		if i == 2 {
			goto done
		}
		s += fmt.Sprint(i)
	}
	s += "never"
done:
	return s + "!"
}

func loopWithGotos(rows [][]int) (sum int) {
	for _, r := range rows {
		j := 0
	again:
		if j < len(r) {
			if r[j] < 0 {
				break
			}
			if r[j] == 0 {
				j++
				continue
			}
			sum += r[j]
			j++
			goto again
		}
	}
	return sum
}

func main() {
	fmt.Println(machine([]int{1, 2, 15, 3, 10, 4}, 0))
	fmt.Println(machine([]int{5}, 1))
	fmt.Println(machine(nil, 0))
	fmt.Println(forward(5), forward(1))
	fmt.Println(loopWithGotos([][]int{{1, 2, 0, 3}, {4, -1, 100}, {0, 5}}))
}
