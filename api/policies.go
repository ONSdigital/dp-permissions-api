package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

// GetPolicyHandler is a handler that gets policy by its ID from DB
func (api *API) GetPolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	policyID := vars["id"]

	policy, err := api.permissionsStore.GetPolicy(ctx, policyID)
	if err != nil {
		return nil, handleGetPolicyError(ctx, err, policyID)
	}

	b, err := json.Marshal(policy)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err, "policy", policy)
	}

	return models.NewSuccessResponse(b, http.StatusOK, nil), nil
}

func handleGetPolicyError(ctx context.Context, err error, policyID string) *models.ErrorResponse {
	logData := log.Data{"policy_id": policyID}
	if err == apierrors.ErrPolicyNotFound {
		return models.NewErrorResponse(http.StatusNotFound,
			nil,
			models.NewError(ctx, err, models.PolicyNotFoundError, models.PolicyNotFoundDescription, logData),
		)
	}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetPolicyError, models.GetPolicyErrorDescription, logData),
	)
}

// DeletePolicyHandler is a handler that deletes policy by its ID from DB
func (api *API) DeletePolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	policyID := vars["id"]

	err := api.permissionsStore.DeletePolicy(ctx, policyID)
	if err != nil {
		return nil, handleDeletePolicyError(ctx, err, policyID)
	}

	return models.NewSuccessResponse(nil, http.StatusNoContent, nil), nil
}

func handleDeletePolicyError(ctx context.Context, err error, policyID string) *models.ErrorResponse {
	logData := log.Data{"policy_id": policyID}
	if err == apierrors.ErrPolicyNotFound {
		return models.NewErrorResponse(http.StatusNotFound,
			nil,
			models.NewError(ctx, err, models.PolicyNotFoundError, models.PolicyNotFoundDescription, logData),
		)
	}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.DeletePolicyError, models.DeletePolicyErrorDescription, logData),
	)
}

// PostPolicyHandler is a handler that creates a new policies in DB
func (api *API) PostPolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	policy, err := models.CreatePolicy(req.Body)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	if err := policy.ValidatePolicy(); err != nil {
		return nil, handleValidatePolicyError(ctx, err, policy)
	}

	newPolicy, err := api.createNewPolicy(ctx, policy)
	if err != nil {
		return nil, handleCreateNewPolicyError(ctx, err)
	}

	b, err := json.Marshal(newPolicy)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err, "new_policy", newPolicy)
	}

	return models.NewSuccessResponse(b, http.StatusCreated, nil), nil
}

func handleValidatePolicyError(ctx context.Context, err error, policy *models.PolicyInfo) *models.ErrorResponse {
	logData := log.Data{}
	logData["policies_parameters"] = *policy
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.InvalidPolicyError, err.Error(), logData),
	)
}

func (api *API) createNewPolicy(ctx context.Context, policy *models.PolicyInfo) (*models.Policy, error) {
	policyuuid, err := uuid.NewV4()
	if err != nil {
		log.Error(ctx, "failed to create a new UUID for policies", err)
		return nil, err
	}

	newPolicy, err := api.permissionsStore.AddPolicy(ctx, policy.GetPolicy(policyuuid.String()))
	if err != nil {
		return nil, err
	}
	return newPolicy, nil
}

func handleCreateNewPolicyError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.CreateNewPolicyError, models.CreateNewPolicyErrorDescription, nil),
	)
}

// UpdatePolicyHandler is a handler that updates policy by its ID from DB
func (api *API) UpdatePolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	policyID := vars["id"]

	updatePolicy, err := models.CreatePolicy(req.Body)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	if err := updatePolicy.ValidatePolicy(); err != nil {
		return nil, handleValidatePolicyError(ctx, err, updatePolicy)
	}

	updateResult, err := api.permissionsStore.UpdatePolicy(ctx, updatePolicy.GetPolicy(policyID))
	if err != nil {
		return nil, handleUpdatePolicyError(ctx, err, policyID)
	}

	if updateResult.ModifiedCount > 0 {
		return models.NewSuccessResponse(nil, http.StatusOK, nil), nil
	} else {
		return models.NewSuccessResponse(nil, http.StatusCreated, nil), nil
	}
}

func handleUpdatePolicyError(ctx context.Context, err error, policyID string) *models.ErrorResponse {
	logData := log.Data{"policy_id": policyID}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.UpdatePolicyError, models.UpdatePolicyErrorDescription, logData),
	)
}
