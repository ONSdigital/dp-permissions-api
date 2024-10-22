package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/ONSdigital/dp-permissions-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI_GetPermissionsBundleHandler(t *testing.T) {
	adminPolicy := &models.BundlePolicy{
		ID: "1234",
		Entities: []string{
			"groups/admin",
		},
		Role: "admin",
	}
	publisherPolicy := &models.BundlePolicy{
		ID: "5678",
		Entities: []string{
			"groups/publisher",
		},
		Role: "publisher",
	}
	expectedBundle := models.Bundle{
		"legacy.read": map[string][]*models.BundlePolicy{
			"groups/admin": {
				adminPolicy,
			},
			"groups/publisher": {
				publisherPolicy,
			},
		},
	}

	Convey("Given a permissions bundler that returns a bundle", t, func() {
		bundler := &mock.PermissionsBundlerMock{
			GetFunc: func(ctx context.Context) (models.Bundle, error) {
				return expectedBundle, nil
			},
		}
		permissionsAPI := setupAPIWithBundler(bundler)

		Convey("When a GET request is made to the /v1/permissions-bundle endpoint", func() {
			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/permissions-bundle", http.NoBody)
			w := httptest.NewRecorder()
			permissionsAPI.Router.ServeHTTP(w, r)

			Convey("Then the bundle data is returned in the response body", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := io.ReadAll(w.Body)
				So(err, ShouldBeNil)
				expectedJSON, err := json.Marshal(expectedBundle)
				So(err, ShouldBeNil)
				So(payload, ShouldResemble, expectedJSON)
			})
		})
	})
}

func TestAPI_GetPermissionsBundleHandler_BundlerError(t *testing.T) {
	Convey("Given a permissions bundler that returns an error", t, func() {
		expectedError := errors.New("bundler error")
		bundler := &mock.PermissionsBundlerMock{
			GetFunc: func(ctx context.Context) (models.Bundle, error) {
				return nil, expectedError
			},
		}
		permissionsAPI := setupAPIWithBundler(bundler)

		Convey("When a GET request is made to the /v1/permissions-bundle endpoint", func() {
			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/permissions-bundle", http.NoBody)
			w := httptest.NewRecorder()
			permissionsAPI.Router.ServeHTTP(w, r)

			Convey("Then a 500 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
