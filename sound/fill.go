package sound

import "math"

// Fill mixes the next frames of the stream into block, interleaved left
// and right, overwriting it: what plays, and silence where nothing does.
// Every sample is kept from -1 to 1, and one that is not a number is
// silence, so no sound can make the speaker click.
func (m *Mixer) Fill(block []float32) {
	clear(block)
	m.mu.Lock()
	rate := float64(m.rate)
	kept := m.playing[:0]
	for _, p := range m.playing {
		if p.mix(block, rate) {
			kept = append(kept, p)
		}
	}
	clear(m.playing[len(kept):])
	m.playing = kept
	for _, l := range m.loops {
		l.mix(block, rate)
	}
	m.mu.Unlock()
	for i, s := range block {
		switch {
		case s != s:
			block[i] = 0
		case s > 1:
			block[i] = 1
		case s < -1:
			block[i] = -1
		}
	}
}

// mix adds p's next frames to block, at a stream of rate, and reports
// whether there is more of it to play.
func (p *playing) mix(block []float32, rate float64) bool {
	s := p.clip.Samples
	step := p.step * float64(p.clip.Rate) / rate
	for i := 0; i+1 < len(block); i += 2 {
		if p.pos >= float64(len(s)-1) {
			return false
		}
		v := sampleAt(s, p.pos)
		block[i] += v * p.left
		block[i+1] += v * p.right
		p.pos += step
	}
	return true
}

// mix adds the loop's next frames to block, at a stream of rate, easing
// its gain to the one it is faded to.
func (l *Loop) mix(block []float32, rate float64) {
	s := l.clip.Samples
	if len(s) < 2 {
		return
	}
	target := l.target
	ease := 1 - math.Exp(-1/(fadeTime*rate))
	step := float64(l.clip.Rate) / rate
	n := float64(len(s))
	for i := 0; i+1 < len(block); i += 2 {
		l.gain += (target - l.gain) * ease
		if l.gain < 1e-4 && target == 0 {
			l.pos = math.Mod(l.pos+step*float64(len(block)-i)/2, n)
			l.gain = 0
			return
		}
		v := sampleAt(s, l.pos) * float32(l.gain)
		block[i] += v
		block[i+1] += v
		if l.pos += step; l.pos >= n {
			l.pos -= n
		}
	}
}

// sampleAt is s at the place pos between two samples, read between them,
// wrapping to the first past the last.
func sampleAt(s []float32, pos float64) float32 {
	i := int(pos)
	f := float32(pos - float64(i))
	a := s[i]
	b := s[0]
	if i+1 < len(s) {
		b = s[i+1]
	}
	return a + (b-a)*f
}
