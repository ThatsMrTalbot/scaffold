package errors

import (
	"fmt"
	"net/http"

	"github.com/ThatsMrTalbot/scaffold"

	"golang.org/x/net/context"
)

// AllStatusCodes is a nice way of saying add for all status codes
const AllStatusCodes = 0

// DefaultErrorHandler is the default error handler
var DefaultErrorHandler = &defaultErrorHandler{}

// ErrorHandler is an error handler
type ErrorHandler interface {
	ServeErrorPage(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error)
}

type defaultErrorHandler struct{}

func (*defaultErrorHandler) ServeErrorPage(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error) {
	http.Error(w, err.Error(), status)
}

// GetErrorHandler gets the error handler from the context or returns the default
func GetErrorHandler(ctx context.Context, status int) ErrorHandler {
	key := fmt.Sprintf("error_handler_%d", status)
	if h, ok := ctx.Value(key).(ErrorHandler); ok {
		return h
	}

	key = fmt.Sprintf("error_handler_%d", AllStatusCodes)
	if h, ok := ctx.Value(key).(ErrorHandler); ok {
		return h
	}

	return DefaultErrorHandler
}

// SetErrorHandler returns Middleware that can be used to set the error handler
func SetErrorHandler(status int, handler ErrorHandler) scaffold.Middleware {
	return scaffold.Middleware(func(next scaffold.Handler) scaffold.Handler {
		return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			key := fmt.Sprintf("error_handler_%d", status)
			ctx = context.WithValue(ctx, key, handler)

			next.CtxServeHTTP(ctx, w, r)
		})
	})
}
