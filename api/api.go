package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router              *mux.Router
	permissionsStore    PermissionsStore
	bundler             PermissionsBundler
	defaultLimit        int
	defaultOffset       int
	maximumDefaultLimit int
}

type baseHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.SuccessResponse, *models.ErrorResponse)

func contextAndErrors(h baseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		response, err := h(ctx, w, req)
		if err != nil {
			writeErrorResponse(ctx, w, err)
			return
		}
		writeSuccessResponse(ctx, w, response)
	}
}

//Setup function sets up the api and returns an api
func Setup(
	cfg *config.Config,
	r *mux.Router,
	permissionsStore PermissionsStore,
	bundler PermissionsBundler,
	auth authorisation.Middleware) *API {

	api := &API{
		Router:              r,
		permissionsStore:    permissionsStore,
		defaultLimit:        cfg.DefaultLimit,
		defaultOffset:       cfg.DefaultOffset,
		maximumDefaultLimit: cfg.MaximumDefaultLimit,
		bundler:             bundler,
	}

	r.HandleFunc("/v1/roles", auth.Require(models.RolesRead, api.GetRolesHandler)).Methods(http.MethodGet)
	r.HandleFunc("/v1/roles/{id}", auth.Require(models.RolesRead, api.GetRoleHandler)).Methods(http.MethodGet)
	r.HandleFunc("/v1/policies", auth.Require(models.PoliciesCreate, api.PostPolicyHandler)).Methods(http.MethodPost)
	r.HandleFunc("/v1/policies/{id}", auth.Require(models.PoliciesRead, api.GetPolicyHandler)).Methods(http.MethodGet)
	r.HandleFunc("/v1/policies/{id}", auth.Require(models.PoliciesUpdate, api.UpdatePolicyHandler)).Methods(http.MethodPut)
	r.HandleFunc("/v1/policies/{id}", auth.Require(models.PoliciesDelete, api.DeletePolicyHandler)).Methods(http.MethodDelete)
	r.HandleFunc("/v1/permissions-bundle", api.GetPermissionsBundleHandler).Methods(http.MethodGet)

	return api
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, errorResponse *models.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	// process custom headers
	if errorResponse.Headers != nil {
		for key := range errorResponse.Headers {
			w.Header().Set(key, errorResponse.Headers[key])
		}
	}
	w.WriteHeader(errorResponse.Status)

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.JSONMarshalError, models.ErrorMarshalFailedDescription)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.WriteResponseError, models.WriteResponseFailedDescription)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}

func writeSuccessResponse(ctx context.Context, w http.ResponseWriter, successResponse *models.SuccessResponse) {
	w.Header().Set("Content-Type", "application/json")
	// process custom headers
	if successResponse.Headers != nil {
		for key := range successResponse.Headers {
			w.Header().Set(key, successResponse.Headers[key])
		}
	}
	w.WriteHeader(successResponse.Status)

	_, err := w.Write(successResponse.Body)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.WriteResponseError, models.WriteResponseFailedDescription)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}
