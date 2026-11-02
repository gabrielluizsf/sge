package bvh

import (
	"errors"
	"fmt"
)

// endSite marks an End Site open on the stack of joints being read.
const endSite = -1

// readHierarchy reads the joints, from HIERARCHY to MOTION, keeping the
// joints open in a stack, not by recursion, as deep as l allows.
func readHierarchy(t *tokens, m *Motion, l Limits) error {
	if err := t.expect("HIERARCHY"); err != nil {
		return err
	}
	var open []int
	for {
		w, err := t.next()
		if err != nil {
			return err
		}
		switch w {
		case "ROOT", "JOINT":
			if err = openJoint(t, m, &open, w == "ROOT", l); err != nil {
				return err
			}
		case "End":
			if err = openEnd(t, &open, l); err != nil {
				return err
			}
		case "OFFSET":
			err = readOffset(t, m, open)
		case "CHANNELS":
			err = readChannels(t, m, open)
		case "}":
			if len(open) == 0 {
				return errors.New("bvh: a brace closes nothing")
			}
			open = open[:len(open)-1]
		case "MOTION":
			if len(open) > 0 || len(m.Joints) == 0 {
				return errors.New("bvh: MOTION with joints open, or none")
			}
			return nil
		default:
			return fmt.Errorf("bvh: %q in the hierarchy", w)
		}
		if err != nil {
			return err
		}
	}
}

// openJoint reads a joint's name and its brace, a root only where no
// joint is open, any other inside one, within l's joints and depth.
func openJoint(t *tokens, m *Motion, open *[]int, root bool, l Limits) error {
	name, err := t.next()
	if err != nil {
		return err
	}
	if err := t.expect("{"); err != nil {
		return err
	}
	parent := -1
	if len(*open) > 0 {
		parent = (*open)[len(*open)-1]
	}
	switch {
	case root != (len(*open) == 0), parent == endSite && !root:
		return fmt.Errorf("bvh: joint %q out of place", name)
	case len(m.Joints) >= l.Joints || len(*open) >= l.Depth:
		return fmt.Errorf("bvh: more than %d joints, or deeper than %d", l.Joints, l.Depth)
	case m.Joint(name) >= 0:
		return fmt.Errorf("bvh: joint %q twice", name)
	}
	m.Joints = append(m.Joints, Joint{Name: name, Parent: parent})
	*open = append(*open, len(m.Joints)-1)
	return nil
}

// openEnd reads an End Site's brace, inside a joint, within l's depth.
func openEnd(t *tokens, open *[]int, l Limits) error {
	if err := t.expect("Site"); err != nil {
		return err
	}
	if err := t.expect("{"); err != nil {
		return err
	}
	if len(*open) == 0 || (*open)[len(*open)-1] == endSite || len(*open) >= l.Depth {
		return errors.New("bvh: an End Site out of place, or too deep")
	}
	*open = append(*open, endSite)
	return nil
}

// readOffset reads the open joint's offset; an End Site's is read and
// let be.
func readOffset(t *tokens, m *Motion, open []int) error {
	if len(open) == 0 {
		return errors.New("bvh: OFFSET outside a joint")
	}
	var v [3]float64
	for i := range v {
		var err error
		if v[i], err = t.number(); err != nil {
			return err
		}
	}
	if j := open[len(open)-1]; j != endSite {
		m.Joints[j].Offset = v
	}
	return nil
}

// readChannels reads the open joint's channels, once, at most six, each
// a channel's name.
func readChannels(t *tokens, m *Motion, open []int) error {
	if len(open) == 0 || open[len(open)-1] == endSite {
		return errors.New("bvh: CHANNELS outside a joint")
	}
	j := &m.Joints[open[len(open)-1]]
	n, err := t.count(6)
	if err != nil {
		return err
	}
	if j.Channels != nil {
		return fmt.Errorf("bvh: joint %q has its channels twice", j.Name)
	}
	j.Channels = make([]Channel, 0, n)
	for range n {
		w, err := t.next()
		if err != nil {
			return err
		}
		c, ok := channelNames[w]
		if !ok {
			return fmt.Errorf("bvh: %q is no channel", w)
		}
		j.Channels = append(j.Channels, c)
	}
	return nil
}
