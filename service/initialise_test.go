package service_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/service"
	"github.com/ONSdigital/dp-permissions-api/service/mock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	. "github.com/smartystreets/goconvey/convey"
)

var errFunc = func() error {
	return errors.New("Server error")
}

var cfg, _ = config.Get()

func TestGetHTTPServer(t *testing.T) {
	Convey("Given a service list that includes a mocked server", t, func() {
		serverMock := &mock.HTTPServerMock{}
		newServiceMock := &mock.InitialiserMock{
			DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
				return serverMock
			},
		}
		r := mux.NewRouter()
		svcList := service.NewServiceList(newServiceMock)
		Convey(" When GetHTTPServer is called", func() {
			server := svcList.GetHTTPServer(cfg.BindAddr, r)
			Convey("Then then mock server is returned and had been initialised with the correct bind address", func() {
				So(newServiceMock.DoGetHTTPServerCalls(), ShouldHaveLength, 1)
				So(newServiceMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, cfg.BindAddr)
				So(server, ShouldEqual, serverMock)
			})
		})
	})

	Convey("Given a service list returns a mocked server that errors on ListenAndServe", t, func() {
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: errFunc,
		}
		newServiceMock := &mock.InitialiserMock{
			DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
				return serverMock
			},
		}
		svcErrors := make(chan error, 1)
		r := mux.NewRouter()
		var err error
		svcList := service.NewServiceList(newServiceMock)
		Convey("When the server is retrieved and started", func() {
			server := svcList.GetHTTPServer(cfg.BindAddr, r)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			go func() {
				if err := server.ListenAndServe(); err != nil {
					svcErrors <- err
				}
			}()

			select {
			case err = <-svcErrors:
				cancel()
			case <-ctx.Done():
				server.Shutdown(context.Background())
				t.Fatal("ListenAndServe returned no error")
			}
			Convey("Then the startup has failed and returns the expected error", func() {
				So(err.Error(), ShouldEqual, "Server error")
			})
		})
	})

	Convey("Given a service list that includes mocked server", t, func() {
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				return nil
			},
		}
		newServiceMock := &mock.InitialiserMock{
			DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
				return serverMock
			},
		}
		r := mux.NewRouter()
		svcList := service.NewServiceList(newServiceMock)
		svcErrors := make(chan error, 1)
		Convey(" When GetHTTPServer is called", func() {
			server := svcList.GetHTTPServer(cfg.BindAddr, r)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			go func() {
				if err := server.ListenAndServe(); err != nil {
					svcErrors <- err
					return
				}
				cancel()
			}()

			var err error
			select {
			case err = <-svcErrors:
				cancel()
			case errDone := <-ctx.Done():
				So(errDone, ShouldBeZeroValue)
			}
			Convey("Then the server returns nil", func() {
				So(newServiceMock.DoGetHTTPServerCalls(), ShouldHaveLength, 1)
				So(serverMock.ListenAndServeCalls(), ShouldHaveLength, 1)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestGetHealthCheck(t *testing.T) {

	Convey("Given a service list that returns a mocked healthchecker", t, func() {

		hcMock := &mock.HealthCheckerMock{}

		newServiceMock := &mock.InitialiserMock{
			DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
				return hcMock, nil
			},
		}
		svcList := service.NewServiceList(newServiceMock)
		Convey("When GetHealthCheck is called", func() {
			hc, err := svcList.GetHealthCheck(cfg, testBuildTime, testGitCommit, testVersion)
			Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
				So(svcList.HealthCheck, ShouldBeTrue)
				So(hc, ShouldEqual, hcMock)
				So(newServiceMock.DoGetHealthCheckCalls(), ShouldHaveLength, 1)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a service list that returns nil for healthcheck", t, func() {
		newServiceMock := &mock.InitialiserMock{
			DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
				return nil, errHealthcheck
			},
		}
		svcList := service.NewServiceList(newServiceMock)
		Convey("When GetHealthCheck is called", func() {
			hc, err := svcList.GetHealthCheck(cfg, testBuildTime, testGitCommit, testVersion)
			Convey("Then the HealthCheck flag is set to false and HealthCheck is nil", func() {
				So(hc, ShouldBeNil)
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.HealthCheck, ShouldBeFalse)
			})
		})
	})
}
func TestGetMongoDB(t *testing.T) {

	Convey("Given a service list that returns a mocked mongo permissions store", t, func() {

		mongoMock := &mock.PermissionsStoreMock{}

		newServiceMock := &mock.InitialiserMock{
			DoGetMongoDBFunc: func(ctx context.Context, cfg *config.Config) (service.PermissionsStore, error) {
				return mongoMock, nil
			},
		}
		svcList := service.NewServiceList(newServiceMock)
		Convey("When GetMongoDB is called", func() {
			m, err := svcList.GetMongoDB(ctx, cfg)
			Convey(" Then the mongo flag is set to true and the mongo permssions store is returned", func() {
				So(svcList.MongoDB, ShouldBeTrue)
				So(m, ShouldEqual, mongoMock)
				So(newServiceMock.DoGetMongoDBCalls(), ShouldHaveLength, 1)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given a service list that returns nil for mongo permissions store", t, func() {

		newServiceMock := &mock.InitialiserMock{
			DoGetMongoDBFunc: func(ctx context.Context, cfg *config.Config) (service.PermissionsStore, error) {
				return nil, errMongoDB
			},
		}
		svcList := service.NewServiceList(newServiceMock)
		Convey("When GetMongo is called", func() {
			mong, err := svcList.GetMongoDB(ctx, cfg)
			Convey("Then the mongo flag is set to false and the mongodb client is nil", func() {
				So(mong, ShouldBeNil)
				So(err, ShouldResemble, errMongoDB)
				So(svcList.MongoDB, ShouldBeFalse)
			})
		})
	})
}
