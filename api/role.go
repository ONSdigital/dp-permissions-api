package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/utils"
	"github.com/ONSdigital/log.go/v2/log"
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
		log.Error(ctx, "getRole Handler: retrieving role from mongoDB returned an error", err, logdata)
		if err == apierrors.ErrRoleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(role)
	if err != nil {
		log.Error(ctx, "getRole Handler: failed to marshal role resource into bytes", err, logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(b); err != nil {
		log.Error(ctx, "getRole Handler: error writing bytes to response", err, logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info(ctx, "getRole Handler: Successfully retrieved role", logdata)

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
			log.Error(ctx, "invalid query parameter: limit", err, logData)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if offsetParameter != "" {
		logData["offset"] = offsetParameter
		offset, err = utils.ValidatePositiveInteger(offsetParameter)
		if err != nil {
			log.Error(ctx, "invalid query parameter: offset", err, logData)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if limit > api.maximumDefaultLimit {
		logData["max_limit"] = api.maximumDefaultLimit
		err = apierrors.ErrorMaximumLimitReached(api.maximumDefaultLimit)
		log.Error(ctx, "invalid query parameter: limit, maximum limit reached", err, logData)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get roles from MongoDB
	listOfRoles, err := api.permissionsStore.GetRoles(ctx, offset, limit)

	if err != nil {
		log.Error(ctx, "getRoles Handler: retrieving roles from MongoDB returned an error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(listOfRoles)
	if err != nil {
		log.Error(ctx, "getRoles Handler: failed to marshal roles resource into bytes", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(b); err != nil {
		log.Error(ctx, "getRoles Handler: error writing bytes to response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info(ctx, "getRoles Handler: Successfully retrieved roles")

}
