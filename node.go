package scaffold

import "sync"

type node struct {
	lock sync.RWMutex

	depth     int
	children  map[string]*node
	handlers  map[string]Handler
	middlware map[string][]Middleware
	notfound  map[string]Handler
}

func newNode(depth int) *node {
	return &node{
		depth:     depth,
		children:  make(map[string]*node),
		handlers:  make(map[string]Handler),
		middlware: make(map[string][]Middleware),
		notfound:  make(map[string]Handler),
	}
}

func (n *node) node(parts []string, vars map[int]string) *node {
	if len(parts) == 0 {
		return n
	}

	next, parts := parts[0], parts[1:]
	if next[0] == ':' {
		vars[n.depth], next = next[1:], ""
	}

	n.lock.Lock()
	if _, ok := n.children[next]; !ok {
		n.children[next] = newNode(n.depth + 1)
	}
	n.lock.Unlock()

	n.lock.RLock()
	defer n.lock.RUnlock()

	return n.children[next].node(parts, vars)
}

func (n *node) handle(method string, parts []string, handlers ...Handler) {
	vars := make(map[int]string)
	node := n.node(parts, vars)

	node.lock.Lock()
	node.handlers[method] = wrapHandler(vars, HandlerList(handlers))
	node.lock.Unlock()
}

func (n *node) use(method string, parts []string, middleware ...Middleware) {
	vars := make(map[int]string)
	node := n.node(parts, vars)

	m := wrapMiddleware(vars, MiddlewareList(middleware))

	node.lock.Lock()
	node.middlware[method] = append(node.middlware[method], m)
	node.lock.Unlock()
}

func (n *node) error(method string, parts []string, handler Handler) {
	vars := make(map[int]string)
	node := n.node(parts, vars)

	node.lock.Lock()
	node.notfound[method] = wrapHandler(vars, handler)
	node.lock.Unlock()
}

func (n *node) resolve(method string, parts []string) (Handler, []Middleware, bool) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	m1 := n.middleware(method)

	if len(parts) == 0 {
		if h, ok := n.handler(method); ok {
			return h, m1, ok
		}
		return n.notFoundHandler(method, nil, nil), m1, false
	}

	next, parts := parts[0], parts[1:]

	var h1, h2 Handler
	var m2 []Middleware
	if c, ok := n.children[next]; ok {
		if h1, m2, ok = c.resolve(method, parts); ok {
			return h1, append(m1, m2...), ok
		}
	}
	if c, ok := n.children[""]; ok {
		if h2, m2, ok = c.resolve(method, parts); ok {
			return h2, append(m1, m2...), ok
		}
	}

	return n.notFoundHandler(method, h1, h2), m1, false
}

func (n *node) handler(method string) (Handler, bool) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	if h, ok := n.handlers[method]; ok {
		return h, true
	}
	if h, ok := n.handlers[""]; ok {
		return h, true
	}
	return nil, false
}

func (n *node) notFoundHandler(method string, h1 Handler, h2 Handler) Handler {
	n.lock.RLock()
	defer n.lock.RUnlock()

	if h1 != nil {
		return h1
	}
	if h2 != nil {
		return h2
	}
	if h, ok := n.notfound[method]; ok {
		return h
	}
	if h, ok := n.notfound[""]; ok {
		return h
	}
	return nil
}

func (n *node) middleware(method string) []Middleware {
	n.lock.RLock()
	defer n.lock.RUnlock()

	middleware := make([]Middleware, 0, 10)

	if m, ok := n.middlware[method]; ok {
		middleware = append(middleware, m...)
	}
	if m, ok := n.middlware[""]; ok {
		middleware = append(middleware, m...)
	}
	return middleware
}
