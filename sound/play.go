package sound

// Play starts playing c as v says. A mixer with as many voices playing as
// it may have cuts the oldest.
func (m *Mixer) Play(c *Clip, v Voice) {
	if c == nil || len(c.Samples) == 0 || c.Rate <= 0 {
		return
	}
	pitch := v.Pitch
	if pitch <= 0 {
		pitch = 1
	}
	p := &playing{clip: c, step: pitch}
	p.left, p.right = gains(v.Gain, v.Pan)
	m.mu.Lock()
	defer m.mu.Unlock()
	if max := m.maxVoices(); len(m.playing) >= max {
		m.playing = append(m.playing[:0], m.playing[len(m.playing)-max+1:]...)
	}
	m.playing = append(m.playing, p)
}

func (m *Mixer) maxVoices() int {
	if m.MaxVoices > 0 {
		return m.MaxVoices
	}
	return 32
}
