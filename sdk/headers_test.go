package sdk

import (
	"net/http"
	"testing"

	dprequest "github.com/ONSdigital/dp-net/v3/request"

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

	Convey("Given Headers with a ServiceAuthToken", t, func() {
		headers := &Headers{
			ServiceAuthToken: "test-auth-token",
		}

		Convey("When Add is called", func() {
			headers.Add(req)
			expectedHeader := dprequest.BearerPrefix + headers.ServiceAuthToken

			Convey("Then an Authorization header is set on the request", func() {
				So(req.Header[dprequest.AuthHeaderKey][0], ShouldEqual, expectedHeader)
			})
		})
	})
}
