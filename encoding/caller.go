package encoding

import (
	"net/http"
	"reflect"

	"github.com/ThatsMrTalbot/scaffold/errors"
	"golang.org/x/net/context"
)

type caller struct {
	e      *Encoder
	params []func(context.Context, http.ResponseWriter, *http.Request) (reflect.Value, error)
	caller reflect.Value
}

func (c *caller) CtxServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	responder := c.e.Responder(r)
	params := make([]reflect.Value, len(c.params))
	err := error(nil)

	for i, p := range c.params {
		params[i], err = p(ctx, w, r)
		if err != nil {
			responder.Respond(500, w, errorObj{Error: err.Error()})
			return
		}
	}

	results := c.caller.Call(params)
	resp := results[0].Interface()

	if err, ok := results[1].Interface().(error); ok {
		s := 500
		if status, ok := err.(errors.ErrorStatus); ok {
			s = status.Status()
		}
		responder.Respond(s, w, errorObj{Error: err.(error).Error()})
		return
	}

	responder.Respond(200, w, resp)
}
