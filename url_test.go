package scaffold

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestURL(t *testing.T) {
	Convey("Given a request", t, func() {
		req, err := http.NewRequest("POST", "http://www.example.com/some/path/here", nil)
		So(err, ShouldBeNil)

		Convey("When the url is split", func() {
			ctx, parts := URLParts(nil, req)

			Convey("Then the parts should be correct", func() {
				So(parts, ShouldResemble, []string{"some", "path", "here"})
			})

			Convey("Then the parts should exist in the context", func() {
				parts, ok := ctx.Value("scaffold_url_parts").([]string)
				So(ok, ShouldBeTrue)
				So(parts, ShouldResemble, []string{"some", "path", "here"})
			})
		})

		Convey("When a url part is a accessed", func() {
			_, part, ok := URLPart(nil, req, 1)

			Convey("Then the part should be correct", func() {
				So(ok, ShouldBeTrue)
				So(part, ShouldEqual, "path")
			})
		})

		Convey("When a url part that does not exist is a accessed", func() {
			_, part, ok := URLPart(nil, req, 10)

			Convey("Then the part should be correct", func() {
				So(ok, ShouldBeFalse)
				So(part, ShouldEqual, "")
			})
		})
	})
}
