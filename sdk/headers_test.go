package sdk

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHeaders_Add(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		Header: http.Header{},
	}

	Convey("Given empty Headers", t, func() {
		var headers Headers

		Convey("When Add is called", func() {
			headers.Add(req)

			Convey("Then no headers are set on the request", func() {
				So(req.Header, ShouldBeEmpty)
			})
		})
	})

	Convey("Given Headers with an Authorization token", t, func() {
		headers := &Headers{
			Authorization: "test-auth-token",
		}

		Convey("When Add is called", func() {
			headers.Add(req)
			expectedHeader := BearerPrefix + headers.Authorization

			Convey("Then an Authorization header is set on the request", func() {
				So(req.Header[Authorization][0], ShouldEqual, expectedHeader)
			})
		})
	})
}
