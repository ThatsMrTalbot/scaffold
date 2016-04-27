package errors

import (
	"errors"
	"net/http"

	"github.com/ThatsMrTalbot/scaffold"
	"golang.org/x/net/context"
)

// Handler is similar to scaffold.Handler with the difference that an error is
// returned
type Handler interface {
	CtxServeHTTP(context.Context, http.ResponseWriter, *http.Request) error
}

// HandlerFunc is similar to scaffold.HandlerFunc with the difference that an
// error is returned
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// CtxServeHTTP implements Handler.CtxServeHTTP
func (h HandlerFunc) CtxServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return h(ctx, w, r)
}

// HandlerBuilder can be used to create a scaffold.Handler based on a Handler
func HandlerBuilder(i interface{}) (scaffold.Handler, error) {
	switch i.(type) {
	case Handler:
		return build(i.(Handler)), nil
	case func(context.Context, http.ResponseWriter, *http.Request):
		return build(HandlerFunc(i.(func(context.Context, http.ResponseWriter, *http.Request) error))), nil
	}

	return nil, errors.New("Invalid type")
}

func build(handler Handler) scaffold.Handler {
	return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		err := handler.CtxServeHTTP(ctx, w, r)
		if err != nil {
			status := 500
			if s, ok := err.(ErrorStatus); ok {
				status = s.Status()
			}
			GetErrorHandler(ctx, status).ServeErrorPage(ctx, w, r, status, err)
		}
	})
}
