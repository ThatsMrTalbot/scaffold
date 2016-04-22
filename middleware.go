package scaffold

// Middleware is middleware for a router
type Middleware func(next Handler) Handler

// MiddlewareList list creates a single middlware callback from an array
func MiddlewareList(middlware []Middleware) Middleware {
	return Middleware(func(next Handler) Handler {
		for _, m := range middlware {
			next = m(next)
		}
		return next
	})
}
