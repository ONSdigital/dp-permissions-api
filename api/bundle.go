package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/models"
)

// GetPermissionsBundleHandler gets and returns the permissions bundle as JSON in the HTTP response body.
func (api *API) GetPermissionsBundleHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {

	bundle, err := api.bundler.Get(ctx)
	if err != nil {
		return nil, handleGetPermissionsBundleError(ctx, err)
	}

	b, err := json.Marshal(bundle)
	if err != nil {
		return nil, handleBodyMarshalError(ctx, err) // TODO: log message is no longer specific to permissions bundle
	}

	return models.NewSuccessResponse(b, http.StatusOK, nil), nil
	// TODO: errors during success response write is no longer specific to permissions bundle

	// log.Info(ctx, "successfully retrieved permissions bundle")  // TODO: happy-path success is no longer logged
}

func handleGetPermissionsBundleError(ctx context.Context, err error) *models.ErrorResponse {
	return models.NewErrorResponse(http.StatusInternalServerError,
		nil,
		models.NewError(ctx, err, models.GetPermissionBundleError, models.GetPermissionBundleErrorDescription),
	)
}
