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

//GetPolicyHandler is a handler that gets policy by its ID from DB
func (api *API) GetPolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	vars := mux.Vars(req)
	policyId := vars["id"]
	// logData := log.Data{"policy_id": policyId} // TODO: add to error logs

	policy, err := api.permissionsStore.GetPolicy(ctx, policyId)
	if err != nil {
		return nil, handleGetPolicyError(ctx, err) // TODO: add logData
	}

	b, err := json.Marshal(policy)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err) // TODO: add logData; log message is no longer policy specific
	}

	return models.NewSuccessResponse(b, http.StatusOK, nil), nil // TODO: add logData to error logs

	// log.Info(ctx, "getPolicy Handler: Successfully retrieved policy", logData)  // TODO: no longer log happy-path success?
}

func handleGetPolicyError(ctx context.Context, err error) *models.ErrorResponse {
	if err == apierrors.ErrPolicyNotFound {
		return models.NewErrorResponse(http.StatusNotFound,
			nil,
			models.NewError(ctx, err, models.PolicyNotFoundError, models.PolicyNotFoundDescription),
		)
	}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetPolicyError, models.GetPolicyErrorDescription),
	)
}

//DeletePolicyHandler is a handler that deletes policy by its ID from DB
func (api *API) DeletePolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	vars := mux.Vars(req)
	policyId := vars["id"]
	// logData := log.Data{"policy_id": policyId}  // TODO: add to logging

	err := api.permissionsStore.DeletePolicy(ctx, policyId)
	if err != nil {
		return nil, handleDeletePolicyError(ctx, err) // TODO: add logData to logging
	}

	return models.NewSuccessResponse(nil, http.StatusNoContent, nil), nil // TODO: logged errors are not policy specific
	// log.Info(ctx, "deletePolicy Handler: Successfully deleted policy", logData)  // TODO: remove happy-path success logging?
}

func handleDeletePolicyError(ctx context.Context, err error) *models.ErrorResponse {
	if err == apierrors.ErrPolicyNotFound {
		return models.NewErrorResponse(http.StatusNotFound,
			nil,
			models.NewError(ctx, err, models.PolicyNotFoundError, models.PolicyNotFoundDescription), // TODO: models.PolicyNotFoundDescription duplicates apierrors.ErrPolicyNotFound, just use the apierror?
		)
	}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.DeletePolicyError, models.DeletePolicyErrorDescription),
	)
}

//PostPolicyHandler is a handler that creates a new policies in DB
func (api *API) PostPolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	policy, err := models.CreatePolicy(req.Body)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	if err := policy.ValidatePolicy(); err != nil {
		// logData := log.Data{}				// TODO: add logData to logs
		// logData["policies_parameters"] = policy
		return nil, handleValidatePolicyError(ctx, err)
	}

	newPolicy, err := api.createNewPolicy(ctx, policy)
	if err != nil {
		return nil, handleCreateNewPolicyError(ctx, err)
	}

	b, err := json.Marshal(newPolicy)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err) // TODO: error description no longer policy specific
	}

	return models.NewSuccessResponse(b, http.StatusCreated, nil), nil
	// TODO: inconsistent happy-path success logging (not included here)
}

func handleValidatePolicyError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.InvalidPolicyError, models.InvalidPolicyDescription+err.Error()), // TODO: is this the best error description?
	)
}

func (api *API) createNewPolicy(ctx context.Context, policy *models.PolicyInfo) (*models.Policy, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Error(ctx, "failed to create a new UUID for policies", err)
		return nil, err
	}

	newPolicy, err := api.permissionsStore.AddPolicy(ctx, policy.GetPolicy(uuid.String()))
	if err != nil {
		return nil, err
	}
	return newPolicy, nil
}

func handleCreateNewPolicyError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.CreateNewPolicyError, models.CreateNewPolicyErrorDescription),
	)
}

//UpdatePolicyHandler is a handler that updates policy by its ID from DB
func (api *API) UpdatePolicyHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	vars := mux.Vars(req)
	policyId := vars["id"]
	// logData := log.Data{"policy_id": policyId} // TODO: re-add log data

	updatePolicy, err := models.CreatePolicy(req.Body)
	if err != nil {
		return nil, handleBodyUnmarshalError(ctx, err)
	}

	if err := updatePolicy.ValidatePolicy(); err != nil {
		// logData := log.Data{}
		// logData["policies_parameters"] = updatePolicy
		return nil, handleValidatePolicyError(ctx, err) // TODO: add log data
	}

	updateResult, err := api.permissionsStore.UpdatePolicy(ctx, updatePolicy.GetPolicy(policyId))
	// writer.Header().Set("Content-Type", "application/json; charset=utf-8")  // TODO: delete this, bug?
	if err != nil {
		return nil, handleUpdatePolicyError(ctx, err) // TODO: add logData to log
	}

	if updateResult.ModifiedCount > 0 {
		// log.Info(ctx, "Updated policy", logData) // TODO: remove happy-path success log?
		return models.NewSuccessResponse(nil, http.StatusOK, nil), nil
	} else {
		// log.Info(ctx, "Created new policy", logData) // TODO: remove happy-path success log?
		return models.NewSuccessResponse(nil, http.StatusCreated, nil), nil
	}

	// log.Info(ctx, "UpdatePolicy Handler: Successfully upserted policy", logData)  // TODO: remove happy path success log?
}

func handleUpdatePolicyError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.UpdatePolicyError, models.UpdatePolicyErrorDescription),
	)
}
