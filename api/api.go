package api

import (
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/config"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"

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
	r.HandleFunc("/v1/permissions-bundle", api.GetPermissionsBundleHandler).Methods(http.MethodGet)

	return api
}
