package bvh

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// tokens are a file's words, one at a time.
type tokens struct{ s *bufio.Scanner }

// next is the next word; io.ErrUnexpectedEOF past the last.
func (t *tokens) next() (string, error) {
	if t.s.Scan() {
		return t.s.Text(), nil
	}
	if err := t.s.Err(); err != nil {
		return "", err
	}
	return "", errTruncated
}

// expect reads the word want, or fails.
func (t *tokens) expect(want string) error {
	w, err := t.next()
	if err == nil && w != want {
		err = fmt.Errorf("bvh: %q where %q was due", w, want)
	}
	return err
}

// number is the next word as a finite number.
func (t *tokens) number() (float64, error) {
	w, err := t.next()
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(w, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, fmt.Errorf("bvh: %q is not a finite number", w)
	}
	return v, nil
}

// count is the next word as a whole number from 0 to most.
func (t *tokens) count(most int) (int, error) {
	w, err := t.next()
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(w)
	if err != nil || n < 0 || n > most {
		return 0, fmt.Errorf("bvh: %q is not a count from 0 to %d", w, most)
	}
	return n, nil
}

var (
	errTruncated = errors.New("bvh: the file ends early")
	errTooLarge  = errors.New("bvh: the file is larger than its limit")
)
