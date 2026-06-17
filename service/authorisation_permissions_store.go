package service

import (
	"context"

	authpermissions "github.com/ONSdigital/dp-authorisation/v2/permissions"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/dp-permissions-api/sdk"
	"github.com/ONSdigital/log.go/v2/log"
)

var _ authpermissions.Store = (*authorisationPermissionsStore)(nil)

type permissionsBundler interface {
	Get(ctx context.Context) (models.Bundle, error)
}

type authorisationPermissionsStore struct {
	bundler permissionsBundler
}

func newAuthorisationPermissionsStore(bundler permissionsBundler) authpermissions.Store {
	return &authorisationPermissionsStore{
		bundler: bundler,
	}
}

func (s *authorisationPermissionsStore) GetPermissionsBundle(ctx context.Context, _ sdk.Headers) (sdk.Bundle, error) {
	bundle, err := s.bundler.Get(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(ctx, "retrieved permissions bundle from local store")

	return convertAuthorisationPermissionsBundle(bundle), nil
}

func convertAuthorisationPermissionsBundle(bundle models.Bundle) sdk.Bundle {
	convertedBundle := make(sdk.Bundle, len(bundle))

	for permission, entityPolicies := range bundle {
		convertedEntityPolicies := make(sdk.EntityIDToPolicies, len(entityPolicies))

		for entity, policies := range entityPolicies {
			convertedPolicies := make([]sdk.Policy, 0, len(policies))

			for _, policy := range policies {
				convertedPolicies = append(convertedPolicies, sdk.Policy{
					ID: policy.ID,
					Condition: sdk.Condition{
						Attribute: policy.Condition.Attribute,
						Operator:  sdk.Operator(policy.Condition.Operator),
						Values:    policy.Condition.Values,
					},
				})
			}

			convertedEntityPolicies[entity] = convertedPolicies
		}

		convertedBundle[permission] = convertedEntityPolicies
	}

	return convertedBundle
}
