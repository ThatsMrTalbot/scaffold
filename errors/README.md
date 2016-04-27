# Scaffold Errors

This package improves error handling in scaffold.

## Examples:
This example uses the handler builder to accept handlers that return errors:
```go
package main

import (
    "net/http"

    "github.com/ThatsMrTalbot/scaffold"
    "github.com/ThatsMrTalbot/scaffold/errors"
	"golang.org/x/net/context"
)

func ExampleHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
    return errors.NewErrorStatus(500, "Oh dear!")
}

func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error) {
    http.Error(w, err.Error(), status)
}

func main() {
    dispatcher := scaffold.DefaultDispatcher()
    router := scaffold.New(dispatcher)

    // Add handler builder to accept handlers that return errors
    router.AddHandlerBuilder(errors.HandlerBuilder)

    // Set the error handler for 500 errors
    router.Use(errors.SetErrorHandlerFunc(500, ErrorHandler))

    // Route the example handler
    router.Handle("/", ExampleHandler)

    http.ListenAndServe(":8080", dispatcher)
}

```

This example shows getting the error handler without using the handler builder:
```go
package main

import (
    "fmt"
    "net/http"

    "github.com/ThatsMrTalbot/scaffold"
    "github.com/ThatsMrTalbot/scaffold/errors"
	"golang.org/x/net/context"
)

func ExampleHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    err := errors.NewErrorStatus(500, "Oh dear!")
    if err != nil {
        GetErrorHandler(ctx, 500).ServeErrorPage(ctx, w, r, 500, err)
        return
    }
}

func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error) {
    http.Error(w, err.Error(), status)
}

func main() {
    dispatcher := scaffold.DefaultDispatcher()
    router := scaffold.New(dispatcher)

    // Set the error handler for all errors
    router.Use(errors.SetErrorHandlerFunc(errors.AllStatusCodes, ErrorHandler))

    // Route the example handler
    router.Handle("/", ExampleHandler)

    http.ListenAndServe(":8080", dispatcher)
}

```