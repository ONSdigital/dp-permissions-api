package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/apierrors"

	"github.com/gorilla/mux"

	"github.com/ONSdigital/dp-permissions-api/api"
	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/ONSdigital/dp-permissions-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testRoleID1 = "roleID1"
)

func dbRole(id string) *models.Role {
	return &models.Role{
		ID:   id,
		Name: "ReadOnly",
		Permissions: []string{
			"read",
		},
	}
}

func TestGetRoleHandler(t *testing.T) {

	Convey("Given a GetRole Handler", t, func() {

		mockedPermissionsStore := &mock.PermissionsStoreMock{
			GetRoleFunc: func(ctx context.Context, id string) (*models.Role, error) {
				switch id {
				case testRoleID1:
					return &models.Role{ID: "testRoleID1", Name: "ReadOnly", Permissions: []string{"read"}}, nil
				default:
					return nil, apierrors.ErrRoleNotFound
				}
			},
		}

		permissionsApi := api.Setup(context.Background(), mux.NewRouter(), mockedPermissionsStore)

		Convey("When an existing role is requested with its Role ID", func() {

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25400/roles/%s", testRoleID1), nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("The the matched role is returned with status code 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedRole := models.Role{}
				err = json.Unmarshal(payload, &returnedRole)
				So(err, ShouldBeNil)
				So(returnedRole, ShouldResemble, *dbRole("testRoleID1"))
			})
		})

		Convey("When a non existing role is requested a Not Found response is returned", func() {

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/roles/inexistent", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusNotFound)

		})
	})

}
