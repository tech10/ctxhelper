// Package ctxhelper provides helpers for executing functions
// once a context is canceled, with proper synchronization.
package ctxhelper

import (
	"context"
	"sync"
)

// H provides context helper functions for the context it was created with.
// H must be created with New.
type H struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates H with a child context from ctx.
// If ctx is nil, a runtime panic will be produced.
func New(ctx context.Context) *H {
	if ctx == nil {
		panic("ctxhelper: nil context not permitted")
	}
	h := &H{}
	h.ctx, h.cancel = context.WithCancel(ctx)
	return h
}

// OnDone waits for ctx to be canceled, then executes fn.
// It increments the WaitGroup before waiting, and decrements it after fn finishes.
// This can be used across multiple goroutines and called multiple times.
// If ctx is already canceled, this will be a no op.
//
// Each call to OnDone will wait for context cancellation and function execution in its own goroutine.
// OnDone is a non-blocking call.
//
// fn must not panic. Any panic recovery is up to the caller of OnDone to implement.
//
// When ctx is canceled, fn can be executed as many times as OnDone has been called,
// but each fn is not executed in any predetermined order.
func (h *H) OnDone(fn func()) {
	if h.IsDone() {
		return
	}
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		<-h.ctx.Done() // wait for cancellation
		fn()
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

// CancelAndWait cancels ctx and waits for any functions to complete their execution.
func (h *H) CancelAndWait() {
	h.cancel()
	h.Wait()
}

// Cancel cancels ctx but does not wait for any functions to complete their execution.
func (h *H) Cancel() {
	h.cancel()
}

// Wait waits for all functions to complete execution on context cancellation.
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

// Done returns the done channel associated with ctx.
func (h *H) Done() <-chan struct{} {
	return h.ctx.Done()
}
