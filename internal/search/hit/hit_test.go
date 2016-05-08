package hit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHit(t *testing.T) {

	Convey("Should create hit if bytes and bounds are as expected.", t, func() {
		h, err := NewHit([]byte("hello"), []int{2, 4})
		So(h, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Should not create hit if the bounds are not in order.", t, func() {
		h, err := NewHit([]byte("hello"), []int{4, 2})
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("New hit should faile to create a hit if the bounds are not within byte offsets.", t, func() {
		h, err := NewHit([]byte("hello"), []int{2, 10})
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("New hit should faile to create a hit if the bounds are not withing byte offsets (negative).", t, func() {
		h, err := NewHit([]byte("hello"), []int{-1, 0})
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("New hit should fail creation for nil bounds.", t, func() {
		h, err := NewHit([]byte("hello"), nil)
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("New hit should fail creation for 1 offset bounds.", t, func() {
		h, err := NewHit([]byte("hello"), []int{})
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("New hit should fail creation for empty bounds.", t, func() {
		h, err := NewHit([]byte("hello"), []int{})
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("New hit should fail creation for empty bytes.", t, func() {
		h, err := NewHit([]byte{}, []int{1, 2})
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})

	Convey("New hit should fail creation for nil bytes.", t, func() {
		h, err := NewHit(nil, []int{1, 2})
		So(h, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}
