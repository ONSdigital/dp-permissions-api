package sdk

import (
	"context"

	"github.com/ONSdigital/dp-permissions-api/models"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Clienter

type Clienter interface {
	GetRoles(ctx context.Context, headers Headers) (*models.Roles, error)
	GetRole(ctx context.Context, id string, headers Headers) (*models.Roles, error)
	PostPolicy(ctx context.Context, policy models.PolicyInfo, headers Headers) (*models.Policy, error)
	PostPolicyWithID(ctx context.Context, id string, policy models.PolicyInfo, headers Headers) (*models.Policy, error)
	DeletePolicy(ctx context.Context, id string, headers Headers) error
	GetPolicy(ctx context.Context, id string, headers Headers) (*models.Policy, error)
	PutPolicy(ctx context.Context, id string, policy models.Policy, headers Headers) error
	GetPermissionsBundle(ctx context.Context, headers Headers) (Bundle, error)
}
