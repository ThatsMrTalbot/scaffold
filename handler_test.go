package scaffold

import (
	"net/http"

	"golang.org/x/net/context"
)

type MockHandler struct {
	Context context.Context
	Request *http.Request
}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (m *MockHandler) CtxServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	m.Context = ctx
	m.Request = r
}
