package pathfind

import "math"

// steps are the eight ways to a neighbour, straight ones first.
var steps = [8][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

// Path is the shortest way from one cell to another, A*: the cells to walk
// through, the last being to, the first not being from. A diagonal step is
// taken only when both cells beside it are open, so no corner is cut. It
// is nil when to is not open or cannot be reached within the grid's Limit,
// and empty when from is to.
func (g Grid) Path(from, to Cell) []Cell {
	if from == to {
		return []Cell{}
	}
	if !g.Open(to.Col, to.Row) {
		return nil
	}
	start, ok := key(from)
	if !ok {
		return nil
	}
	came := map[int]Cell{start: from}
	cost := map[int]float64{start: 0}
	q := queue{{c: from, rank: guess(from, to)}}
	for looked := 0; len(q) > 0 && looked < g.Limit; looked++ {
		n := q.pop()
		if n.c == to {
			return walkBack(came, from, to)
		}
		k, _ := key(n.c)
		if n.cost > cost[k] {
			continue // a cheaper way here was found after this was queued
		}
		for i, s := range steps {
			next := Cell{n.c.Col + s[0], n.c.Row + s[1]}
			if !g.Open(next.Col, next.Row) || i >= 4 && !(g.Open(n.c.Col+s[0], n.c.Row) && g.Open(n.c.Col, n.c.Row+s[1])) {
				continue
			}
			nk, ok := key(next)
			if !ok {
				continue
			}
			c := n.cost + map[bool]float64{true: 1, false: math.Sqrt2}[i < 4]
			if old, seen := cost[nk]; seen && old <= c {
				continue
			}
			cost[nk], came[nk] = c, n.c
			q.push(node{c: next, cost: c, rank: c + guess(next, to)})
		}
	}
	return nil
}

// guess is how far a is from b at the least, in steps of the grid: the
// diagonal steps there must be, and the straight ones after.
func guess(a, b Cell) float64 {
	dx, dy := math.Abs(float64(a.Col-b.Col)), math.Abs(float64(a.Row-b.Row))
	return max(dx, dy) + (math.Sqrt2-1)*min(dx, dy)
}

// walkBack is the way from from to to, from the cell each was come from.
func walkBack(came map[int]Cell, from, to Cell) []Cell {
	var way []Cell
	for c := to; c != from; {
		way = append(way, c)
		k, _ := key(c)
		c = came[k]
	}
	for i, j := 0, len(way)-1; i < j; i, j = i+1, j-1 {
		way[i], way[j] = way[j], way[i]
	}
	return way
}
