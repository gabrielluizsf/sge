//go:build gots_overlay

// Overlay for handshakeContext: Go's holds handshakeMutex while the
// handshake runs, which under gots waits on node:tls, and a mutex of gots
// cannot be waited for. A second caller waits on a channel instead, which
// can be, until the first is done.

package tls

import (
	"context"
	"errors"
)

func (c *Conn) handshakeContext(ctx context.Context) (ret error) {
	if c.isHandshakeComplete.Load() {
		return nil
	}
	p := c.pump()
	if p.handshaking != nil {
		select {
		case <-p.handshaking:
			return c.handshakeErr
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	p.handshaking = make(chan struct{})
	defer close(p.handshaking)

	handshakeCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if ctx.Done() != nil {
		stop := context.AfterFunc(ctx, func() {
			_ = c.conn.Close()
		})
		defer func() {
			if !stop() {
				ret = ctx.Err()
			}
		}()
	}

	c.handshakeErr = c.handshakeFn(handshakeCtx)
	if c.handshakeErr == nil {
		c.handshakes++
		if !c.isHandshakeComplete.Load() {
			c.handshakeErr = errors.New("tls: internal error: handshake should have had a result")
		}
	}
	return c.handshakeErr
}
