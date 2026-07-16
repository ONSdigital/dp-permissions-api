package api

import (
	"context"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

const unknownUser = "unknown user"

// logAuditEvent produces protective monitoring logging for the API endpoints given successful requests.
// should we be logging a failed request we also log the reason for this.
func logAuditEvent(ctx context.Context, message string, authEntityData *authorisation.AuthEntityData, action models.Action,
	endpoint string, outcome models.Outcome, errReason string) {
	identityType := log.USER

	if authEntityData != nil && authEntityData.IsServiceAuth {
		identityType = log.SERVICE
	}

	data := log.Data{
		"action":   action,
		"endpoint": endpoint,
		"outcome":  outcome,
	}

	if errReason != "" {
		data["reason"] = errReason
	}

	var userID string
	if authEntityData != nil {
		userID = authEntityData.EntityData.UserID
	} else {
		userID = unknownUser
	}

	log.Info(
		ctx,
		message,
		log.Classification(log.ProtectiveMonitoring),
		log.Auth(identityType, userID),
		data,
	)
}
