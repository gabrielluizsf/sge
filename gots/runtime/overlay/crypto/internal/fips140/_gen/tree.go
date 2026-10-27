package main

// tree is the sum of terms, parenthesized as a balanced tree.
func tree(terms []string) string {
	if len(terms) == 1 {
		return terms[0]
	}
	h := len(terms) / 2
	return group(terms[:h]) + " + " + group(terms[h:])
}

func group(terms []string) string {
	if len(terms) == 1 {
		return terms[0]
	}
	return "(" + tree(terms) + ")"
}
