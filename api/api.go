package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/config"

	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router           *mux.Router
	permissionsStore PermissionsStore
	defaultLimit     int
	defaultOffset    int
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg *config.Config, r *mux.Router, permissionsStore PermissionsStore) *API {
	api := &API{
		Router:           r,
		permissionsStore: permissionsStore,
		defaultLimit:     cfg.MongoConfig.Limit,
		defaultOffset:    cfg.MongoConfig.Offset,
	}

	r.HandleFunc("/roles/{id}", api.GetRoleHandler).Methods(http.MethodGet)
	r.HandleFunc("/roles", api.GetRolesHandler).Methods(http.MethodGet)

	return api
}
