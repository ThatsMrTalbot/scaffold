package encoding

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/context"

	"github.com/ThatsMrTalbot/scaffold"
	. "github.com/smartystreets/goconvey/convey"
)

type MockHandler struct {
	Called  bool
	Obj     MockObject
	Context context.Context
	Request *http.Request
	Writer  http.ResponseWriter
}

func (m *MockHandler) Handler(obj MockObject, ctx context.Context, w http.ResponseWriter, r *http.Request) (MockObject, error) {
	m.Obj = obj
	m.Called = true
	m.Context = ctx
	m.Request = r
	m.Writer = w

	return obj, nil
}

func TestEncoder(t *testing.T) {
	Convey("Given the default encoder", t, func() {
		Convey("When a request with no Content-Type header is parsed", func() {
			req, err := http.NewRequest("GET", "http://www.example.com", nil)
			So(err, ShouldBeNil)

			parser := DefaultEncoder.Parser(req)

			Convey("The default parser should be returned", func() {
				So(parser, ShouldEqual, JSONEncoding)
			})
		})

		Convey("When a request with valid Content-Type header is parsed", func() {
			req, err := http.NewRequest("GET", "http://www.example.com", nil)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/xml")
			parser := DefaultEncoder.Parser(req)

			Convey("The correct parser should be returned", func() {
				So(parser, ShouldEqual, XMLEncoding)
			})
		})

		Convey("When a request with invalid Content-Type header is parsed", func() {
			req, err := http.NewRequest("GET", "http://www.example.com", nil)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/invalid")
			parser := DefaultEncoder.Parser(req)

			Convey("The default parser should be returned", func() {
				So(parser, ShouldEqual, JSONEncoding)
			})
		})

		Convey("When a request with no Accept header is parsed", func() {
			req, err := http.NewRequest("GET", "http://www.example.com", nil)
			So(err, ShouldBeNil)

			responder := DefaultEncoder.Responder(req)

			Convey("The default responder should be returned", func() {
				So(responder, ShouldEqual, JSONEncoding)
			})
		})

		Convey("When a request with a Content-Type header but no Accept header is parsed", func() {
			req, err := http.NewRequest("GET", "http://www.example.com", nil)
			So(err, ShouldBeNil)

			req.Header.Set("Content-Type", "application/xml")
			responder := DefaultEncoder.Responder(req)

			Convey("The responder returned should match the parser", func() {
				So(responder, ShouldEqual, XMLEncoding)
			})
		})

		Convey("When a request with invalid Accept header is parsed", func() {
			req, err := http.NewRequest("GET", "http://www.example.com", nil)
			So(err, ShouldBeNil)

			req.Header.Set("Accept", "application/invalid")
			responder := DefaultEncoder.Responder(req)

			Convey("The default responder should be returned", func() {
				So(responder, ShouldEqual, JSONEncoding)
			})
		})

		Convey("When a request with valid Accept header is parsed", func() {
			req, err := http.NewRequest("GET", "http://www.example.com", nil)
			So(err, ShouldBeNil)

			req.Header.Set("Accept", "application/xml, application/json; a=b;q=.2, invalid")
			responder := DefaultEncoder.Responder(req)

			Convey("The default responder should be returned", func() {
				So(responder, ShouldEqual, XMLEncoding)
			})
		})

		Convey("When a handler is built", func() {
			mock := &MockHandler{}
			handler, err := DefaultEncoder.HandlerBuilder(mock.Handler)
			So(err, ShouldBeNil)

			Convey("The handler should be valid", func() {
				So(handler, ShouldImplement, (*scaffold.Handler)(nil))
			})
		})

		Convey("When a built handler is called", func() {
			mock := &MockHandler{}
			handler, err := DefaultEncoder.HandlerBuilder(mock.Handler)
			So(err, ShouldBeNil)

			ctx := context.Background()
			r, err := http.NewRequest("GET", "http://www.example.com", strings.NewReader(`{"a":"123","b":321}`))
			So(err, ShouldBeNil)

			w := httptest.NewRecorder()

			handler.CtxServeHTTP(ctx, w, r)

			Convey("The result should be valid", func() {
				So(mock.Called, ShouldBeTrue)
				So(mock.Context, ShouldEqual, ctx)
				So(mock.Request, ShouldEqual, r)
				So(mock.Writer, ShouldEqual, w)
				So(w.Body.String(), ShouldEqual, `{"a":"123","b":321}`)
			})
		})
	})
}
