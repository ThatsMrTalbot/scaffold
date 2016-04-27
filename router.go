package scaffold

import (
	"net/http"
	"path"

	"golang.org/x/net/context"
)

// Router is a HTTP router
type Router struct {
	route      Route
	dispatcher Dispatcher
	builders   []*func(interface{}) (Handler, error)
}

// New creates new router
func New(d Dispatcher) *Router {
	return &Router{
		dispatcher: d,
	}
}

// Scaffold creates a router and passes it to a platorm
func Scaffold(d Dispatcher, platform Platform) Handler {
	router := New(d)
	platform.Routes(router)
	return router.dispatcher
}

// Host lets you specify the host(s) the route is available for
func (r *Router) Host(hosts ...string) *Router {
	c := r.route.clone()
	c.Hosts = hosts
	return r.clone(c)
}

// Route returns the subrouter for a pettern
func (r *Router) Route(pattern string) *Router {
	c := r.pattern(pattern)
	return r.clone(c)
}

// Group calls the specified function with the subrouter for the given pattern
func (r *Router) Group(pattern string, group func(*Router)) {
	c := r.pattern(pattern)
	group(r.clone(c))
}

// Platform routes the platform object to the given pattern
func (r *Router) Platform(pattern string, platform Platform) {
	c := r.pattern(pattern)
	platform.Routes(r.clone(c))
}

// Handle all methods with a given pattern
func (r *Router) Handle(pattern string, handlers ...interface{}) *Router {
	c := r.pattern(pattern)
	clone := r.clone(c)
	clone.handle(handlers)
	return clone
}

// Options handles OPTIONS methods with a given pattern
func (r *Router) Options(pattern string, handlers ...interface{}) *Router {
	c := r.pattern(pattern)
	c.Method = "OPTIONS"
	clone := r.clone(c)
	clone.handle(handlers)
	return clone
}

// Get handles GET methods with a given pattern
func (r *Router) Get(pattern string, handlers ...interface{}) *Router {
	c := r.pattern(pattern)
	c.Method = "GET"
	clone := r.clone(c)
	clone.handle(handlers)
	return clone
}

// Head handles HEAD methods with a given pattern
func (r *Router) Head(pattern string, handlers ...interface{}) *Router {
	c := r.pattern(pattern)
	c.Method = "HEAD"
	clone := r.clone(c)
	clone.handle(handlers)
	return clone
}

// Post handles POST methods with a given pattern
func (r *Router) Post(pattern string, handlers ...interface{}) *Router {
	c := r.pattern(pattern)
	c.Method = "POST"
	clone := r.clone(c)
	clone.handle(handlers)
	return clone
}

// Put handles PUT methods with a given pattern
func (r *Router) Put(pattern string, handlers ...interface{}) *Router {
	c := r.pattern(pattern)
	c.Method = "PUT"
	clone := r.clone(c)
	clone.handle(handlers)
	return clone
}

// Delete handles DELETE methods with a given pattern
func (r *Router) Delete(pattern string, handlers ...interface{}) *Router {
	c := r.pattern(pattern)
	c.Method = "DELETE"
	clone := r.clone(c)
	clone.handle(handlers)
	return clone
}

// Use attaches middleware to a route
func (r *Router) Use(middleware ...interface{}) {
	r.dispatcher.Middleware(r.route, r.buildMiddlewares(middleware)...)
}

// NotFound specifys a not found handler for a route
func (r *Router) NotFound(i interface{}) {
	handler := r.buildHandler(i)
	r.dispatcher.NotFoundHandler(r.route, handler)
}

func (r *Router) clone(route Route) *Router {
	return &Router{
		dispatcher: r.dispatcher,
		route:      route,
		builders:   r.builders,
	}
}

// AddHandlerBuilder adds a builder to construct handlers
func (r *Router) AddHandlerBuilder(builder func(interface{}) (Handler, error)) {
	r.builders = append(r.builders, &builder)
}

func (r *Router) handle(i []interface{}) {
	handlers := r.buildHandlers(i)
	r.dispatcher.Handle(r.route, handlers...)
}

func (r *Router) buildMiddleware(i interface{}) Middleware {
	switch i.(type) {
	case Middleware:
		return i.(Middleware)
	case func(Handler) Handler:
		return Middleware(i.(func(Handler) Handler))
	case func(http.Handler) http.Handler:
		h := i.(func(http.Handler) http.Handler)
		return Middleware(func(next Handler) Handler {
			return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				n := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.CtxServeHTTP(ctx, w, r)
				})
				h(n).ServeHTTP(w, r)
			})
		})
	}

	panic("Invalid middleware passsed to router")
}

func (r *Router) buildMiddlewares(i []interface{}) []Middleware {
	middleware := make([]Middleware, len(i))
	for j, h := range i {
		middleware[j] = r.buildMiddleware(h)
	}
	return middleware
}

func (r *Router) buildHandler(i interface{}) Handler {
	switch i.(type) {
	case Handler:
		return i.(Handler)
	case func(context.Context, http.ResponseWriter, *http.Request):
		return HandlerFunc(i.(func(context.Context, http.ResponseWriter, *http.Request)))
	case http.Handler:
		h := i.(http.Handler)
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	case func(http.ResponseWriter, *http.Request):
		h := i.(func(http.ResponseWriter, *http.Request))
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			h(w, r)
		})
	}

	for _, b := range r.builders {
		h, err := (*b)(i)
		if err == nil {
			return h
		}
	}

	panic("Invalid handler passsed to router")
}

func (r *Router) buildHandlers(i []interface{}) []Handler {
	handlers := make([]Handler, len(i))
	for j, h := range i {
		handlers[j] = r.buildHandler(h)
	}
	return handlers
}

func (r *Router) pattern(pattern string) Route {
	c := r.route.clone()
	c.Pattern = path.Join(r.route.Pattern, pattern)
	return c
}
