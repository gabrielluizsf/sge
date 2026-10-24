package terrain

// Memo is in, remembering its answers for up to size cells, then
// forgetting them all and starting over: an area asked about the same
// cells many times, as Open and Near ask, while a chunk is made.
func Memo(in func(col, row int) bool, size int) func(col, row int) bool {
	seen := map[int]bool{}
	return func(col, row int) bool {
		k, ok := cellKey(col, row)
		if !ok {
			return in(col, row)
		}
		if v, ok := seen[k]; ok {
			return v
		}
		if len(seen) >= size {
			clear(seen)
		}
		v := in(col, row)
		seen[k] = v
		return v
	}
}
