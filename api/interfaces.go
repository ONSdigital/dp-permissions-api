package api

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-permissions-api/models"
)

//go:generate moq -out mock/permissionsStore.go -pkg mock . PermissionsStore

//PermissionsStore defines the behaviour of a PermissionsStore
type PermissionsStore interface {
	Checker(ctx context.Context, state *healthcheck.CheckState) error
	Close(ctx context.Context) error
	GetRole(ctx context.Context, id string) (*models.Role, error)
	GetRoles(ctx context.Context) ([]models.Role, error)
}
