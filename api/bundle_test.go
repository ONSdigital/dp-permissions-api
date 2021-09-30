package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/ONSdigital/dp-permissions-api/models"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_GetPermissionsBundleHandler(t *testing.T) {
	adminPolicy := &models.Policy{
		Entities: []string{
			"groups/admin",
		},
		Role: "admin",
	}
	publisherPolicy := &models.Policy{
		Entities: []string{
			"groups/publisher",
		},
		Role: "publisher",
	}
	expectedBundle := &models.Bundle{
		PermissionToEntityLookup: map[string]models.EntityIDToPolicies{
			"legacy.read": map[string][]*models.Policy{
				"groups/admin": {
					adminPolicy,
				},
				"groups/publisher": {
					publisherPolicy,
				},
			},
		},
	}

	Convey("Given a permissions bundler that returns a bundle", t, func() {
		bundler := &mock.PermissionsBundlerMock{
			GetFunc: func(ctx context.Context) (*models.Bundle, error) {
				return expectedBundle, nil
			},
		}
		permissionsApi := setupAPIWithBundler(bundler)

		Convey("When a GET request is made to the /v1/permissions-bundle endpoint", func() {

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/permissions-bundle", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("Then the bundle data is returned in the response body", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				actualBundle := &models.Bundle{}
				err = json.Unmarshal(payload, &actualBundle)
				So(err, ShouldBeNil)
				So(actualBundle, ShouldResemble, expectedBundle)
			})
		})
	})
}

func TestAPI_GetPermissionsBundleHandler_BundlerError(t *testing.T) {

	Convey("Given a permissions bundler that returns an error", t, func() {
		expectedError := errors.New("bundler error")
		bundler := &mock.PermissionsBundlerMock{
			GetFunc: func(ctx context.Context) (*models.Bundle, error) {
				return nil, expectedError
			},
		}
		permissionsApi := setupAPIWithBundler(bundler)

		Convey("When a GET request is made to the /v1/permissions-bundle endpoint", func() {

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/permissions-bundle", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("Then a 500 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
