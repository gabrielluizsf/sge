//go:build gots_overlay

// Overlay for io/pipe.go: a pipe whose writes go into a buffer the reader
// drains, since without goroutines a writer cannot wait for a reader.

package io

type pipe struct {
	buf    []byte
	rerr   onceError
	werr   onceError
	closed bool
}

func (p *pipe) read(b []byte) (n int, err error) {
	if err := p.rerr.Load(); err != nil {
		return 0, ErrClosedPipe
	}
	if len(p.buf) == 0 {
		if p.closed {
			return 0, p.readCloseError()
		}
		return 0, nil
	}
	n = copy(b, p.buf)
	p.buf = p.buf[n:]
	return n, nil
}

func (p *pipe) closeRead(err error) error {
	if err == nil {
		err = ErrClosedPipe
	}
	p.rerr.Store(err)
	return nil
}

func (p *pipe) write(b []byte) (n int, err error) {
	if p.closed {
		return 0, p.writeCloseError()
	}
	if err := p.rerr.Load(); err != nil {
		return 0, err
	}
	p.buf = append(p.buf, b...)
	return len(b), nil
}

func (p *pipe) closeWrite(err error) error {
	if err == nil {
		err = EOF
	}
	p.werr.Store(err)
	p.closed = true
	return nil
}

func (p *pipe) readCloseError() error {
	if werr := p.werr.Load(); werr != nil {
		return werr
	}
	return ErrClosedPipe
}

func (p *pipe) writeCloseError() error {
	if rerr := p.rerr.Load(); rerr != nil {
		return rerr
	}
	return ErrClosedPipe
}

func Pipe() (*PipeReader, *PipeWriter) {
	pw := &PipeWriter{r: PipeReader{pipe: pipe{}}}
	return &pw.r, pw
}
