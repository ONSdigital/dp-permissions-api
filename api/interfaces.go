package api

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-permissions-api/models"
)

//go:generate moq -out mock/permissionsStore.go -pkg mock . PermissionsStore
//go:generate moq -out mock/bundler.go -pkg mock . PermissionsBundler

//PermissionsStore defines the behaviour of a PermissionsStore
type PermissionsStore interface {
	Checker(ctx context.Context, state *healthcheck.CheckState) error
	Close(ctx context.Context) error
	GetRole(ctx context.Context, id string) (*models.Role, error)
	GetRoles(ctx context.Context, offset, limit int) (*models.Roles, error)
	AddPolicy(ctx context.Context, policy *models.Policy) (*models.Policy, error)
	GetPolicy(ctx context.Context, id string) (*models.Policy, error)
}

// PermissionsBundler defines the functions used by the API to get permissions bundles
type PermissionsBundler interface {
	Get(ctx context.Context) (models.Bundle, error)
}
