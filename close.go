package errorx

import (
	"errors"
	"sync"
	"sync/atomic"
)

type CloseErr struct {
	mu     sync.RWMutex
	err    error
	done   chan struct{}
	closed atomic.Bool

	// ErrClosed is returned by repeated Close calls after a successful close;
	// defaults to ErrClosed.
	ErrClosed error
}

var ErrClosed = errors.New("closed")

// Close invokes fn to release resources and keeps the first error it returns.
func (c *CloseErr) Close(fn func() (errs []error)) error {
	if c.closed.Load() {
		return c.Error()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed.Load() {
		return c.err
	}
	c.closed.Store(true)

	if c.done != nil {
		close(c.done)
	}
	if fn != nil {
		for _, e := range fn() {
			if e != nil {
				c.err = e
				break // keep only the first error
			}
		}
	}

	if c.err != nil {
		return c.err
	} else {
		if c.ErrClosed == nil {
			c.ErrClosed = ErrClosed
		}
		c.err = c.ErrClosed
		return nil // first return nil
	}
}

// Error returns the error recorded by Close, or nil if Close was never called.
func (c *CloseErr) Error() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.err != nil && c.err == c.ErrClosed {
		return WithT(c.err, newStack())
	} else {
		return c.err
	}
}

func (c *CloseErr) Closed() bool { return c.closed.Load() }

// Done returns a channel that is closed once Close is called.
func (c *CloseErr) Done() <-chan struct{} {
	c.mu.RLock()
	done := c.done
	c.mu.RUnlock()
	if done != nil {
		return done
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.done != nil {
		return c.done
	}
	c.done = make(chan struct{}, 1)
	if c.err != nil {
		close(c.done)
	}
	return c.done
}
