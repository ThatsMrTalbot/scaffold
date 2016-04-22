package scaffold

// Platform is a object that can specify its own routes
type Platform interface {
	Routes(*Router)
}
