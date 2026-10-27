package sound

// Loop is a clip played over and over, as long as the mixer runs, at a
// gain that eases to the one it is faded to: silent until faded in.
type Loop struct {
	m            *Mixer
	clip         *Clip
	pos          float64
	gain, target float64
}

// Loop starts looping c, silent.
func (m *Mixer) Loop(c *Clip) *Loop {
	l := &Loop{m: m, clip: c}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loops = append(m.loops, l)
	return l
}
