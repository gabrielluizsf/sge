package sound

// Seamless makes s loop without a click where its end meets its start:
// its last n samples are crossfaded into its first n, and cut off. It
// is the loop, n samples shorter.
func Seamless(s []float32, n int) []float32 {
	n = min(n, len(s)/2)
	end := len(s) - n
	for i := range n {
		f := float32(i) / float32(n)
		s[i] = s[i]*f + s[end+i]*(1-f)
	}
	return s[:end]
}
