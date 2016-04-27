package errors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	. "github.com/smartystreets/goconvey/convey"
)

func MockHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return NewErrorStatus(500, "Some error message")
}

func TestBuild(t *testing.T) {
	Convey("Given an handler that returns an error", t, func() {
		h, err := HandlerBuilder(MockHandler)
		So(err, ShouldBeNil)

		Convey("When the handler is called", func() {
			ctx := context.Background()
			r, _ := http.NewRequest("GET", "http://www.example.com", nil)
			w := httptest.NewRecorder()

			h.CtxServeHTTP(ctx, w, r)

			Convey("Then an error should be retuned", func() {
				So(w.Body.String(), ShouldContainSubstring, "Some error message")
				So(w.Code, ShouldEqual, 500)
			})
		})
	})
}
