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
	// logdata := log.Data{"role-id": roleID} TODO: add back in

	//get role from mongoDB by id
	role, err := api.permissionsStore.GetRole(ctx, roleID)
	if err != nil {
		return nil, handleGetRoleError(ctx, err) // TODO: add logdata
	}

	b, err := json.Marshal(role)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err) // TODO: add logdata; error message is no longer role specific
	}

	return models.NewSuccessResponse(b, http.StatusOK, nil), nil // TODO: logdata is not passed to writeErrorResponse; error message is no longer role specific

	// log.Info(ctx, "getRole Handler: Successfully retrieved role", logdata)  // TODO: happy-path success is no longer logged
}

func handleGetRoleError(ctx context.Context, err error) *models.ErrorResponse {
	if err == apierrors.ErrRoleNotFound {
		return models.NewErrorResponse(http.StatusNotFound,
			nil,
			models.NewError(ctx, err, models.RoleNotFoundError, models.RoleNotFoundDescription), // TODO: models.RoleNotFoundDescription duplicates apierrors.ErrRoleNotFound - use apierrors instead?
		)
	}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetRoleError, models.GetRoleErrorDescription),
	)
}

//GetRolesHandler is a handler that gets all roles from MongoDB
func (api *API) GetRolesHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	logData := log.Data{}

	offsetParameter := req.URL.Query().Get("offset")
	limitParameter := req.URL.Query().Get("limit")

	offset := api.defaultOffset
	limit := api.defaultLimit
	var err error

	if limitParameter != "" {
		logData["limit"] = limitParameter
		limit, err = utils.ValidatePositiveInteger(limitParameter)
		if err != nil {
			return nil, handleInvalidQueryParameterError("limit", ctx, err)
		}
	}

	if offsetParameter != "" {
		logData["offset"] = offsetParameter
		offset, err = utils.ValidatePositiveInteger(offsetParameter)
		if err != nil {
			return nil, handleInvalidQueryParameterError("offset", ctx, err) // TODO: parameterised error ok?
		}
	}

	if limit > api.maximumDefaultLimit {
		logData["max_limit"] = api.maximumDefaultLimit
		err = apierrors.ErrorMaximumLimitReached(api.maximumDefaultLimit)
		return nil, handleInvalidLimitQueryParameterMaxExceededError(ctx, err) // TODO: use handleInvalidQueryParameterError() instead?
	}

	// get roles from MongoDB
	listOfRoles, err := api.permissionsStore.GetRoles(ctx, offset, limit)
	if err != nil {
		return nil, handleGetRolesError(ctx, err)
	}

	b, err := json.Marshal(listOfRoles)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err) // TODO: error is no longer roles specific
	}

	//Set headers
	// w.Header().Set("Content-Type", "application/json; charset=utf-8")	// TODO: now handled by writeSuccessResponse but charset is missing there

	return models.NewSuccessResponse(b, http.StatusOK, nil), nil // TODO: errors are no longer role specific

	// log.Info(ctx, "getRoles Handler: Successfully retrieved roles") // TODO: there's no happy-path success logging in dp-identity-api; it's inconsistently applied in this service - remove?
}

func handleInvalidLimitQueryParameterMaxExceededError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.InvalidLimitQueryParameterMaxExceededError, models.InvalidLimitQueryParameterMaxExceededDescription),
	)
}

func handleGetRolesError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetRolesError, models.GetRolesErrorDescription),
	)
}
