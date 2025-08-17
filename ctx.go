// Package ctxhelper provides helpers for executing functions
// once a context is canceled, with proper synchronization.
package ctxhelper

import (
	"context"
	"sync"
	"time"
)

// H provides context helpers for the context it was created with.
// H must be created with New.
type H struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates a new helper with the provided context and sync.WaitGroup.
// If ctx is nil, a runtime panic will be produced.
func New(ctx context.Context) *H {
	if ctx == nil {
		panic("ctxhelper: nil context not permitted")
	}
	h := &H{}
	h.ctx, h.cancel = context.WithCancel(ctx)
	return h
}

// OnDone waits for the context to be canceled, then executes fn.
// It increments the WaitGroup before waiting, and decrements it after fn finishes.
func (h *H) OnDone(fn func()) {
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

// Cancel cancels H and waits for any functions to complete their execution.
func (h *H) Cancel() {
	h.cancel()
	h.Wait()
}

// Wait waits for all functions to be called on context cancellation.
func (h *H) Wait() {
	h.wg.Wait()
}

// Context returns the underlying context managed by H.
func (h *H) Context() context.Context {
	return h.ctx
}

// Close cancels and waits, making H usable as an io.Closer.
func (h *H) Close() error {
	h.Cancel()
	return nil
}

// Deadline returns ctx.Deadline values.
func (h *H) Deadline() (time.Time, bool) {
	return h.ctx.Deadline()
}

// Done returns the done channel associated with ctx.
func (h *H) Done() <-chan struct{} {
	return h.ctx.Done()
}
