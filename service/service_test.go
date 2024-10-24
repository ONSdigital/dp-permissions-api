package service_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/service"
	"github.com/ONSdigital/dp-permissions-api/service/mock"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"

	authorisationMock "github.com/ONSdigital/dp-authorisation/v2/authorisation/mock"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
	errServer     = errors.New("HTTP Server error")
)

var (
	errMongoDB     = errors.New("mongoDB error")
	errHealthcheck = errors.New("healthCheck error")
)

var funcDoGetMongoDBErr = func(ctx context.Context, cfg *config.Config) (service.PermissionsStore, error) {
	return nil, errMongoDB
}

var funcDoGetHealthcheckErr = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

var funcDoGetHTTPServerNil = func(bindAddr string, router http.Handler) service.HTTPServer {
	return nil
}

func TestRun(t *testing.T) {
	Convey("Having a set of mocked dependencies", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}

		failingServerMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return errServer
			},
		}

		funcDoGetHealthcheckOk := func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		funcDoGetHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		funcDoGetFailingHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return failingServerMock
		}

		funcDoGetMongoDBOk := func(ctx context.Context, cfg *config.Config) (service.PermissionsStore, error) {
			return &mock.PermissionsStoreMock{
				CloseFunc: func(ctx context.Context) error { return nil },
			}, nil
		}

		funcDoGetAuthorisationMiddleware := func(ctx context.Context, authorisationConfig *authorisation.Config) (authorisation.Middleware, error) {
			return &authorisationMock.MiddlewareMock{
				RequireFunc: func(permission string, handlerFunc http.HandlerFunc) http.HandlerFunc {
					return handlerFunc
				},
			}, nil
		}

		Convey("Given that initialising mongoDB returns an error", func() {
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc:              funcDoGetHTTPServerNil,
				DoGetMongoDBFunc:                 funcDoGetMongoDBErr,
				DoGetHealthCheckFunc:             funcDoGetHealthcheckOk,
				DoGetAuthorisationMiddlewareFunc: funcDoGetAuthorisationMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set. No further initialisations are attempted", func() {
				So(err, ShouldResemble, errMongoDB)
				So(svcList.MongoDB, ShouldBeFalse)
				So(svcList.HealthCheck, ShouldBeFalse)
			})
		})

		Convey("Given that initialising healthcheck returns an error", func() {
			// setup (run before each `Convey` at this scope / indentation):
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc:              funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc:             funcDoGetHealthcheckErr,
				DoGetMongoDBFunc:                 funcDoGetMongoDBOk,
				DoGetAuthorisationMiddlewareFunc: funcDoGetAuthorisationMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.MongoDB, ShouldBeTrue)
				So(svcList.HealthCheck, ShouldBeFalse)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {
			// setup (run before each `Convey` at this scope / indentation):
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc:              funcDoGetHTTPServer,
				DoGetHealthCheckFunc:             funcDoGetHealthcheckOk,
				DoGetMongoDBFunc:                 funcDoGetMongoDBOk,
				DoGetAuthorisationMiddlewareFunc: funcDoGetAuthorisationMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run succeeds and all the flags are set", func() {
				So(err, ShouldBeNil)
				So(svcList.MongoDB, ShouldBeTrue)
				So(svcList.HealthCheck, ShouldBeTrue)
			})

			Convey("The checkers are registered and the healthcheck and http server started", func() {
				So(len(hcMock.AddCheckCalls()), ShouldEqual, 3)
				So(hcMock.AddCheckCalls()[0].Name, ShouldResemble, "Mongo DB")
				So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
				So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, "localhost:25400")
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				So(initMock.DoGetAuthorisationMiddlewareCalls(), ShouldHaveLength, 1)
				//!!! a call needed to stop the server, maybe ?
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		// ADD CODE: put this code in, if you have Checkers to register
		Convey("Given that Checkers cannot be registered", func() {
			// setup (run before each `Convey` at this scope / indentation):
			errAddheckFail := errors.New("Error(s) registering checkers for health check")
			hcMockAddFail := &mock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddheckFail },
				StartFunc:    func(ctx context.Context) {},
			}

			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMockAddFail, nil
				},
				DoGetMongoDBFunc:                 funcDoGetMongoDBOk,
				DoGetAuthorisationMiddlewareFunc: funcDoGetAuthorisationMiddleware,
				// ADD CODE: add the checkers that you want to register here
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails, but all checks try to register", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, fmt.Sprintf("unable to register checkers: %s", errAddheckFail.Error()))
				So(svcList.HealthCheck, ShouldBeTrue)
				So(len(hcMockAddFail.AddCheckCalls()), ShouldEqual, 3)
				So(hcMockAddFail.AddCheckCalls()[0].Name, ShouldResemble, "Mongo DB")
			})
			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised but the http server fails", func() {
			// setup (run before each `Convey` at this scope / indentation):
			initMock := &mock.InitialiserMock{
				DoGetHealthCheckFunc:             funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:              funcDoGetFailingHTTPServer,
				DoGetMongoDBFunc:                 funcDoGetMongoDBOk,
				DoGetAuthorisationMiddlewareFunc: funcDoGetAuthorisationMiddleware,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			Convey("Then the error is returned in the error channel", func() {
				sErr := <-svcErrors
				So(sErr.Error(), ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
				So(len(failingServerMock.ListenAndServeCalls()), ShouldEqual, 1)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that initialisation of the authorisation middleware fails", func() {
			expectedError := errors.New("failed to init authorisation middleware")
			initMock := &mock.InitialiserMock{
				DoGetHealthCheckFunc: funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:  funcDoGetFailingHTTPServer,
				DoGetMongoDBFunc:     funcDoGetMongoDBOk,
				DoGetAuthorisationMiddlewareFunc: func(ctx context.Context, authorisationConfig *authorisation.Config) (authorisation.Middleware, error) {
					return nil, expectedError
				},
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the expected error", func() {
				So(err, ShouldEqual, expectedError)
			})
		})
	})
}

func TestClose(t *testing.T) {
	Convey("Having a correctly initialised service", t, func() {
		cfg, err := config.Get()
		serverStopped := false
		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &mock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				serverStopped = true
				return nil
			},
		}

		funcDoGetAuthorisationMiddleware := func(ctx context.Context, authorisationConfig *authorisation.Config) (authorisation.Middleware, error) {
			return &authorisationMock.MiddlewareMock{
				RequireFunc: func(permission string, handlerFunc http.HandlerFunc) http.HandlerFunc {
					return handlerFunc
				},
				CloseFunc: func(ctx context.Context) error {
					return nil
				},
			}, nil
		}

		// mongoDB Close will fail if healthcheck and http server are not already closed
		mongoDBMock := &mock.PermissionsStoreMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			CloseFunc: func(ctx context.Context) error {
				if !hcStopped || !serverStopped {
					return errors.New("MongoDB closed before stopping healthcheck or HTTP server")
				}
				return nil
			},
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {
			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetMongoDBFunc: func(ctx context.Context, cfg *config.Config) (service.PermissionsStore, error) {
					return mongoDBMock, nil
				},
				DoGetAuthorisationMiddlewareFunc: funcDoGetAuthorisationMiddleware,
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(serverMock.ShutdownCalls()), ShouldEqual, 1)
			So(len(mongoDBMock.CloseCalls()), ShouldEqual, 1)
		})

		Convey("If services fail to stop, the Close operation tries to close all dependencies and returns an error", func() {
			failingserverMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					return errors.New("Failed to stop http server")
				},
			}

			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return failingserverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetMongoDBFunc: func(ctx context.Context, cfg *config.Config) (service.PermissionsStore, error) {
					return mongoDBMock, nil
				},
				DoGetAuthorisationMiddlewareFunc: funcDoGetAuthorisationMiddleware,
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(failingserverMock.ShutdownCalls()), ShouldEqual, 1)
			So(len(mongoDBMock.CloseCalls()), ShouldEqual, 1)
		})

		Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			timeoutServerMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					time.Sleep(3 * time.Millisecond)
					return nil
				},
			}

			svcList := service.NewServiceList(nil)
			svcList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				ServiceList: svcList,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "context deadline exceeded")
			So(len(hcMock.StopCalls()), ShouldEqual, 1)
			So(len(timeoutServerMock.ShutdownCalls()), ShouldEqual, 1)
		})
	})
}
