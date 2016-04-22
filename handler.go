package scaffold

import (
	"net/http"

	"golang.org/x/net/context"
)

// Handler is a context aware handler
type Handler interface {
	CtxServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

// HandlerFunc wraps a handler letting it implement Handler
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// CtxServeHTTP implements Handler.CtxServeHTTP
func (h HandlerFunc) CtxServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h(ctx, w, r)
}

// HandlerList is an array of handlers
type HandlerList []Handler

// CtxServeHTTP implements Handler.CtxServeHTTP
func (h HandlerList) CtxServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	for _, handler := range h {
		handler.CtxServeHTTP(ctx, w, r)
	}
}
