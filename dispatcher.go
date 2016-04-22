package scaffold

import (
	"net/http"
	"sync"

	"golang.org/x/net/context"
)

// Dispatcher dipatches requests to routes
type Dispatcher interface {
	http.Handler
	Handler

	Handle(route Route, handlers ...Handler)
	Middleware(route Route, middleware ...Middleware)
	NotFoundHandler(route Route, handler Handler)
}

type dispatcher struct {
	lock sync.RWMutex

	*node
	hosts map[string]*node
}

// DefaultDispacher implements dispatcher allowing for basic url params
func DefaultDispacher() Dispatcher {
	return &dispatcher{
		node:  newNode(0),
		hosts: make(map[string]*node),
	}
}

func (d *dispatcher) NotFoundHandler(route Route, handler Handler) {
	method := route.Method
	parts := pathSplit(route.Pattern)

	if len(route.Hosts) != 0 {
		for _, host := range route.Hosts {
			d.host(host).error(method, parts, handler)
		}
		return
	}

	d.error(method, parts, handler)
}

func (d *dispatcher) Middleware(route Route, middleware ...Middleware) {
	if len(middleware) == 0 {
		return
	}

	method := route.Method
	parts := pathSplit(route.Pattern)

	if len(route.Hosts) != 0 {
		for _, host := range route.Hosts {
			d.host(host).use(method, parts, middleware...)
		}
		return
	}

	d.use(method, parts, middleware...)
}

func (d *dispatcher) Handle(route Route, handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}

	method := route.Method
	parts := pathSplit(route.Pattern)

	if len(route.Hosts) != 0 {
		for _, host := range route.Hosts {
			d.host(host).handle(method, parts, handlers...)
		}
		return
	}

	d.handle(method, parts, handlers...)
}

func (d *dispatcher) CtxServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ctx, parts := URLParts(ctx, r)

	var h1, h2 Handler
	var m1, m2 []Middleware

	if host, ok := d.hosts[r.URL.Host]; ok {
		h1, m1, _ = host.resolve(r.Method, parts)
	}
	h2, m2, _ = d.resolve(r.Method, parts)

	if h := d.notFoundHandler(r.Method, h1, h2); h != nil {
		for _, m := range m1 {
			h = m(h)
		}
		for _, m := range m2 {
			h = m(h)
		}
		h.CtxServeHTTP(ctx, w, r)
		return
	}

	http.NotFound(w, r)
}

func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	d.CtxServeHTTP(ctx, w, r)
}

func (d *dispatcher) host(host string) *node {
	d.lock.RLock()
	if n, ok := d.hosts[host]; ok {
		return n
	}
	d.lock.RUnlock()

	d.lock.Lock()
	n := newNode(0)
	d.hosts[host] = n
	d.lock.Unlock()

	return n
}
