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
func (b Bundler) Get(ctx context.Context) (*models.Bundle, error) {
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

func createBundle(policies []*models.Policy, roles []*models.Role) *models.Bundle {
	roleIDToPolicies := createRoleToPoliciesMap(policies)
	bundle := &models.Bundle{PermissionToEntityLookup: map[string]models.EntityIDToPolicies{}}

	for _, role := range roles {

		policiesForRole := roleIDToPolicies[role.ID]

		for _, permission := range role.Permissions {

			entityLookup, ok := bundle.PermissionToEntityLookup[permission]
			if !ok {
				entityLookup = map[string][]*models.Policy{}
				bundle.PermissionToEntityLookup[permission] = entityLookup
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

//func mapToOptimisedLookupFormat(permission string, roles []Role, policies []Policy) PermissionLookup {
//	lookup := PermissionLookup{entityToPolicies: map[string][]Policy{}}
//
//	// work backwards from the permission -> role -> policy, and add the conditions for each entity associated with the permission
//	for _, role := range roles {
//		if role.HasPermission(permission) {
//			for _, policy := range policies {
//				if policy.HasRole(role.ID) {
//					for _, entity := range policy.Entities {
//
//						if lookup.entityToPolicies[entity] == nil {
//							lookup.entityToPolicies[entity] = []Policy{}
//						}
//
//						lookup.entityToPolicies[entity] = append(lookup.entityToPolicies[entity], policy)
//					}
//				}
//			}
//		}
//	}
//	return lookup
//}
