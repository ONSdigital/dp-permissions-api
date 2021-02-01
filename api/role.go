package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

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

	//get limit from query parameters, or default value
	limit, err := getPaginationQueryParameter(req.URL.Query(), "limit", api.defaultLimit)
	if err != nil {
		log.Event(ctx, "failed to obtain limit from request query parameters", log.ERROR)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//get offset from query parameters, or default value
	offset, err := getPaginationQueryParameter(req.URL.Query(), "offset", api.defaultOffset)
	if err != nil {
		log.Event(ctx, "failed to obtain limit from request query parameters", log.ERROR)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get roles from MongoDB
	listOfRoles, err := api.permissionsStore.GetRoles(ctx, offset, limit)

	if err != nil {
		log.Event(ctx, "getRoles Handler: retrieving roles from MongoDB returned an error", log.ERROR, log.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var b []byte
	b, err = json.Marshal(listOfRoles)
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

func getPaginationQueryParameter(queryVars url.Values, varKey string, defaultValue int) (val int, err error) {
	strVal, found := queryVars[varKey]
	if !found {
		return defaultValue, nil
	}
	val, err = strconv.Atoi(strVal[0])
	if err != nil {
		return -1, apierrors.ErrInvalidQueryParameter
	}
	if val < 0 {
		return 0, nil
	}
	return val, nil
}
