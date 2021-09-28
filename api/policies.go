package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gofrs/uuid"
	"net/http"
)

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
