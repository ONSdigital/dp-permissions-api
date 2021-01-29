package api_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/api"
	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		mongoMock := &mock.PermissionsStoreMock{}

		r := mux.NewRouter()
		ctx := context.Background()
		api := api.Setup(ctx, r, mongoMock)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/roles/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/roles", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
