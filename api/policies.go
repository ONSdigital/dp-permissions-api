package api

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"

	"encoding/json"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
)

func (api *API) PostPolicesHandler(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	policy, err := models.CreateNewPolicy(request.Body)
	if err != nil {
		log.Error(ctx, "unable to unmarshal request body", err)
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
	logData := log.Data{}

	if err := policy.ValidateNewPolicy(); err != nil {
		logData["policies_parameters"] = policy
		log.Error(ctx, "policies parameters failed validation", err, logData)
		return nil, errors.New("invalid request body")
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Error(ctx, "failed to create a new UUID for policies", err, logData)
		return nil, err
	}

	newPolicy.ID = uuid.String()
	newPolicy.Entities = policy.Entities
	newPolicy.Roles = policy.Roles
	newPolicy.Conditions = policy.Conditions
	logData["new_policy"] = newPolicy

	newPolicy, err = api.permissionsStore.AddPolicy(ctx, newPolicy)
	if err != nil {
		log.Error(ctx, "failed to create new policy", err, logData)
		return nil, err
	}
	return newPolicy, nil

}
