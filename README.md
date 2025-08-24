This package is a small, but convenient wrapper around contexts. With it, one of the more common uses of contexts with a sync.WaitGroup can be more easily dealt with and condensed. For me, the more I've dealt with this pattern, the more useful it is to use with this package.

## With ctxhelper

```go
package main

import (
	"context"
	"fmt"

	"github.com/tech10/ctxhelper"
)

func main() {
	h := ctxhelper.New(context.Background())
	h.OnDone(func() {
		fmt.Println("Done.")
	})
	h.CancelAndWait()
	fmt.Println("All function calls complete.")
	// Output:
	// Done.
	// All function calls complete.
}
```

## Without ctxhelper

```go
package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		fmt.Println("Done.")
	}()
	cancel()
	wg.Wait()
	fmt.Println("All function calls complete.")
	// Output:
	// Done.
	// All function calls complete.
}
```

# OnDoneWithCancel

This operates exactly as OnDone does in the example above, but it will return a context.CancelFunc for individually canceling that specific goroutine. This will result in that specific function not executing when the context is canceled. See package documentation for a full example of how this works.

# Future additions

There may be other helper functions added in the future, though they will not all be documented here. See the [docs for this package](https://pkg.go.dev/github.com/tech10/ctxhelper) for full documentation.

# Contributing

As always, gofmt, and create a PR.