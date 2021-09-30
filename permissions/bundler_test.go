package permissions_test

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/dp-permissions-api/permissions"
	"github.com/ONSdigital/dp-permissions-api/permissions/mock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBundler_Get(t *testing.T) {
	ctx := context.Background()
	roles := []*models.Role{
		{ID: "admin", Permissions: []string{"legacy.read", "legacy.update", "users.add"}},
		{ID: "publisher", Permissions: []string{"legacy.read", "legacy.update"}},
		{ID: "viewer", Permissions: []string{"legacy.read"}},
	}
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
	viewerPolicy := &models.Policy{
		Entities: []string{
			"groups/viewer",
		},
		Role: "viewer",
		Conditions: []models.Condition{
			{
				Attributes: []string{"collection-id"},
				Operator:   "=",
				Values:     []string{"collection-765"},
			},
		},
	}
	policies := []*models.Policy{
		adminPolicy,
		publisherPolicy,
		viewerPolicy,
	}
	expectedBundle := models.Bundle{
		"legacy.read": map[string][]*models.Policy{
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
		"legacy.update": map[string][]*models.Policy{
			"groups/admin": {
				adminPolicy,
			},
			"groups/publisher": {
				publisherPolicy,
			},
		},
		"users.add": map[string][]*models.Policy{
			"groups/admin": {
				adminPolicy,
			},
		},
	}

	Convey("Given a store that returns permissions data", t, func() {
		store := &mock.StoreMock{
			GetAllPoliciesFunc: func(ctx context.Context) ([]*models.Policy, error) {
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
			GetAllPoliciesFunc: func(ctx context.Context) ([]*models.Policy, error) {
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
