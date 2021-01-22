package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

/*
//GetRolesHandler is a handler that gets all roles from MongoDB
func (api *API) GetRolesHandler(w http.ResponseWriter, req *http.Request, mongoConf config.MongoConfiguration) {
	ctx := req.Context()

	// get roles from MongoDB
	items, err := api.mongoDB.GetRoles(ctx, mongoConf)
	if err != nil {
		log.Event(ctx, "api endpoint getRoles returned an error", log.ERROR, log.Error(err))
		return nil, err
	}

	roles := models.Roles{
		Items:      items,
		Count:      len(items),
		Limit:      len(items),
		TotalCount: len(items),
	}

	var b []byte
	b, err = json.Marshal(roles)

}
*/

//GetRoleHandler is a handler that gets a role by its ID from MongoDB
func (api *API) GetRoleHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	roleID := vars["id"]
	logdata := log.Data{"role-id": roleID}

	//get role from mongoDB by id
	role, err := api.mongoDB.GetRole(ctx, roleID)
	if err != nil {
		log.Event(ctx, "getRole Handler: retrieving role from mongoDB returned an error", log.ERROR, log.Error(err), logdata)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var b []byte
	b, err = json.Marshal(role)
	if err != nil {
		log.Event(ctx, "getRole Handler: filed to marshal role resource into bytes", log.ERROR, log.Error(err), logdata)
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
