package scaffold

// Route is a route used in handling requests
type Route struct {
	Hosts   []string
	Pattern string
	Method  string
}

func (r Route) clone() Route {
	return Route{
		Hosts:   r.Hosts,
		Pattern: r.Pattern,
		Method:  r.Method,
	}
}
