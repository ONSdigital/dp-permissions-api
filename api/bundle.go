package api

import (
	"encoding/json"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
)

func (api *API) GetPermissionsBundleHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	bundle, err := api.bundler.Get(ctx)
	if err != nil {
		log.Error(ctx, "failed to get permissions bundle", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(bundle)
	if err != nil {
		log.Error(ctx, "failed to marshal permissions bundle to json", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(b); err != nil {
		log.Error(ctx, "error writing permissions bundle response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info(ctx, "successfully retrieved permissions bundle")
}
