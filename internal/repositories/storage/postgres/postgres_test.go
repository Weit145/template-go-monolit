package postgres

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPingWithRetryRecovers(t *testing.T) {
	calls := 0
	var contexts []context.Context
	err := pingWithRetry(context.Background(), func(ctx context.Context) error {
		contexts = append(contexts, ctx)
		if _, ok := ctx.Deadline(); !ok {
			t.Error("ping attempt has no deadline")
		}
		calls++
		if calls < 3 {
			return errors.New("database is starting")
		}
		return nil
	}, 5, time.Second, 0)
	if err != nil || calls != 3 {
		t.Fatalf("err=%v, attempts=%d; want success on third attempt", err, calls)
	}
	for _, ctx := range contexts {
		if ctx.Err() == nil {
			t.Error("attempt context was not canceled after ping")
		}
	}
}

func TestPingWithRetryExhausted(t *testing.T) {
	failure := errors.New("database unavailable")
	calls := 0
	err := pingWithRetry(context.Background(), func(context.Context) error {
		calls++
		return failure
	}, 3, time.Second, 0)
	if !errors.Is(err, failure) || calls != 3 {
		t.Fatalf("err=%v, attempts=%d; want last error after three attempts", err, calls)
	}
}

func TestPingWithRetryAttemptTimeout(t *testing.T) {
	calls := 0
	err := pingWithRetry(context.Background(), func(ctx context.Context) error {
		calls++
		<-ctx.Done()
		return ctx.Err()
	}, 2, time.Millisecond, 0)
	if !errors.Is(err, context.DeadlineExceeded) || calls != 2 {
		t.Fatalf("err=%v, attempts=%d; want two timed out attempts", err, calls)
	}
}

func TestPingWithRetryCanceledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := pingWithRetry(ctx, func(context.Context) error {
		t.Error("ping called with canceled parent context")
		return nil
	}, 3, time.Second, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v, want context cancellation", err)
	}
}

func TestPingWithRetryCanceledDuringPause(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- pingWithRetry(ctx, func(context.Context) error {
			close(called)
			return errors.New("database unavailable")
		}, 3, time.Second, time.Hour)
	}()
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("first ping was not called")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("err=%v, want context cancellation", err)
		}
	case <-time.After(time.Second):
		t.Fatal("cancellation did not interrupt retry pause")
	}
}

func TestPingWithRetryDoesNotPauseAfterLastAttempt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	failure := errors.New("database unavailable")
	err := pingWithRetry(ctx, func(context.Context) error { return failure }, 1, time.Second, time.Hour)
	if !errors.Is(err, failure) {
		t.Fatalf("err=%v, want ping error without waiting after the last attempt", err)
	}
}
