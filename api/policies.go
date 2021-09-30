package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

//GetPolicyHandler is a handler that gets policy by its ID from DB
func (api *API) GetPolicyHandler(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	vars := mux.Vars(request)
	policyId := vars["id"]
	logData := log.Data{"policy-id": policyId}

	policy, err := api.permissionsStore.GetPolicy(ctx, policyId)
	if err != nil {
		log.Error(ctx, "getPolicy Handler: retrieving policy from DB returned an error", err, logData)
		if err == apierrors.ErrPolicyNotFound {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(policy)
	if err != nil {
		log.Error(ctx, "getPolicy Handler: failed to marshal policy resource into bytes", err, logData)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := writer.Write(b); err != nil {
		log.Error(ctx, "getPolicy Handler: error writing bytes to response", err, logData)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info(ctx, "getPolicy Handler: Successfully retrieved policy", logData)
}

//PostPolicyHandler is a handler that creates a new policies in DB
func (api *API) PostPolicyHandler(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	policy, err := models.CreateNewPolicy(request.Body)
	if err != nil {
		log.Error(ctx, "unable to unmarshal request body", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := policy.ValidateNewPolicy(); err != nil {
		logData := log.Data{}
		logData["policies_parameters"] = policy
		log.Error(ctx, "policies parameters failed validation", err, logData)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	newPolicy, err := api.createNewPolicy(ctx, policy)
	if err != nil {
		log.Error(ctx, "failed to create new policy", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(newPolicy)
	if err != nil {
		log.Error(ctx, "failed to marshal new policy", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_, err = writer.Write(bytes)
	if err != nil {
		log.Error(ctx, "failed to write bytes for http response", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) createNewPolicy(ctx context.Context, policy *models.NewPolicy) (*models.Policy, error) {
	newPolicy := &models.Policy{}
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Error(ctx, "failed to create a new UUID for policies", err)
		return nil, err
	}

	newPolicy.ID = uuid.String()
	newPolicy.Entities = policy.Entities
	newPolicy.Role = policy.Role
	newPolicy.Conditions = policy.Conditions

	newPolicy, err = api.permissionsStore.AddPolicy(ctx, newPolicy)
	if err != nil {
		return nil, err
	}
	return newPolicy, nil

}
