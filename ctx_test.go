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
	h.CancelAndWait()
	mu.Lock()
	num := count
	mu.Unlock()
	if num != 2 {
		t.Fatalf("expected 2, got %d", num)
	}
}

func TestQuit(t *testing.T) {
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
	h.QuitAndWait()
	mu.Lock()
	num := count
	mu.Unlock()
	if num != 0 {
		t.Fatalf("expected 0, got %d", num)
	}
	if h.IsDone() {
		t.Fatal("expected no context cancellation on quit")
	}

	h.CancelAndWait()
	mu.Lock()
	num = count
	mu.Unlock()
	if num != 0 {
		t.Fatalf("expected 0, got %d: functions should not be executed with context cancellation after Quit has been called", num)
	}
}

func TestIsDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	h := New(ctx)
	if h.IsDone() {
		t.Fatal("expected false on IsDone when context is not canceled")
	}
	cancel()
	if !h.IsDone() {
		t.Fatal("expected true on IsDone when context is canceled")
	}
}

func TestIsQuit(t *testing.T) {
	h := New(context.Background())
	if h.IsQuit() {
		t.Fatal("expected false on IsQuit when H has been created")
	}
	h.Quit()
	if !h.IsQuit() {
		t.Fatal("expected true on IsQuit when H has been quit")
	}
}

func TestQuitCallMultiple(t *testing.T) {
	var wg sync.WaitGroup
	h := New(context.Background())
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.Quit()
		}()
	}
	wg.Wait()
	if !h.IsQuit() {
		t.Fatal("H should be quit on multiple calls to quit")
	}
	t.Log("no panic occurred")
}
