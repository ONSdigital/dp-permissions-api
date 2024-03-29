package permissions_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/dp-permissions-api/permissions"
	"github.com/ONSdigital/dp-permissions-api/permissions/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBundler_Get(t *testing.T) {
	ctx := context.Background()
	roles := []*models.Role{
		{ID: "admin", Permissions: []string{"legacy.read", "legacy.update", "users.add"}},
		{ID: "publisher", Permissions: []string{"legacy.read", "legacy.update"}},
		{ID: "viewer", Permissions: []string{"legacy.read"}},
	}
	adminPolicy := &models.BundlePolicy{
		Entities: []string{
			"groups/admin",
		},
		Role: "admin",
	}
	publisherPolicy := &models.BundlePolicy{
		Entities: []string{
			"groups/publisher",
		},
		Role: "publisher",
	}
	viewerPolicy := &models.BundlePolicy{
		Entities: []string{
			"groups/viewer",
		},
		Role: "viewer",
		Condition: models.Condition{
			Attribute: "collection-id",
			Operator:  models.OperatorStringEquals,
			Values:    []string{"collection-765"},
		},
	}
	policies := []*models.BundlePolicy{
		adminPolicy,
		publisherPolicy,
		viewerPolicy,
	}
	expectedBundle := models.Bundle{
		"legacy.read": map[string][]*models.BundlePolicy{
			"groups/admin": {
				adminPolicy,
			},
			"groups/publisher": {
				publisherPolicy,
			},
			"groups/viewer": {
				viewerPolicy,
			},
		},
		"legacy.update": map[string][]*models.BundlePolicy{
			"groups/admin": {
				adminPolicy,
			},
			"groups/publisher": {
				publisherPolicy,
			},
		},
		"users.add": map[string][]*models.BundlePolicy{
			"groups/admin": {
				adminPolicy,
			},
		},
	}

	Convey("Given a store that returns permissions data", t, func() {
		store := &mock.StoreMock{
			GetAllBundlePoliciesFunc: func(ctx context.Context) ([]*models.BundlePolicy, error) {
				return policies, nil
			},
			GetAllRolesFunc: func(ctx context.Context) ([]*models.Role, error) {
				return roles, nil
			},
		}
		bundler := permissions.NewBundler(store)

		Convey("When the Get function is called", func() {
			bundle, err := bundler.Get(ctx)

			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the expected bundle is returned", func() {
				So(bundle, ShouldNotBeNil)
				So(bundle, ShouldResemble, expectedBundle)
			})
		})
	})
}

func TestBundler_Get_StoreError(t *testing.T) {
	ctx := context.Background()

	Convey("Given a store that returns an error", t, func() {
		expectedErr := errors.New("store is broken")
		store := &mock.StoreMock{
			GetAllBundlePoliciesFunc: func(ctx context.Context) ([]*models.BundlePolicy, error) {
				return nil, expectedErr
			},
		}
		bundler := permissions.NewBundler(store)

		Convey("When the Get function is called", func() {
			bundle, err := bundler.Get(ctx)

			Convey("Then the expected error should be returned", func() {
				So(err, ShouldEqual, expectedErr)
			})

			Convey("Then the returned bundle should be nil", func() {
				So(bundle, ShouldBeNil)
			})
		})
	})
}
