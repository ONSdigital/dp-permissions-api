package service

import (
	"context"

	"github.com/ONSdigital/dp-permissions-api/permissions"

	"github.com/ONSdigital/dp-permissions-api/api"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
)

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	Config                  *config.Config
	Server                  HTTPServer
	Router                  *mux.Router
	api                     *api.API
	ServiceList             *ExternalServiceList
	HealthCheck             HealthChecker
	MongoDB                 PermissionsStore
	AuthorisationMiddleware authorisation.Middleware
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {

	log.Info(ctx, "running service")

	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	// Get HTTP Server and ... // ADD CODE: Add any middleware that your service requires
	r := mux.NewRouter()

	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	// Get MongoDB client
	mongoDB, err := serviceList.GetMongoDB(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise mongo DB", err)
		return nil, err
	}

	bundler := permissions.NewBundler(mongoDB)

	authorisationMiddleware, err := serviceList.GetAuthorisationMiddleware(ctx, cfg.AuthorisationConfig)
	if err != nil {
		log.Fatal(ctx, "could not instantiate authorisation middleware", err)
		return nil, err
	}

	// Setup the API
	a := api.Setup(cfg, r, mongoDB, bundler, authorisationMiddleware)

	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)

	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, hc, mongoDB, authorisationMiddleware); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		Config:                  cfg,
		Router:                  r,
		api:                     a,
		HealthCheck:             hc,
		ServiceList:             serviceList,
		Server:                  s,
		MongoDB:                 mongoDB,
		AuthorisationMiddleware: authorisationMiddleware,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		if svc.ServiceList.MongoDB {
			if err := svc.MongoDB.Close(ctx); err != nil {
				log.Error(ctx, "error closing mongo db", err)
				hasShutdownError = true
			}
		}

		if svc.ServiceList.AuthorisationMiddleware {
			if err := svc.AuthorisationMiddleware.Close(ctx); err != nil {
				log.Error(ctx, "failed to close authorisation middleware", err)
				hasShutdownError = true
			}
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context,
	hc HealthChecker,
	permissionsStore PermissionsStore,
	authorisationMiddleware authorisation.Middleware) (err error) {

	hasErrors := false

	if err = hc.AddCheck("Mongo DB", permissionsStore.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for mongo db", err)
	}

	if err := hc.AddCheck("permissions cache health check", authorisationMiddleware.HealthCheck); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for permissions cache", err)
	}

	if err := hc.AddCheck("jwt keys state health check", authorisationMiddleware.IdentityHealthCheck); err != nil {
		hasErrors = true
		log.Error(ctx, "error getting jwt keys from identity service", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for health check")
	}
	return nil
}
