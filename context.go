package scaffold

import (
	"net/http"

	"golang.org/x/net/context"
)

type paramContext struct {
	context.Context
	params map[string]Param
}

func newParamContext(ctx context.Context, r *http.Request, vars map[int]string) *paramContext {
	params := make(map[string]Param)
	for i, p := range vars {
		k, v := paramName(p), ""
		ctx, v, _ = URLPart(ctx, r, i)
		params[k] = Param(v)
	}
	return &paramContext{
		Context: ctx,
		params:  params,
	}
}

func (p *paramContext) Value(key interface{}) interface{} {
	if p.params != nil {
		if k, ok := key.(string); ok {
			if param, ok := p.params[k]; ok {
				return param
			}
		}
	}

	return p.Context.Value(key)
}

func wrapMiddleware(vars map[int]string, middleware Middleware) Middleware {
	return Middleware(func(next Handler) Handler {
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			p := newParamContext(ctx, r, vars)
			teardown := HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				p.params = nil
				next.CtxServeHTTP(ctx, w, r)
			})
			middleware(teardown).CtxServeHTTP(p, w, r)
		})
	})
}

func wrapHandler(vars map[int]string, handler Handler) Handler {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		p := newParamContext(ctx, r, vars)
		handler.CtxServeHTTP(p, w, r)
		p.params = nil
	})
}
