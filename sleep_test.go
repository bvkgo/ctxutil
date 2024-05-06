// Copyright (c) 2024 BVK Chaitanya

package ctxutil

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSleepTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sleepCtx, sleepCancel := context.WithTimeoutCause(ctx, time.Millisecond, errTimedout)
	defer sleepCancel()

	cancel()
	<-sleepCtx.Done()

	if err := context.Cause(sleepCtx); errors.Is(err, errTimedout) {
		t.Fatalf("want %v, got %v", context.Canceled, err)
	}
}
