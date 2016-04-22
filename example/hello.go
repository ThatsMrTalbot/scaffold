package main

import (
	"fmt"
	"net/http"

	"github.com/ThatsMrTalbot/scaffold"
	"golang.org/x/net/context"
)

type Hello struct {
}

func (h *Hello) World(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (h *Hello) Name(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	name := scaffold.GetParam(ctx, "name")

	msg := fmt.Sprintf("Hello %s!", name)
	w.Write([]byte(msg))
}

func (h *Hello) Routes(router *scaffold.Router) {
	router.Get("", scaffold.HandlerFunc(h.World))
	router.Get(":name", scaffold.HandlerFunc(h.Name))
}
