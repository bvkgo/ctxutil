// Copyright (c) 2024 BVK Chaitanya

package ctxutil

import (
	"context"
	"errors"
	"time"
)

var (
	errTimedout = errors.New("timeout expired")
	errDeadline = errors.New("deadline reached")
)

// SleepTimeout blocks the caller for given timeout duration or the input
// context is canceled. Returns true if sleep has returned because timeout has
// expired and false because context was canceled.
func SleepTimeout(ctx context.Context, d time.Duration) bool {
	sctx, scancel := context.WithTimeoutCause(ctx, d, errTimedout)
	<-sctx.Done()
	scancel()
	return errors.Is(context.Cause(sctx), errTimedout)
}

// SleepDeadline blocks the caller till a deadline time has reached or the
// input context is canceled. Returns true if sleep has returned because
// deadline has reached and false because context was canceled.
func SleepDeadline(ctx context.Context, t time.Time) bool {
	sctx, scancel := context.WithDeadlineCause(ctx, t, errDeadline)
	<-sctx.Done()
	scancel()
	return errors.Is(context.Cause(sctx), errDeadline)
}
