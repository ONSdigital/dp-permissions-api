package sdk

import (
	"context"

	"github.com/ONSdigital/dp-permissions-api/models"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Clienter

type Clienter interface {
	GetRoles(ctx context.Context) (*models.Roles, error)
	GetRole(ctx context.Context, id string) (*models.Roles, error)
	PostPolicy(ctx context.Context, policy models.PolicyInfo) (*models.Policy, error)
	PostPolicyWithID(ctx context.Context, headers Headers, id string, policy models.PolicyInfo) (*models.Policy, error)
	DeletePolicy(ctx context.Context, id string) error
	GetPolicy(ctx context.Context, id string) (*models.Policy, error)
	PutPolicy(ctx context.Context, id string, policy models.Policy) error
	GetPermissionsBundle(ctx context.Context) (Bundle, error)
}
