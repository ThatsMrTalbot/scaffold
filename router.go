package scaffold

import "path"

// Router is a HTTP router
type Router struct {
	route      Route
	dispatcher Dispatcher
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

// Host lets you specify the host(s) the route is avaiable for
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
func (r *Router) Handle(pattern string, handlers ...Handler) *Router {
	c := r.pattern(pattern)
	r.dispatcher.Handle(c, handlers...)
	return r.clone(c)
}

// Options handles OPTIONS methods with a given pattern
func (r *Router) Options(pattern string, handlers ...Handler) *Router {
	c := r.pattern(pattern)
	c.Method = "OPTIONS"
	r.dispatcher.Handle(c, handlers...)
	return r.clone(c)
}

// Get handles GET methods with a given pattern
func (r *Router) Get(pattern string, handlers ...Handler) *Router {
	c := r.pattern(pattern)
	c.Method = "GET"
	r.dispatcher.Handle(c, handlers...)
	return r.clone(c)
}

// Head handles HEAD methods with a given pattern
func (r *Router) Head(pattern string, handlers ...Handler) *Router {
	c := r.pattern(pattern)
	c.Method = "HEAD"
	r.dispatcher.Handle(c, handlers...)
	return r.clone(c)
}

// Post handles POST methods with a given pattern
func (r *Router) Post(pattern string, handlers ...Handler) *Router {
	c := r.pattern(pattern)
	c.Method = "POST"
	r.dispatcher.Handle(c, handlers...)
	return r.clone(c)
}

// Put handles PUT methods with a given pattern
func (r *Router) Put(pattern string, handlers ...Handler) *Router {
	c := r.pattern(pattern)
	c.Method = "PUT"
	r.dispatcher.Handle(c, handlers...)
	return r.clone(c)
}

// Delete handles DELETE methods with a given pattern
func (r *Router) Delete(pattern string, handlers ...Handler) *Router {
	c := r.pattern(pattern)
	c.Method = "DELETE"
	r.dispatcher.Handle(c, handlers...)
	return r.clone(c)
}

// Use attaches middleware to a route
func (r *Router) Use(middleware ...Middleware) {
	r.dispatcher.Middleware(r.route, middleware...)
}

// UseFunc attaches middlware to a route
func (r *Router) UseFunc(middleware ...func(next Handler) Handler) {
	m1 := make([]Middleware, len(middleware))
	for i, m2 := range middleware {
		m1[i] = Middleware(m2)
	}
	r.Use(m1...)
}

// NotFound specifys a not found handler for a route
func (r *Router) NotFound(handler Handler) {
	r.dispatcher.NotFoundHandler(r.route, handler)
}

func (r *Router) clone(route Route) *Router {
	return &Router{
		dispatcher: r.dispatcher,
		route:      route,
	}
}

func (r *Router) pattern(pattern string) Route {
	c := r.route.clone()
	c.Pattern = path.Join(r.route.Pattern, pattern)
	return c
}
