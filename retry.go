// Copyright (c) 2024 BVK Chaitanya

package ctxutil

import (
	"context"
	"time"
)

// Retry runs the input function till it succeeds with given retry intervals or
// the input context is canceled.
//
// If the retry intervals are omitted a fixed one second interval is
// assumed. When the last retry interval is reached, it is continued for all
// future retries, which enables users to select exponential backoff like
// intervals with a maximum limit.
//
// Returns nil if the input function is successful or the last non-nil error
// from the function after the context has expired.
func Retry(ctx context.Context, f func(context.Context) error, intervals ...time.Duration) (err error) {
	if len(intervals) == 0 {
		intervals = []time.Duration{time.Second}
	}

	i := 0
	for err = f(ctx); err != nil && context.Cause(ctx) == nil; err = f(ctx) {
		interval := intervals[i%len(intervals)]
		SleepTimeout(ctx, interval)
		if i < len(intervals)-1 {
			i++
		}
	}
	return
}
