package ctxhelper_test

import (
	"context"
	"fmt"

	"github.com/tech10/ctxhelper"
)

// Example provides a full example for the most common pattern using H.
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

// ExampleH_IsDone demonstrates IsDone to check context cancellation.
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

// ExampleH_IsQuit demonstrates IsQuit to check if H has been quit.
// A call to Quit will prevent function execution and terminate the waiting goroutines.
func ExampleH_IsQuit() {
	h := ctxhelper.New(context.Background())
	h.OnDone(func() {
		fmt.Println("This function should never be executed.")
	})
	fmt.Println("Quit:", h.IsQuit())
	h.QuitAndWait()
	// h.Quit() then h.Wait() could also be called here.
	fmt.Println("Quit:", h.IsQuit())
	// Output:
	// Quit: false
	// Quit: true
}
