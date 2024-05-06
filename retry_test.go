// Copyright (c) 2024 BVK Chaitanya

package ctxutil

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	n := 0
	f := func(ctx context.Context) error {
		if n < 10 {
			t.Logf("n == %d: needs retry (%s)", n, time.Now())
			n++
			return errors.New("retry")
		}
		return nil
	}
	intervals := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		4 * time.Millisecond,
		8 * time.Millisecond,
		16 * time.Millisecond,
		32 * time.Millisecond,
	}
	Retry(context.Background(), f, intervals...)
	if n != 10 {
		t.Fatalf("want 10, got %d", n)
	}
}
