package api

import (
	"context"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

// logPolicyAuditEvent produces protective monitoring logging for the API endpoints given successful requests.
// should we be logging a failed request we also log the reason for this.
func logPolicyAuditEvent(ctx context.Context, message string, authEntityData *authorisation.AuthEntityData, action models.Action,
	endpoint string, outcome models.Outcome, errReason string) {
	identityType := log.USER
	if authEntityData.IsServiceAuth {
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

	log.Info(
		ctx,
		message,
		log.Classification(log.ProtectiveMonitoring),
		log.Auth(identityType, authEntityData.EntityData.UserID),
		data,
	)
}
