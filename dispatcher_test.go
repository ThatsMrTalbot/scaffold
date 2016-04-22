package scaffold

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func NewRoute(host string, method string, pattern string) Route {
	hosts := []string{}
	if host != "" {
		hosts = []string{host}
	}

	return Route{Hosts: hosts, Method: method, Pattern: pattern}
}

func TestDispatcherRoutes(t *testing.T) {
	hosts := []string{"", "example.com", "www.example.com"}
	methods := []string{"", "OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE"}

	patterns := []string{"some/path", ":var1/something/:var2"}
	paths := []string{"some/path", "value1/something/value2"}

	for _, host := range hosts {
		for _, method := range methods {
			for i, pattern := range patterns {

				convey := fmt.Sprintf("Given a route with host: `%s`, method: `%s` and pattern `%s`",
					host,
					method,
					pattern,
				)

				Convey(convey, t, func() {

					d := DefaultDispacher()
					r := NewRoute(host, method, pattern)
					handler := NewMockHandler()

					d.Handle(r, handler)

					for _, h := range hosts {
						if h == "" {
							continue
						}

						u := &url.URL{Scheme: "http", Host: h, Path: paths[i]}

						for _, m := range methods {
							if m == "" {
								continue
							}

							convey = fmt.Sprintf("When dispatching %s with method %s",
								u,
								m,
							)

							Convey(convey, func() {

								req, err := http.NewRequest(m, u.String(), nil)
								So(err, ShouldBeNil)

								d.ServeHTTP(httptest.NewRecorder(), req)

								correctHost := host == "" || h == host
								correctMethod := method == "" || m == method

								if correctHost && correctMethod {
									Convey("Then the handler should be hit", func() {
										So(handler.Request, ShouldNotBeNil)
										So(handler.Request, ShouldEqual, req)
									})
								} else {
									Convey("Then the handler should not be hit", func() {
										So(handler.Request, ShouldBeNil)
									})
								}
							})
						}
					}
				})
			}
		}
	}

	// Code coverage OCD

	Convey("Given a route", t, func() {
		d := DefaultDispacher()
		r := NewRoute("", "", "")

		Convey("When no handlers are provied", func() {
			d.Handle(r)

			Convey("Then no handlers should be added", func() {
				// Code coverage
			})
		})
	})

	Convey("Given a route with a host", t, func() {
		d := DefaultDispacher()
		r := NewRoute("example.com", "", "")
		h := NewMockHandler()

		Convey("When the host already exists", func() {
			d.Handle(r, h)
			d.Handle(r, h)

			Convey("Then no new host should not be created", func() {
				// Code coverage
			})
		})
	})
}

func TestDispatcherNotFound(t *testing.T) {
	hosts := []string{"", "example.com", "www.example.com"}
	methods := []string{"", "OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE"}

	patterns := []string{"some/path", "some/path", ""}
	paths := []string{"some/path", "some/path/here", "invalid"}

	for _, host := range hosts {
		for _, method := range methods {
			for i, pattern := range patterns {

				convey := fmt.Sprintf("Given a not found handler with host: `%s`, method: `%s` and pattern `%s`",
					host,
					method,
					pattern,
				)

				Convey(convey, t, func() {

					d := DefaultDispacher()
					r := NewRoute(host, method, pattern)
					handler := NewMockHandler()

					d.NotFoundHandler(r, handler)

					for _, h := range hosts {
						if h == "" {
							continue
						}

						u := &url.URL{Scheme: "http", Host: h, Path: paths[i]}

						for _, m := range methods {
							if m == "" {
								continue
							}

							convey = fmt.Sprintf("When dispatching %s with method %s",
								u,
								m,
							)

							Convey(convey, func() {

								req, err := http.NewRequest(m, u.String(), nil)
								So(err, ShouldBeNil)

								d.ServeHTTP(httptest.NewRecorder(), req)

								correctHost := host == "" || h == host
								correctMethod := method == "" || m == method

								if correctHost && correctMethod {
									Convey("Then the handler should be hit", func() {
										So(handler.Request, ShouldNotBeNil)
										So(handler.Request, ShouldEqual, req)
									})
								} else {
									Convey("Then the handler should not be hit", func() {
										So(handler.Request, ShouldBeNil)
									})
								}
							})
						}
					}
				})
			}
		}
	}
}

func TestDispacherMiddleware(t *testing.T) {
	hosts := []string{"", "example.com", "www.example.com"}
	methods := []string{"", "OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE"}

	middleware := []string{"some/path", "some/path", ""}
	patterns := []string{"some/path", "some/path/here", "some"}
	paths := []string{"some/path", "some/path/here", "some"}

	for _, host := range hosts {
		for _, method := range methods {
			for i, pattern := range patterns {

				convey := fmt.Sprintf("Given a not found handler with host: `%s`, method: `%s` and pattern `%s`",
					host,
					method,
					pattern,
				)

				Convey(convey, t, func() {

					d := DefaultDispacher()
					r1 := NewRoute(host, method, middleware[i])
					r2 := NewRoute(host, method, pattern)
					middleware := NewMockMiddleware()
					handler := NewMockHandler()

					d.Middleware(r1, middleware.Middleware())
					d.Handle(r2, handler)

					for _, h := range hosts {
						if h == "" {
							continue
						}

						u := &url.URL{Scheme: "http", Host: h, Path: paths[i]}

						for _, m := range methods {
							if m == "" {
								continue
							}

							convey = fmt.Sprintf("When dispatching %s with method %s",
								u,
								m,
							)

							Convey(convey, func() {

								req, err := http.NewRequest(m, u.String(), nil)
								So(err, ShouldBeNil)

								d.ServeHTTP(httptest.NewRecorder(), req)

								correctHost := host == "" || h == host
								correctMethod := method == "" || m == method

								if correctHost && correctMethod {
									Convey("Then the middleware should be hit", func() {
										So(middleware.Called, ShouldBeTrue)
									})
									Convey("Then the handler should be hit", func() {
										So(handler.Request, ShouldNotBeNil)
										So(handler.Request, ShouldEqual, req)
									})
								} else {
									Convey("Then the handler should not be hit", func() {
										So(handler.Request, ShouldBeNil)
									})
								}
							})
						}
					}
				})
			}
		}
	}

	// Code coverage OCD

	Convey("Given a route", t, func() {
		d := DefaultDispacher()
		r := NewRoute("", "", "")

		Convey("When no middleware is provied", func() {
			d.Middleware(r)

			Convey("Then no middleware should be added", func() {
				// Code coverage
			})
		})
	})
}

type MockDispatcher struct {
	Dispatcher

	HandlerHasBeenCalled         bool
	MiddlewareHasBeenCalled      bool
	NotFoundHandlerHasBeenCalled bool

	NotFound        Handler
	HandlerItems    []Handler
	MiddlewareItems []Middleware
	Route           *Route
}

func NewMockDispatcher() *MockDispatcher {
	return &MockDispatcher{
		Dispatcher: DefaultDispacher(),
	}
}

func (m *MockDispatcher) Handle(route Route, handlers ...Handler) {
	m.Route = &route
	m.HandlerItems = handlers
	m.HandlerHasBeenCalled = true

	m.Dispatcher.Handle(route, handlers...)
}

func (m *MockDispatcher) Middleware(route Route, middleware ...Middleware) {
	m.Route = &route
	m.MiddlewareItems = middleware
	m.MiddlewareHasBeenCalled = true

	m.Dispatcher.Middleware(route, middleware...)
}

func (m *MockDispatcher) NotFoundHandler(route Route, handler Handler) {
	m.Route = &route
	m.NotFound = handler
	m.NotFoundHandlerHasBeenCalled = true

	m.Dispatcher.NotFoundHandler(route, handler)
}
