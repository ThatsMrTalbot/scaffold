package scaffold

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouter(t *testing.T) {
	Convey("Given a router", t, func() {
		h := NewMockHandler()
		d := NewMockDispatcher()
		r := New(d)

		methods := map[string]func(pattern string, handlers ...Handler) *Router{
			"OPTIONS": r.Options,
			"GET":     r.Get,
			"HEAD":    r.Head,
			"POST":    r.Post,
			"PUT":     r.Put,
			"DELETE":  r.Delete,
			"":        r.Handle,
		}

		for method, handler := range methods {
			when := fmt.Sprintf("When a route with method `%s` is added", method)
			Convey(when, func() {
				r = handler("somepath", h)

				then := fmt.Sprintf("Then the route method should be `%s`", method)
				Convey(then, func() {
					So(d.HandlerHasBeenCalled, ShouldBeTrue)
					So(d.Route, ShouldNotBeNil)
					So(d.Route.Method, ShouldEqual, method)
				})

				Convey("Then the route pattern should be correct", func() {
					So(d.HandlerHasBeenCalled, ShouldBeTrue)
					So(d.Route, ShouldNotBeNil)
					So(d.Route.Pattern, ShouldEqual, "somepath")
				})
			})
		}

		Convey("When a host is set", func() {
			r = r.Host("example.com", "www.example.com")

			Convey("Then the route should contain the hosts", func() {
				So(r.route.Hosts, ShouldContain, "example.com")
				So(r.route.Hosts, ShouldContain, "www.example.com")
			})
		})

		Convey("When a route pattern is set", func() {
			r = r.Route("somepath")

			Convey("Then the route pattern should be correct", func() {
				So(r.route.Pattern, ShouldEqual, "somepath")
			})
		})

		Convey("When a group is used", func() {
			r.Group("somepath", func(r *Router) {
				Convey("Then the route pattern should be correct", func() {
					So(r.route.Pattern, ShouldEqual, "somepath")
				})
			})
		})

		Convey("When a plaform is added", func() {
			platform := NewMockPlatform()
			r.Platform("somepath", platform)

			Convey("Then the route pattern should be correct", func() {
				So(platform.Route, ShouldNotBeNil)
				So(platform.Route.Pattern, ShouldEqual, "somepath")
			})
		})

		Convey("When middleware is added", func() {
			m := NewMockMiddleware()
			middleware := m.Middleware()
			r.Use(middleware)

			Convey("Then the middleware should be added", func() {
				So(d.MiddlewareHasBeenCalled, ShouldBeTrue)
				So(d.MiddlewareItems, ShouldContain, middleware)
			})
		})

		Convey("When middleware function is added", func() {
			m := NewMockMiddleware()
			middleware := m.Middleware()
			r.UseFunc(middleware)

			Convey("Then the middleware should be added", func() {
				So(d.MiddlewareHasBeenCalled, ShouldBeTrue)
				So(d.MiddlewareItems, ShouldContain, middleware)
			})
		})

		Convey("When a not found handler is added", func() {
			h := NewMockHandler()
			r.NotFound(h)

			Convey("Then the not found handler should be set", func() {
				So(d.NotFound, ShouldEqual, h)
			})
		})
	})

	Convey("Given a platform", t, func() {
		d := NewMockDispatcher()
		platform := NewMockPlatform()

		Convey("When the platform is scaffolded", func() {
			Scaffold(d, platform)

			Convey("The platform should be routed", func() {
				So(platform.Route, ShouldNotBeNil)
			})
		})
	})
}
