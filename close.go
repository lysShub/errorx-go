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

	// 当Close时没有错误, 后续重复调用Close时将返回此错误, 默认值为[ErrClosed]
	ErrClosed error
}

var ErrClosed = errors.New("closed")

// Close 关闭, 在回调函数中关闭各个资源, 并把错误追加到errs中
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
				break // 只收集最初的错误
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

// Error 返回Close时的err, 没有Close时返回nil
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

// Done 关闭通知
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
