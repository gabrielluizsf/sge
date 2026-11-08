package bvh

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// Read reads a motion from r within l: its hierarchy of joints, then its
// frames.
func Read(r io.Reader, l Limits) (*Motion, error) {
	data, err := io.ReadAll(io.LimitReader(r, l.Bytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > l.Bytes {
		return nil, errTooLarge
	}
	s := bufio.NewScanner(bytes.NewReader(data))
	s.Split(bufio.ScanWords)
	t := &tokens{s: s}
	m := &Motion{}
	if err := readHierarchy(t, m, l); err != nil {
		return nil, err
	}
	if err := readFrames(t, m, int64(len(data)), l); err != nil {
		return nil, err
	}
	return m, nil
}

// readFrames reads the frames after MOTION: as many as the file says, no
// more than l allows, each with a value for every channel, and nothing
// after them. The values must fit in the file's size, two bytes each at
// the least, before room is made for them.
func readFrames(t *tokens, m *Motion, size int64, l Limits) error {
	if err := t.expect("Frames:"); err != nil {
		return err
	}
	n, err := t.count(l.Frames)
	if err != nil {
		return err
	}
	if err := t.expect("Frame"); err != nil {
		return err
	}
	if err := t.expect("Time:"); err != nil {
		return err
	}
	if m.FrameTime, err = t.number(); err != nil {
		return err
	}
	channels := 0
	for _, j := range m.Joints {
		channels += len(j.Channels)
	}
	if int64(n)*int64(channels) > size/2 {
		return fmt.Errorf("bvh: %d frames of %d channels cannot fit in %d bytes", n, channels, size)
	}
	m.Frames = make([][]float64, n)
	for f := range m.Frames {
		m.Frames[f] = make([]float64, channels)
		for c := range m.Frames[f] {
			if m.Frames[f][c], err = t.number(); err != nil {
				return err
			}
		}
	}
	if w, err := t.next(); err == nil {
		return fmt.Errorf("bvh: %q after the last frame", w)
	}
	return nil
}
