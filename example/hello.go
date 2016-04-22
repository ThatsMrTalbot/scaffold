package main

import (
	"fmt"
	"net/http"

	"github.com/ThatsMrTalbot/scaffold"
	"golang.org/x/net/context"
)

// Hello is an example scaffold platform
type Hello struct {
}

// World is a handler
func (h *Hello) World(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// Name is a handler
func (h *Hello) Name(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	name := scaffold.GetParam(ctx, "name")

	msg := fmt.Sprintf("Hello %s!", name)
	w.Write([]byte(msg))
}

// Routes implements scaffold.Platform.Routes
func (h *Hello) Routes(router *scaffold.Router) {
	router.Get("", scaffold.HandlerFunc(h.World))
	router.Get(":name", scaffold.HandlerFunc(h.Name))
}
