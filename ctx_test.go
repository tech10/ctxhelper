package ctxhelper

import (
	"context"
	"sync"
	"testing"
)

func TestPanicOnNewWithNilContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on New with nil context")
		}
	}()
	New(nil)
}

func TestCancel(t *testing.T) {
	var mu sync.Mutex
	count := 0
	h := New(context.Background())
	h.OnDone(func() {
		mu.Lock()
		defer mu.Unlock()
		count++
		t.Logf("called %d", count)
	})
	h.OnDone(func() {
		mu.Lock()
		defer mu.Unlock()
		count++
		t.Logf("called %d", count)
	})
	h.Cancel()
	mu.Lock()
	num := count
	mu.Unlock()
	if num != 2 {
		t.Fatalf("expected 2, got %d", num)
	}
}

func TestIsDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	h := New(ctx)
	if h.IsDone() {
		t.Fatal("expected false when context is not canceled.")
	}
	cancel()
	if !h.IsDone() {
		t.Fatal("expected true when context is canceled.")
	}
}

func TestIsNotContext(t *testing.T) {
	h := New(context.Background())
	if _, ok := interface{}(h).(context.Context); ok {
		t.Fatal("H should not be a context.Context interface")
	}
}
