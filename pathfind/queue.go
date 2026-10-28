package pathfind

// node is a cell on the way being searched: its cost so far, and its
// cost so far plus a guess of what is left.
type node struct {
	c          Cell
	cost, rank float64
}

// queue is the cells to look at next, the lowest rank first: a binary
// heap.
type queue []node

func (q *queue) push(n node) {
	*q = append(*q, n)
	h := *q
	for i := len(h) - 1; i > 0; {
		p := (i - 1) / 2
		if h[p].rank <= h[i].rank {
			break
		}
		h[p], h[i] = h[i], h[p]
		i = p
	}
}

func (q *queue) pop() node {
	h := *q
	top := h[0]
	last := len(h) - 1
	h[0] = h[last]
	h = h[:last]
	for i := 0; ; {
		l, r, m := 2*i+1, 2*i+2, i
		if l < len(h) && h[l].rank < h[m].rank {
			m = l
		}
		if r < len(h) && h[r].rank < h[m].rank {
			m = r
		}
		if m == i {
			break
		}
		h[m], h[i] = h[i], h[m]
		i = m
	}
	*q = h
	return top
}
