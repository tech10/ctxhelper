// Package ctxhelper provides helpers for executing functions
// once a context is canceled, with WaitGroup synchronization.
package ctxhelper

import (
	"context"
	"sync"
)

// H provides context helper functions for the context it was created with.
// H must be created with New.
type H struct {
	ctx      context.Context
	cancel   context.CancelFunc
	quitch   chan struct{}
	quitOnce sync.Once
	wg       sync.WaitGroup
}

// New creates H with a child context from ctx.
// If ctx is nil, a runtime panic will be produced.
func New(ctx context.Context) *H {
	if ctx == nil {
		panic("ctxhelper: nil context not permitted")
	}
	h := &H{
		quitch: make(chan struct{}),
	}
	h.ctx, h.cancel = context.WithCancel(ctx)
	return h
}

// OnDone waits for ctx to be canceled, then executes fn.
// It increments the WaitGroup before waiting, and decrements it after fn finishes.
// If Quit is called on H, fn will not be executed, but the internal WaitGroup will still be incremented and decremented as needed.
//
// OnDone can be used across multiple goroutines and called multiple times.
// If ctx is already canceled or H is terminated via Quit,
// calling OnDone will be a no op.
//
// Each call to OnDone will wait for ctx cancellation and function execution, or a call to Quit, in its own goroutine.
// OnDone is a non-blocking call.
//
// fn must not panic. Any panic recovery is up to the caller of OnDone to implement.
//
// When ctx is canceled, fn will be executed as many times as OnDone has been called,
// but each fn is not executed in any predetermined order.
//
// Once OnDone is called, any functions being executed on ctx cancellation cannot be removed.
// Before any functions are executed via context cancellation,
// you can quit all function termination by calling Quit.
func (h *H) OnDone(fn func()) {
	if h.IsDone() || h.IsQuit() {
		return
	}
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		select {
		case <-h.ctx.Done(): // wait for context cancellation
			fn()
		case <-h.quitch:
			return
		}
	}()
}

// IsDone returns true if the context has been canceled, false if not.
func (h *H) IsDone() bool {
	select {
	case <-h.ctx.Done():
		return true
	default:
		return false
	}
}

// IsQuit returns true if H has been quit, false if not.
func (h *H) IsQuit() bool {
	select {
	case <-h.quitch:
		return true
	default:
		return false
	}
}

// Quit quits function execution.
// This works like Cancel, but Quit will ensure the functions pending execution will not be called.
func (h *H) Quit() {
	h.quitOnce.Do(func() {
		close(h.quitch)
	})
}

// QuitAndWait quits all function execution and waits for goroutine termination.
func (h *H) QuitAndWait() {
	h.Quit()
	h.Wait()
}

// Cancel cancels ctx but does not wait for any functions to complete their execution.
func (h *H) Cancel() {
	h.cancel()
}

// CancelAndWait cancels ctx and waits for any functions to complete their execution.
func (h *H) CancelAndWait() {
	h.Cancel()
	h.Wait()
}

// Wait waits for all functions to complete execution on ctx cancellation, or waits for all pending goroutines to terminate on Quit.
func (h *H) Wait() {
	h.wg.Wait()
}

// Context returns the underlying context within H.
func (h *H) Context() context.Context {
	return h.ctx
}

// Close cancels ctx and waits for function execution, making H usable as an io.Closer.
func (h *H) Close() error {
	h.CancelAndWait()
	return nil
}
