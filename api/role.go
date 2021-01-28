package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/models"

	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

//GetRoleHandler is a handler that gets a role by its ID from MongoDB
func (api *API) GetRoleHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	roleID := vars["id"]
	logdata := log.Data{"role-id": roleID}

	//get role from mongoDB by id
	role, err := api.permissionsStore.GetRole(ctx, roleID)
	if err != nil {
		log.Event(ctx, "getRole Handler: retrieving role from mongoDB returned an error", log.ERROR, log.Error(err), logdata)
		if err == apierrors.ErrRoleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var b []byte
	b, err = json.Marshal(role)
	if err != nil {
		log.Event(ctx, "getRole Handler: failed to marshal role resource into bytes", log.ERROR, log.Error(err), logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(b); err != nil {
		log.Event(ctx, "getRole Handler: error writing bytes to response", log.ERROR, log.Error(err), logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Event(ctx, "getRole Handler: Successfully retrieved role", log.INFO, logdata)

}

//GetRolesHandler is a handler that gets all roles from MongoDB
func (api *API) GetRolesHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	log.Event(ctx, "getRoles Handler: came in here", log.INFO)
	// get roles from MongoDB
	listOfRoles, err := api.permissionsStore.GetRoles(ctx)
	if err != nil {
		log.Event(ctx, "getRoles Handler: retrieving roles from MongoDB returned an error", log.ERROR, log.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logdata := log.Data{"role-id": listOfRoles}
	log.Event(ctx, "getRoles Handler: came in here", logdata)

	roles := models.Roles{
		Items:      listOfRoles,
		Count:      len(listOfRoles),
		Limit:      len(listOfRoles),
		TotalCount: len(listOfRoles),
	}

	var b []byte
	b, err = json.Marshal(roles)
	if err != nil {
		log.Event(ctx, "getRoles Handler: failed to marshal roles resource into bytes", log.ERROR, log.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(b); err != nil {
		log.Event(ctx, "getRoles Handler: error writing bytes to response", log.ERROR, log.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Event(ctx, "getRoles Handler: Successfully retrieved roles", log.INFO)

}
