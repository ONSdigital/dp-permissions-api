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
	"github.com/ONSdigital/dp-permissions-api/config"

	"github.com/gorilla/mux"

	"github.com/ONSdigital/dp-permissions-api/api"
	"github.com/ONSdigital/dp-permissions-api/api/mock"
	"github.com/ONSdigital/dp-permissions-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testRoleID1 = "testRoleID1"
	testRoleID2 = "testRoleID2"
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

var imageList = models.Roles{
	Items:      []models.Role{*dbRole(testRoleID1), *dbRole(testRoleID2)},
	Count:      2,
	Limit:      20,
	TotalCount: 2,
	Offset:     0,
}

var emptyImageList = models.Roles{
	Items:      []models.Role{},
	Count:      0,
	Limit:      20,
	TotalCount: 0,
	Offset:     0,
}

var paginatedImageList = models.Roles{
	Items:      []models.Role{*dbRole(testRoleID2)},
	Count:      1,
	Limit:      1,
	TotalCount: 1,
	Offset:     1,
}

var negativeQueryImageList = models.Roles{
	Items:      []models.Role{*dbRole(testRoleID1), *dbRole(testRoleID2)},
	Count:      2,
	Limit:      0,
	TotalCount: 2,
	Offset:     0,
}

var cfg = &config.Config{
	DefaultLimit:        20,
	DefaultOffset:       0,
	MaximumDefaultLimit: 1000,
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

		permissionsApi := api.Setup(context.Background(), &config.Config{}, mux.NewRouter(), mockedPermissionsStore)

		Convey("When an existing role is requested with its Role ID", func() {

			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25400/roles/%s", testRoleID1), nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("The matched role is returned with status code 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedRole := models.Role{}
				err = json.Unmarshal(payload, &returnedRole)
				So(err, ShouldBeNil)
				So(returnedRole, ShouldResemble, *dbRole(testRoleID1))
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

func TestGetRolesHandler(t *testing.T) {

	Convey("Given a GetRoles Handler", t, func() {

		Convey("When existing roles are requested using the default offset and limit values", func() {

			mockedPermissionsStore := &mock.PermissionsStoreMock{
				GetRolesFunc: func(ctx context.Context, offset int, limit int) (*models.Roles, error) {
					return &models.Roles{
						Count:      2,
						Offset:     0,
						Limit:      20,
						Items:      []models.Role{{ID: "testRoleID1", Name: "ReadOnly", Permissions: []string{"read"}}, {ID: "testRoleID2", Name: "ReadOnly", Permissions: []string{"read"}}},
						TotalCount: 2,
					}, nil
				},
			}

			permissionsApi := api.Setup(context.Background(), cfg, mux.NewRouter(), mockedPermissionsStore)

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/roles", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("The list of roles are returned with status code 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedRoles := models.Roles{}
				err = json.Unmarshal(payload, &returnedRoles)
				So(err, ShouldBeNil)
				So(returnedRoles, ShouldResemble, imageList)
			})

		})

		Convey("When existing roles are requested using valid user defined offset and limit values", func() {

			mockedPermissionsStore := &mock.PermissionsStoreMock{
				GetRolesFunc: func(ctx context.Context, offset int, limit int) (*models.Roles, error) {
					return &models.Roles{
						Count:      1,
						Offset:     1,
						Limit:      1,
						Items:      []models.Role{{ID: "testRoleID2", Name: "ReadOnly", Permissions: []string{"read"}}},
						TotalCount: 1,
					}, nil
				},
			}

			permissionsApi := api.Setup(context.Background(), cfg, mux.NewRouter(), mockedPermissionsStore)

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/roles?offset=1&limit=1", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("The paginated list of roles are returned with status code 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedRoles := models.Roles{}
				err = json.Unmarshal(payload, &returnedRoles)
				So(err, ShouldBeNil)
				So(returnedRoles, ShouldResemble, paginatedImageList)
			})

		})

		Convey("When existing roles are requested with a limit that is higher than the maximum default limit value", func() {

			mockedPermissionsStore := &mock.PermissionsStoreMock{
				GetRolesFunc: func(ctx context.Context, offset int, limit int) (*models.Roles, error) {
					return &models.Roles{
						Count:      1,
						Offset:     0,
						Limit:      1500,
						Items:      []models.Role{{ID: "testRoleID2", Name: "ReadOnly", Permissions: []string{"read"}}},
						TotalCount: 1,
					}, nil
				},
			}

			permissionsApi := api.Setup(context.Background(), cfg, mux.NewRouter(), mockedPermissionsStore)

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/roles?limit=1500", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("A status code of 400 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})

		})

		Convey("When existing roles are requested using negative user defined offset and limit values", func() {

			mockedPermissionsStore := &mock.PermissionsStoreMock{
				GetRolesFunc: func(ctx context.Context, offset int, limit int) (*models.Roles, error) {
					return &models.Roles{
						Count:      2,
						Offset:     0,
						Limit:      0,
						Items:      []models.Role{{ID: "testRoleID1", Name: "ReadOnly", Permissions: []string{"read"}}, {ID: "testRoleID2", Name: "ReadOnly", Permissions: []string{"read"}}},
						TotalCount: 2,
					}, nil
				},
			}

			permissionsApi := api.Setup(context.Background(), cfg, mux.NewRouter(), mockedPermissionsStore)

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/roles?offset=-1&limit=-1", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("A status code of 400 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})

		})

		Convey("When existing roles are requested using invalid user defined offset and limit values", func() {

			mockedPermissionsStore := &mock.PermissionsStoreMock{
				GetRolesFunc: func(ctx context.Context, offset int, limit int) (*models.Roles, error) {
					return &models.Roles{
						Count:      2,
						Offset:     0,
						Limit:      0,
						Items:      []models.Role{{ID: "testRoleID1", Name: "ReadOnly", Permissions: []string{"read"}}, {ID: "testRoleID2", Name: "ReadOnly", Permissions: []string{"read"}}},
						TotalCount: 2,
					}, nil
				},
			}

			permissionsApi := api.Setup(context.Background(), cfg, mux.NewRouter(), mockedPermissionsStore)

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/roles?offset=h&limit=i", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("A status code of 400 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})

		})

		Convey("When non existing roles are requested using the default offset and limit values", func() {

			mockedPermissionsStore := &mock.PermissionsStoreMock{
				GetRolesFunc: func(ctx context.Context, offset int, limit int) (*models.Roles, error) {
					return &models.Roles{
						Count:      0,
						Offset:     0,
						Limit:      20,
						Items:      []models.Role{},
						TotalCount: 0,
					}, nil
				},
			}

			permissionsApi := api.Setup(context.Background(), &config.Config{}, mux.NewRouter(), mockedPermissionsStore)

			r := httptest.NewRequest(http.MethodGet, "http://localhost:25400/v1/roles", nil)
			w := httptest.NewRecorder()
			permissionsApi.Router.ServeHTTP(w, r)

			Convey("The list of roles are returned with status code 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedRoles := models.Roles{}
				err = json.Unmarshal(payload, &returnedRoles)
				So(err, ShouldBeNil)
				So(returnedRoles, ShouldResemble, emptyImageList)
			})

		})
	})
}
