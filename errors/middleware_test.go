package errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/ThatsMrTalbot/scaffold"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMiddleware(t *testing.T) {
	Convey("Given error middleware", t, func() {
		handler := scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			GetErrorHandler(ctx, 404).ServeErrorPage(ctx, w, r, 404, errors.New("Some error message"))
		})
		errHandler := ErrorHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error) {
			w.WriteHeader(status)
			w.Write([]byte(err.Error()))
		})

		middleware := SetErrorHandler(404, errHandler)
		h := middleware(handler)

		Convey("When a wrapped handler is called", func() {
			ctx := context.Background()
			r, _ := http.NewRequest("GET", "http://www.example.com", nil)
			w := httptest.NewRecorder()

			h.CtxServeHTTP(ctx, w, r)

			Convey("Then the response should contain the error", func() {
				So(w.Body.String(), ShouldEqual, "Some error message")
			})
		})

	})

	Convey("Given no error middleware", t, func() {
		handler := scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			GetErrorHandler(ctx, 404).ServeErrorPage(ctx, w, r, 404, errors.New("Some error message"))
		})

		Convey("When a handler is called", func() {
			ctx := context.Background()
			r, _ := http.NewRequest("GET", "http://www.example.com", nil)
			w := httptest.NewRecorder()

			handler.CtxServeHTTP(ctx, w, r)

			Convey("Then the response should contain the default error", func() {
				So(w.Body.String(), ShouldContainSubstring, "Some error message")
			})
		})

	})
}
