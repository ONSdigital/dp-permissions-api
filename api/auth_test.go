package api

import (
	"errors"
	"net/http"
	"testing"

	authorisation "github.com/ONSdigital/dp-authorisation/v2/authorisation/mock"
	"github.com/ONSdigital/dp-net/v3/request"
	permissionsAPISDK "github.com/ONSdigital/dp-permissions-api/sdk"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAuthEntityData(t *testing.T) {
	Convey("Given an API with auth middleware", t, func() {
		const token = "test-token"

		authMiddleware := &authorisation.MiddlewareMock{
			ParseFunc: func(token string) (*permissionsAPISDK.EntityData, error) {
				return &permissionsAPISDK.EntityData{
					UserID: "test-user",
					Groups: []string{"role-admin"},
				}, nil
			},
		}

		permissionsAPI := &API{authMiddleware: authMiddleware}
		req, err := http.NewRequest(http.MethodGet, "/v1/policies/test-policy", http.NoBody)
		So(err, ShouldBeNil)
		req.Header.Set(request.AuthHeaderKey, request.BearerPrefix+token)

		Convey("When auth entity data is requested", func() {
			authEntityData, err := permissionsAPI.getAuthEntityData(req)

			Convey("Then the bearer token is parsed", func() {
				So(err, ShouldBeNil)
				So(authMiddleware.ParseCalls(), ShouldHaveLength, 1)
				So(authMiddleware.ParseCalls()[0].Token, ShouldEqual, token)
			})

			Convey("Then user auth entity data is returned", func() {
				So(authEntityData.IsServiceAuth, ShouldBeFalse)
				So(authEntityData.EntityData.UserID, ShouldEqual, "test-user")
				So(authEntityData.EntityData.Groups, ShouldResemble, []string{"role-admin"})
			})
		})
	})
}

func TestGetAuthEntityDataParseError(t *testing.T) {
	Convey("Given auth middleware fails to parse the token", t, func() {
		authMiddleware := &authorisation.MiddlewareMock{
			ParseFunc: func(token string) (*permissionsAPISDK.EntityData, error) {
				return nil, errors.New("invalid token")
			},
		}

		permissionsAPI := &API{authMiddleware: authMiddleware}
		req, err := http.NewRequest(http.MethodGet, "/v1/policies/test-policy", http.NoBody)
		So(err, ShouldBeNil)
		req.Header.Set(request.AuthHeaderKey, request.BearerPrefix+"bad-token")

		Convey("When auth entity data is requested", func() {
			authEntityData, err := permissionsAPI.getAuthEntityData(req)

			Convey("Then an error is returned", func() {
				So(authEntityData, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to parse access token")
			})
		})
	})
}

func TestCreateAuthEntityData(t *testing.T) {
	Convey("Given entity data", t, func() {
		entityData := &permissionsAPISDK.EntityData{
			UserID: "test-user",
			Groups: []string{"role-admin"},
		}

		Convey("When auth entity data is created", func() {
			authEntityData := CreateAuthEntityData(entityData, true)

			Convey("Then values are copied", func() {
				So(authEntityData.IsServiceAuth, ShouldBeTrue)
				So(authEntityData.EntityData.UserID, ShouldEqual, "test-user")
				So(authEntityData.EntityData.Groups, ShouldResemble, []string{"role-admin"})
			})
		})
	})
}
