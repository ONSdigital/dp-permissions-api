package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-net/v3/request"
	permissionsAPISDK "github.com/ONSdigital/dp-permissions-api/sdk"
)

// AuthEntityData holds the entity data for an authenticated request along with
// whether the request was made by a service account or user
type AuthEntityData struct {
	EntityData    *permissionsAPISDK.EntityData
	IsServiceAuth bool
}

// getAuthEntityData returns the AuthEntityData associated with the provided access token.
func (api *API) getAuthEntityData(r *http.Request) (*AuthEntityData, error) {
	accessToken := strings.TrimPrefix(r.Header.Get(request.AuthHeaderKey), request.BearerPrefix)
	entityData, err := api.authMiddleware.Parse(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}
	return CreateAuthEntityData(entityData, false), nil
}

// CreateAuthEntityData creates an AuthEntityData from the provided EntityData and
// a bool indicating whether the token belongs to a service account
func CreateAuthEntityData(entityData *permissionsAPISDK.EntityData, isService bool) *AuthEntityData {
	return &AuthEntityData{
		EntityData: &permissionsAPISDK.EntityData{
			UserID: entityData.UserID,
			Groups: entityData.Groups,
		},
		IsServiceAuth: isService,
	}
}
