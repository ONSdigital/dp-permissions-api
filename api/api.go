package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// API provides a struct to wrap the api around
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

// Setup function sets up the api and returns an api
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

	r.HandleFunc("/v1/roles", auth.Require(models.RolesRead, contextAndErrors(api.GetRolesHandler))).Methods(http.MethodGet)
	r.HandleFunc("/v1/roles/{id}", auth.Require(models.RolesRead, contextAndErrors(api.GetRoleHandler))).Methods(http.MethodGet)
	r.HandleFunc("/v1/policies", auth.Require(models.PoliciesCreate, contextAndErrors(api.PostPolicyHandler))).Methods(http.MethodPost)
	r.HandleFunc("/v1/policies/{id}", auth.Require(models.PoliciesCreate, contextAndErrors(api.PostPolicyWithIDHandler))).Methods(http.MethodPost)
	r.HandleFunc("/v1/policies/{id}", auth.Require(models.PoliciesRead, contextAndErrors(api.GetPolicyHandler))).Methods(http.MethodGet)
	r.HandleFunc("/v1/policies/{id}", auth.Require(models.PoliciesUpdate, contextAndErrors(api.UpdatePolicyHandler))).Methods(http.MethodPut)
	r.HandleFunc("/v1/policies/{id}", auth.Require(models.PoliciesDelete, contextAndErrors(api.DeletePolicyHandler))).Methods(http.MethodDelete)
	r.HandleFunc("/v1/permissions-bundle", contextAndErrors(api.GetPermissionsBundleHandler)).Methods(http.MethodGet)

	return api
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, errorResponse *models.ErrorResponse) {
	// override internal server error response to prevent leaking sensitive data
	if errorResponse.Status == http.StatusInternalServerError {
		http.Error(w, models.InternalServerErrorDescription, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// process custom headers
	for key, value := range errorResponse.Headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(errorResponse.Status)

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.JSONMarshalError, models.ErrorMarshalFailedDescription, nil)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.WriteResponseError, models.WriteResponseFailedDescription, nil)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}

func writeSuccessResponse(ctx context.Context, w http.ResponseWriter, successResponse *models.SuccessResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// process custom headers
	for key, value := range successResponse.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(successResponse.Status)

	_, err := w.Write(successResponse.Body)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.WriteResponseError, models.WriteResponseFailedDescription, nil)
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}

func handleInvalidQueryParameterError(ctx context.Context, err error, name, value string) *models.ErrorResponse {
	logData := log.Data{name: value}
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.InvalidQueryParameterError, models.InvalidQueryParameterDescription, logData),
	)
}

func handleBodyMarshalError(ctx context.Context, err error, name string, value interface{}) *models.ErrorResponse {
	logData := log.Data{name: value}
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.JSONMarshalError, models.MarshalFailedDescription, logData),
	)
}

func handleBodyUnmarshalError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusBadRequest,
		nil,
		models.NewError(ctx, err, models.JSONUnmarshalError, models.UnmarshalFailedDescription, nil),
	)
}
