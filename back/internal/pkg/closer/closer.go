package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Closer interface {
	Close(ctx context.Context) error
	Add(f CloseFunc)
}

type closer struct {
	m          sync.Mutex
	closeFuncs []CloseFunc
}

type CloseFunc func() error

func NewCloser() Closer {
	return &closer{
		m:          sync.Mutex{},
		closeFuncs: make([]CloseFunc, 0),
	}
}

func (c *closer) Close(ctx context.Context) error {
	c.m.Lock()
	defer c.m.Unlock()

	var (
		msgs     = make([]string, 0, len(c.closeFuncs))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.closeFuncs {
			if err := f(); err != nil {
				msgs = append(msgs, err.Error())
			}
		}
		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("error shutdown canceled: %w", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf("error shutdown error(s): \n%s",
			strings.Join(msgs, "\n"))
	}

	return nil
}

func (c *closer) Add(f CloseFunc) {
	c.m.Lock()
	defer c.m.Unlock()
	c.closeFuncs = append(c.closeFuncs, f)
}
