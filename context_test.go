package scaffold

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func TestContext(t *testing.T) {
	Convey("Given wrapped and unwapped middlware functions", t, func() {
		vars := map[int]string{0: "param1", 2: "param2"}

		m1 := Middleware(func(next Handler) Handler {
			return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				p1 := GetParam(ctx, "param1")
				p2 := GetParam(ctx, "param2")
				So(p1, ShouldEqual, "value1")
				So(p2, ShouldEqual, "value2")
				next.CtxServeHTTP(ctx, w, r)
			})
		})
		m1 = wrapMiddleware(vars, m1)

		m2 := Middleware(func(next Handler) Handler {
			return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				p1 := GetParam(ctx, "param1")
				p2 := GetParam(ctx, "param2")
				So(p1, ShouldBeEmpty)
				So(p2, ShouldBeEmpty)
				next.CtxServeHTTP(ctx, w, r)
			})
		})

		Convey("When the handler is called", func() {
			handler := Handler(NewMockHandler())
			handler = m2(m1(m2(handler)))
			req, err := http.NewRequest("GET", "http://example.com/value1/and/value2/end", nil)
			So(err, ShouldBeNil)

			Convey("Then the params should be in only the wrapped middleware context", func() {
				handler.CtxServeHTTP(context.Background(), httptest.NewRecorder(), req)
			})

		})
	})

	Convey("Given a wrapped handler", t, func() {
		vars := map[int]string{0: "param1", 2: "param2"}

		handler := Handler(HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			p1 := GetParam(ctx, "param1")
			p2 := GetParam(ctx, "param2")
			So(p1, ShouldEqual, "value1")
			So(p2, ShouldEqual, "value2")
		}))
		handler = wrapHandler(vars, handler)

		Convey("When the handler is called", func() {
			req, err := http.NewRequest("GET", "http://example.com/value1/and/value2/end", nil)
			So(err, ShouldBeNil)

			Convey("Then the params should be in the handler context", func() {
				handler.CtxServeHTTP(context.Background(), httptest.NewRecorder(), req)
			})

		})

	})
}
