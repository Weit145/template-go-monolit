package closer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
)

type closeFunc struct {
	name string
	fn   func(context.Context) error
}

type Closer struct {
	funcs []closeFunc
	once  sync.Once
	err   error
}

func (c *Closer) Add(name string, fn func(context.Context) error) {
	c.funcs = append(c.funcs, closeFunc{name: name, fn: fn})
}

func (c *Closer) CloseAll(ctx context.Context) error {
	c.once.Do(func() {
		for _, v := range slices.Backward(c.funcs) {
			slog.InfoContext(ctx, "closing resource", "name", v.name)
			if err := v.fn(ctx); err != nil {
				c.err = errors.Join(c.err, fmt.Errorf("close %s: %w", v.name, err))
			}
		}
		c.funcs = nil
	})
	return c.err
}
