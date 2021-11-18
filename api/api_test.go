package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/api"
	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	authorisation "github.com/ONSdigital/dp-authorisation/v2/authorisation/mock"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		mongoMock := &mock.PermissionsStoreMock{}
		bundlerMock := &mock.PermissionsBundlerMock{}

		cfg := &config.Config{}
		r := mux.NewRouter()
		api := api.Setup(cfg, r, mongoMock, bundlerMock, newAuthMiddlwareMock())

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/v1/roles", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/roles/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/policies", "POST"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/policies/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/policies/{id}", "PUT"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/policies/{id}", "DELETE"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/permissions-bundle", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

var cfg = &config.Config{
	DefaultLimit:        20,
	DefaultOffset:       0,
	MaximumDefaultLimit: 1000,
}

func setupAPI() *api.API {
	return setupAPIWithStore(&mock.PermissionsStoreMock{})
}

func setupAPIWithStore(permissionsStore api.PermissionsStore) *api.API {
	return api.Setup(cfg, mux.NewRouter(), permissionsStore, &mock.PermissionsBundlerMock{}, newAuthMiddlwareMock())
}

func setupAPIWithBundler(bundler api.PermissionsBundler) *api.API {
	return api.Setup(cfg, mux.NewRouter(), &mock.PermissionsStoreMock{}, bundler, newAuthMiddlwareMock())
}

func newAuthMiddlwareMock() *authorisation.MiddlewareMock {
	return &authorisation.MiddlewareMock{
		RequireFunc: func(permission string, handlerFunc http.HandlerFunc) http.HandlerFunc {
			return handlerFunc
		},
	}
}
