package api

import (
	"context"
	"net/http"

	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router  *mux.Router
	mongoDB PermissionsStore
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, mongoDB PermissionsStore) *API {
	api := &API{
		Router:  r,
		mongoDB: mongoDB,
	}

	log.Event(ctx, "remove hello endpoint")
	r.HandleFunc("/hello", HelloHandler()).Methods("GET")
	r.HandleFunc("/role/{id}", api.GetRoleHandler).Methods(http.MethodGet)
	return api
}
