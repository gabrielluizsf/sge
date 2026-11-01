package bakes

import "sync"

// each runs do for 0 to n-1, as many at once as the baker's Workers
// allow, one when it has none set, and gives the first error.
func (b *Baker) each(n int, do func(i int) error) error {
	workers := max(1, min(b.Workers, n))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var first error
	next := 0
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				i := next
				next++
				stop := first != nil
				mu.Unlock()
				if i >= n || stop {
					return
				}
				if err := do(i); err != nil {
					mu.Lock()
					if first == nil {
						first = err
					}
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	return first
}
