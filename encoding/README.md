# Scaffold Encoding

This package adds encoding and decoding to scaffold for apis.
It uses the Content-Type and Accept headers to work out what data type to return.

## Example:
This example uses the handler builder to auto decode the request and encode the response:
```go
package main

import (
    "net/http"

    "github.com/ThatsMrTalbot/scaffold"
    "github.com/ThatsMrTalbot/scaffold/encoding"
    "github.com/ThatsMrTalbot/scaffold/errors"
	"golang.org/x/net/context"
)

type ExampleRequest struct {
    SomeValue string `json:"someValue" xml:"SomeValue"`
    OtherValue int `json:"otherValue" xml:"OtherValue"`
}

type ExampleResponse struct {
    SomeValue string `json:"someValue" xml:"SomeValue"`
    OtherValue int `json:"otherValue" xml:"OtherValue"`
}

// This handler will receive the decoded request in `request`
// The response will be encoded and sent to the client
// Eg `{"someValue":"Some value","otherValue":5}`
func ExampleHandler1(request ExampleRequest) (ExampleResponse, error) {
    return ExampleResponse{
        SomeValue: request.SomeValue,
        OtherValue: request.OtherValue,
    }, nil
}

// This handler will receive the decoded request in `request`
// The error will be sent to the client, eg `{"error":"Some error"}`
//
// If the handler has context.Context, *http.Request or http.ResponseWriter
// parameters they will be passed in.
func ExampleHandler2(request ExampleRequest, ctx context.Context) (ExampleResponse, error) {
    return ExampleResponse{}, errors.NewErrorStatus(500, "Some error")
}

func main() {
    dispatcher := scaffold.DefaultDispatcher()
    router := scaffold.New(dispatcher)

    // Add handler builder to accept handlers that need decoding/encoding
    router.AddHandlerBuilder(encoding.DefaultHandlerBuilder)

    // Route the example handlers
    router.Handle("/example1", ExampleHandler1)
    router.Handle("/example2", ExampleHandler2)

    http.ListenAndServe(":8080", dispatcher)
}
```