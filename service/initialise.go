package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/mongo"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck bool
	Init        Initialiser
	MongoDB     bool
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck: false,
		Init:        initialiser,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// GetMongoDB creates a mongoDB client and sets the Mongo flag to true
func (e *ExternalServiceList) GetMongoDB(ctx context.Context, cfg *config.Config) (PermissionsStore, error) {
	mongoDB, err := e.Init.DoGetMongoDB(ctx, cfg)
	if err != nil {
		return nil, err
	}
	e.MongoDB = true
	return mongoDB, nil
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

// DoGetMongoDB returns a MongoDB
func (e *Init) DoGetMongoDB(ctx context.Context, cfg *config.Config) (PermissionsStore, error) {
	mongoDB, err := mongo.NewMongoStore(ctx, cfg.MongoDB)
	if err != nil {
		return nil, err
	}

	return mongoDB, nil
}

// DoGetAuthorisationMiddleware creates authorisation middleware for the given config
func (e *Init) DoGetAuthorisationMiddleware(ctx context.Context, authorisationConfig *authorisation.Config) (authorisation.Middleware, error) {
	return authorisation.NewFeatureFlaggedMiddleware(ctx, authorisationConfig, nil)
}

// GetAuthorisationMiddleware creates a new instance of authorisation.Middlware
func (e *ExternalServiceList) GetAuthorisationMiddleware(ctx context.Context, authorisationConfig *authorisation.Config) (authorisation.Middleware, error) {
	return e.Init.DoGetAuthorisationMiddleware(ctx, authorisationConfig)
}
