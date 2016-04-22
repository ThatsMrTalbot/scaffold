package scaffold

type MockPlatform struct {
	Route *Route
}

func (m *MockPlatform) Routes(r *Router) {
	m.Route = &r.route
}

func NewMockPlatform() *MockPlatform {
	return &MockPlatform{}
}
