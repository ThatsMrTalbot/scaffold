package errors

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {
	Convey("Given an error", t, func() {
		err := errors.New("Some error message")
		Convey("When it is converted to include a status", func() {
			err = ConvertErrorStatus(404, err)
			Convey("Then it should contain the status", func() {
				So(err, ShouldImplement, (*ErrorStatus)(nil))
				So(err.(ErrorStatus).Status(), ShouldEqual, 404)
			})
		})
	})

	Convey("Given an error string", t, func() {
		msg := "Some error message"
		Convey("When it is converted to include a status", func() {
			err := NewErrorStatus(404, msg)
			Convey("Then it should contain the status", func() {
				So(err, ShouldImplement, (*ErrorStatus)(nil))
				So(err.(ErrorStatus).Status(), ShouldEqual, 404)
			})
		})
	})
}
