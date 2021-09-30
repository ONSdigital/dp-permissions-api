package permissions

import (
	"context"
	"github.com/ONSdigital/dp-permissions-api/models"
)

//go:generate moq -out mock/store.go -pkg mock . Store

//Store defines the behaviour of a PermissionsStore as used by the Bundler type.
type Store interface {
	GetAllRoles(ctx context.Context) ([]*models.Role, error)
	GetAllPolicies(ctx context.Context) ([]*models.Policy, error)
}

// Bundler creates permission bundle data - a format optimised for evaluating user permissions.
type Bundler struct {
	store Store
}

// NewBundler creates a new Bundler instance.
func NewBundler(store Store) *Bundler {
	return &Bundler{
		store: store,
	}
}

// Get the latest bundle data.
func (b Bundler) Get(ctx context.Context) (models.Bundle, error) {
	policies, err := b.store.GetAllPolicies(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := b.store.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	bundle := createBundle(policies, roles)

	return bundle, nil
}

func createBundle(policies []*models.Policy, roles []*models.Role) models.Bundle {
	roleIDToPolicies := createRoleToPoliciesMap(policies)
	bundle := models.Bundle{}

	for _, role := range roles {

		policiesForRole := roleIDToPolicies[role.ID]

		for _, permission := range role.Permissions {

			entityLookup, ok := bundle[permission]
			if !ok {
				entityLookup = map[string][]*models.Policy{}
				bundle[permission] = entityLookup
			}

			for _, policy := range policiesForRole {
				for _, entity := range policy.Entities {
					if entityLookup[entity] == nil {
						entityLookup[entity] = []*models.Policy{}
					}

					entityLookup[entity] = append(entityLookup[entity], policy)
				}
			}
		}
	}
	return bundle
}

func createRoleToPoliciesMap(policies []*models.Policy) map[string][]*models.Policy {
	roleIDToPolicies := map[string][]*models.Policy{}
	for _, policy := range policies {
		roleIDToPolicies[policy.Role] = append(roleIDToPolicies[policy.Role], policy)
	}
	return roleIDToPolicies
}
