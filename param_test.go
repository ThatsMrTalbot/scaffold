package scaffold

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func TestParam(t *testing.T) {
	Convey("Given a context", t, func() {
		ctx := context.Background()

		Convey("When a value is stored in a context", func() {
			ctx = StoreParam(ctx, "somekey", "somevalue")
			Convey("Then the value should be in the context", func() {
				val, err := GetParam(ctx, "somekey").String()
				So(err, ShouldBeNil)
				So(val, ShouldEqual, "somevalue")
			})
		})

		Convey("When a value that does not exists is retrieved", func() {
			Convey("Then an empty value should be retrieved", func() {
				val := GetParam(ctx, "missingkey")
				So(val, ShouldBeEmpty)
			})
		})

		Convey("When a numeric value is retrieved from a context", func() {
			ctx = StoreParam(ctx, "somekey", "42")

			Convey("Then the param should be convertable to numeric types", func() {
				p := GetParam(ctx, "somekey")

				i, err := p.Int()
				So(err, ShouldBeNil)
				i32, err := p.Int32()
				So(err, ShouldBeNil)
				i64, err := p.Int64()
				So(err, ShouldBeNil)

				u, err := p.UInt()
				So(err, ShouldBeNil)
				u32, err := p.UInt32()
				So(err, ShouldBeNil)
				u64, err := p.UInt64()
				So(err, ShouldBeNil)

				f32, err := p.Float32()
				So(err, ShouldBeNil)
				f64, err := p.Float64()
				So(err, ShouldBeNil)

				So(i, ShouldEqual, 42)
				So(i32, ShouldEqual, 42)
				So(i64, ShouldEqual, 42)
				So(u, ShouldEqual, 42)
				So(u32, ShouldEqual, 42)
				So(u64, ShouldEqual, 42)
				So(f32, ShouldEqual, 42)
				So(f64, ShouldEqual, 42)
			})
		})
	})
}
