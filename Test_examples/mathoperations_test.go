package mathoperations

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSum(t *testing.T) {
	a := 1
	b := 2
	Convey("Given a = 1 and b = 2", t, func() {
		res := Sum(a, b)
		Convey("Their sum should be 3", func() {
			So(res, ShouldEqual, 3)
		})
	})
}

func TestDif(t *testing.T) {
	a := 9
	b := 6
	Convey("Given a = 9 and b = 6", t, func() {
		res := Dif(a, b)
		Convey("Their difference should be 3", func() {
			So(res, ShouldEqual, 3)
		})
	})
}

func TestProd(t *testing.T) {
	a := 4
	b := 16
	Convey("Given a = 4 and b = 16", t, func() {
		res := Prod(a, b)
		Convey("Their product should be 64", func() {
			So(res, ShouldEqual, 64)
		})
	})
}

func TestDiv(t *testing.T) {
	a := 40
	b := 8
	Convey("Given a = 40 and b = 8", t, func() {
		res := Div(a, b)
		Convey("Their quotient should be 5", func() {
			So(res, ShouldEqual, 5)
		})
	})
}
