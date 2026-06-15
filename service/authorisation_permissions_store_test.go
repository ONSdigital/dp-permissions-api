package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/dp-permissions-api/sdk"
	. "github.com/smartystreets/goconvey/convey"
)

type permissionsBundlerMock struct {
	getFunc func(ctx context.Context) (models.Bundle, error)
}

func (m *permissionsBundlerMock) Get(ctx context.Context) (models.Bundle, error) {
	return m.getFunc(ctx)
}

func TestAuthorisationPermissionsStore_GetPermissionsBundle(t *testing.T) {
	ctx := context.Background()

	Convey("Given a local permissions store with a bundler that returns a bundle", t, func() {
		bundle := models.Bundle{
			"legacy.read": map[string][]*models.BundlePolicy{
				"groups/viewer": {
					{
						ID:       "policy-id",
						Entities: []string{"groups/viewer"},
						Role:     "viewer",
						Condition: models.Condition{
							Attribute: "collection_id",
							Operator:  models.OperatorStringEquals,
							Values:    []string{"collection-1"},
						},
					},
				},
			},
		}
		expectedBundle := sdk.Bundle{
			"legacy.read": map[string][]sdk.Policy{
				"groups/viewer": {
					{
						ID: "policy-id",
						Condition: sdk.Condition{
							Attribute: "collection_id",
							Operator:  sdk.OperatorStringEquals,
							Values:    []string{"collection-1"},
						},
					},
				},
			},
		}
		bundler := &permissionsBundlerMock{
			getFunc: func(ctx context.Context) (models.Bundle, error) {
				return bundle, nil
			},
		}

		store := newAuthorisationPermissionsStore(bundler)

		Convey("When GetPermissionsBundle is called", func() {
			permissionsBundle, err := store.GetPermissionsBundle(ctx, sdk.Headers{})

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the converted bundle is returned", func() {
				So(permissionsBundle, ShouldResemble, expectedBundle)
			})
		})
	})

	Convey("Given a local permissions store with a bundler that returns an error", t, func() {
		expectedErr := errors.New("bundler error")
		bundler := &permissionsBundlerMock{
			getFunc: func(ctx context.Context) (models.Bundle, error) {
				return nil, expectedErr
			},
		}
		store := newAuthorisationPermissionsStore(bundler)

		Convey("When GetPermissionsBundle is called", func() {
			permissionsBundle, err := store.GetPermissionsBundle(ctx, sdk.Headers{})

			Convey("Then the expected error is returned", func() {
				So(err, ShouldEqual, expectedErr)
			})

			Convey("Then the returned bundle is nil", func() {
				So(permissionsBundle, ShouldBeNil)
			})
		})
	})
}
