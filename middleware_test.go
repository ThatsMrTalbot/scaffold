package scaffold

import (
	"net/http"

	"golang.org/x/net/context"
)

type MockMiddleware struct {
	Called  bool
	Context context.Context
	Request *http.Request
}

func NewMockMiddleware() *MockMiddleware {
	return &MockMiddleware{}
}

func (m *MockMiddleware) Middleware() Middleware {
	return Middleware(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			m.Called = true
			m.Context = ctx
			m.Request = r
			next.CtxServeHTTP(ctx, w, r)
		})
	})
}
