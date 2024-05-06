// Copyright (c) 2024 BVK Chaitanya

package ctxutil

import (
	"context"
	"os"
	"sync"
)

// A Group represents a collection of goroutines sharing a cancellation
// context, so that when the group is closed all goroutines sharing the group
// context are canceled and cleaned up properly.
type Group struct {
	closeCtx context.Context

	mu sync.Mutex

	cancelFunc context.CancelCauseFunc

	wg sync.WaitGroup
}

// NewGroup creates a group derived from the input context.
func NewGroup(ctx context.Context) *Group {
	gctx, gcancel := context.WithCancelCause(ctx)
	g := &Group{
		closeCtx:   gctx,
		cancelFunc: gcancel,
	}
	return g
}

// Close cancels the group context with os.ErrClosed and waits for all
// goroutines to finish.
func (g *Group) Close() {
	g.Cancel(os.ErrClosed)
	g.wg.Wait()
}

// Wait waits for all goroutines to finish.
func (g *Group) Wait() {
	g.wg.Wait()
}

// Cancel canels the group context with a specific error. This method does not
// wait for the goroutines to finish.
func (g *Group) Cancel(cause error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.cancelFunc(cause)
}

// GroupContext returns the group context.
func (g *Group) GroupContext() context.Context {
	return g.closeCtx
}

// Go runs the input function with a new goroutine under the control of the
// group context. This method can be run concurrently with Cancel/Close
// methods, in which case, the input function may not run if group context is
// canceled before acquiring the lock.
func (g *Group) Go(f func(ctx context.Context)) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.closeCtx.Err() == nil {
		g.wg.Add(1)
		go func() {
			f(g.closeCtx)
			g.wg.Done()
		}()
	}
}
