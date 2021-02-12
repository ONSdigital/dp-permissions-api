package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/utils"
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

	b, err := json.Marshal(role)
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
	logData := log.Data{}

	offsetParameter := req.URL.Query().Get("offset")
	limitParameter := req.URL.Query().Get("limit")

	offset := api.defaultOffset
	limit := api.defaultLimit
	var err error

	if limitParameter != "" {
		logData["limit"] = limitParameter
		limit, err = utils.ValidatePositiveInteger(limitParameter)
		if err != nil {
			log.Event(ctx, "invalid query parameter: limit", log.ERROR, log.Error(err), logData)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if offsetParameter != "" {
		logData["offset"] = offsetParameter
		offset, err = utils.ValidatePositiveInteger(offsetParameter)
		if err != nil {
			log.Event(ctx, "invalid query parameter: offset", log.ERROR, log.Error(err), logData)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if limit > api.maximumDefaultLimit {
		logData["max_limit"] = api.maximumDefaultLimit
		err = apierrors.ErrorMaximumLimitReached(api.maximumDefaultLimit)
		log.Event(ctx, "invalid query parameter: limit, maximum limit reached", log.ERROR, log.Error(err), logData)
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

	b, err := json.Marshal(listOfRoles)
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
