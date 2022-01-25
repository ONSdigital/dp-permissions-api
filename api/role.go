package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/dp-permissions-api/utils"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

//GetRoleHandler is a handler that gets a role by its ID from MongoDB
func (api *API) GetRoleHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	roleID := vars["id"]

	//get role from mongoDB by id
	role, err := api.permissionsStore.GetRole(ctx, roleID)
	if err != nil {
		return nil, handleGetRoleError(ctx, err, roleID)
	}

	b, err := json.Marshal(role)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err, "role", role)
	}

	return models.NewSuccessResponse(b, http.StatusOK, nil), nil
}

func handleGetRoleError(ctx context.Context, err error, roleID string) *models.ErrorResponse {
	logData := log.Data{"role_id": roleID}
	if err == apierrors.ErrRoleNotFound {
		return models.NewErrorResponse(http.StatusNotFound,
			nil,
			models.NewError(ctx, err, models.RoleNotFoundError, models.RoleNotFoundDescription, logData),
		)
	}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetRoleError, models.GetRoleErrorDescription, logData),
	)
}

//GetRolesHandler is a handler that gets all roles from MongoDB
func (api *API) GetRolesHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	offsetParameter := req.URL.Query().Get("offset")
	limitParameter := req.URL.Query().Get("limit")

	offset := api.defaultOffset
	limit := api.defaultLimit
	var err error

	if limitParameter != "" {
		limit, err = utils.ValidatePositiveInteger(limitParameter)
		if err != nil {
			return nil, handleInvalidQueryParameterError(ctx, err, "limit", limitParameter)
		}
	}

	if offsetParameter != "" {
		offset, err = utils.ValidatePositiveInteger(offsetParameter)
		if err != nil {
			return nil, handleInvalidQueryParameterError(ctx, err, "offset", offsetParameter)
		}
	}

	if limit > api.maximumDefaultLimit {
		err = apierrors.ErrorMaximumLimitReached(api.maximumDefaultLimit)
		return nil, handleInvalidLimitQueryParameterMaxExceededError(ctx, err, limit, api.maximumDefaultLimit)
	}

	// get roles from MongoDB
	listOfRoles, err := api.permissionsStore.GetRoles(ctx, offset, limit)
	if err != nil {
		return nil, handleGetRolesError(ctx, err)
	}

	b, err := json.Marshal(listOfRoles)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err, "list_of_roles", listOfRoles)
	}

	return models.NewSuccessResponse(b, http.StatusOK, nil), nil
}

func handleInvalidLimitQueryParameterMaxExceededError(ctx context.Context, err error, value int, max int) *models.ErrorResponse {
	logData := log.Data{
		"limit":     value,
		"max_limit": max,
	}
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.InvalidLimitQueryParameterMaxExceededError, models.InvalidLimitQueryParameterMaxExceededDescription, logData),
	)
}

func handleGetRolesError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetRolesError, models.GetRolesErrorDescription, nil),
	)
}
