package ctxhelper_test

import (
	"context"
	"fmt"

	"github.com/tech10/ctxhelper"
)

// Example provides a full example on using H.
func Example() {
	h := ctxhelper.New(context.Background())
	h.OnDone(func() {
		fmt.Println("Done.")
	})
	h.Cancel()
	h.Wait()
	// h.CancelAndWait() could be used here instead of the individual calls above.
	fmt.Println("All function calls complete.")
	// Output:
	// Done.
	// All function calls complete.
}

// ExampleH_IsDone checks to see if the context is canceled.
// Returns true if so, false if not.
func ExampleH_IsDone() {
	ctx, cancel := context.WithCancel(context.Background())
	h := ctxhelper.New(ctx)
	fmt.Println("Done:", h.IsDone())
	cancel()
	fmt.Println("Done:", h.IsDone())
	// Output:
	// Done: false
	// Done: true
}
