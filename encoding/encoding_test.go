package encoding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MockObject struct {
	A string `json:"a" xml:"A"`
	B int    `json:"b" xml:"B"`
}

func TestEncoding(t *testing.T) {
	Convey("Given a json encoder", t, func() {
		Convey("When a value is encoded", func() {
			obj := MockObject{
				A: "123",
				B: 321,
			}
			w := httptest.NewRecorder()
			JSONEncoding.Respond(200, w, obj)

			Convey("Then when decoded the values should match", func() {
				var result MockObject
				r, err := http.NewRequest("GET", "http://www.example.com", w.Body)
				So(err, ShouldBeNil)

				JSONEncoding.Parse(&result, r)
				So(result, ShouldResemble, obj)
			})
		})
	})

	Convey("Given a xml encoder", t, func() {
		Convey("When a value is encoded", func() {
			obj := MockObject{
				A: "123",
				B: 321,
			}
			w := httptest.NewRecorder()
			XMLEncoding.Respond(200, w, obj)

			Convey("Then when decoded the values should match", func() {
				var result MockObject
				r, err := http.NewRequest("GET", "http://www.example.com", w.Body)
				So(err, ShouldBeNil)

				XMLEncoding.Parse(&result, r)
				So(result, ShouldResemble, obj)
			})
		})
	})
}
