package api

import (
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/config"

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
	bundler PermissionsBundler) *API {

	api := &API{
		Router:              r,
		permissionsStore:    permissionsStore,
		defaultLimit:        cfg.DefaultLimit,
		defaultOffset:       cfg.DefaultOffset,
		maximumDefaultLimit: cfg.MaximumDefaultLimit,
		bundler:             bundler,
	}

	r.HandleFunc("/roles/{id}", api.GetRoleHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/roles", api.GetRolesHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/policies", api.PostPolicyHandler).Methods(http.MethodPost)
	r.HandleFunc("/v1/policies/{id}", api.GetPolicyHandler).Methods(http.MethodGet)
	r.HandleFunc("/v1/permissions-bundle", api.GetPermissionsBundleHandler).Methods(http.MethodGet)

	return api
}
