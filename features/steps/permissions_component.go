package steps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-component-test/utils"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/mongodb"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/mongo"
	"github.com/ONSdigital/dp-permissions-api/service"
	serviceMock "github.com/ONSdigital/dp-permissions-api/service/mock"
	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// PermissionsComponent holds the initialized http server, mongo client and configs required for running component tests.
type PermissionsComponent struct {
	ErrorFeature   componenttest.ErrorFeature
	svc            *service.Service
	errorChan      chan error
	MongoClient    *mongo.Mongo
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
}

// NewPermissionsComponent initializes mock server and inmemory mongodb used for running component tests.
func NewPermissionsComponent(mongoFeature *componenttest.MongoFeature) (*PermissionsComponent, error) {

	f := &PermissionsComponent{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error

	f.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	getMongoURI := fmt.Sprintf("localhost:%d", mongoFeature.Server.Port())
	databaseName := utils.RandomDatabase()

	f.Config.MongoConfig.Database = databaseName
	f.Config.MongoConfig.BindAddr = getMongoURI
	f.Config.MongoConfig.Username, f.Config.MongoConfig.Password, err = createCredsInDB(getMongoURI, databaseName)
	if err != nil {
		return nil, err
	}

	mongodb := &mongo.Mongo{}

	if err := mongodb.Init(f.Config.MongoConfig); err != nil {
		return nil, err
	}

	f.MongoClient = mongodb

	return f, nil
}

func createCredsInDB(getMongoURI string, databaseName string) (string, string, error) {
	username := "admin"
	password, _ := uuid.NewV4()
	mongoConnectionConfig := &dpMongoDriver.MongoConnectionConfig{
		IsSSL:                   false,
		ConnectTimeoutInSeconds: 15,
		QueryTimeoutInSeconds:   15,

		Username:        "",
		Password:        "",
		ClusterEndpoint: getMongoURI,
		Database:        databaseName,
	}
	mongoConnection, err := dpMongoDriver.Open(mongoConnectionConfig)
	if err != nil {
		return username, password.String(), errors.New(fmt.Sprintf("expected db connection to be opened: %+v", err))
	}
	mongoDatabaseSelection := mongoConnection.
		GetMongoCollection().
		Database()
	createCollectionResponse := mongoDatabaseSelection.RunCommand(context.TODO(), bson.D{
		{"create", "test"},
	})
	if createCollectionResponse.Err() != nil {
		return username, password.String(), errors.New(fmt.Sprintf("expected database creation to go through: %+v", err))
	}
	userCreationResponse := mongoDatabaseSelection.RunCommand(context.TODO(), bson.D{
		{"createUser", username},
		{"pwd", password.String()},
		{"roles", []bson.M{
			{"role": "root", "db": "admin"},
		}},
	})
	if userCreationResponse.Err() != nil {
		return username, password.String(), errors.New(fmt.Sprintf("expected user creation to go through: %+v", err))
	}
	return username, password.String(), nil
}

func (f *PermissionsComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have this roles:$`, f.iHaveTheseRoles)
	ctx.Step(`^I have these policies:$`, f.iHaveThesePolicies)
}

func (f *PermissionsComponent) Reset() *PermissionsComponent {
	f.MongoClient.Database = utils.RandomDatabase()
	f.MongoClient.Init(f.Config.MongoConfig)
	return f
}

func (f *PermissionsComponent) Close() error {
	if f.svc != nil && f.ServiceRunning {
		f.svc.Close(context.Background())
		f.ServiceRunning = false
	}
	return nil
}

func (f *PermissionsComponent) InitialiseService() (http.Handler, error) {
	initMock := &serviceMock.InitialiserMock{
		DoGetMongoDBFunc:     f.DoGetMongoDB,
		DoGetHealthCheckFunc: f.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:  f.DoGetHTTPServer,
	}

	if service, err := service.Run(context.Background(), f.Config, service.NewServiceList(initMock), "1", "", "", f.errorChan); err != nil {
		return nil, err
	} else {
		f.svc = service
	}
	f.ServiceRunning = true
	return f.HTTPServer.Handler, nil
}

func (f *PermissionsComponent) DoGetHealthcheckOk(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &serviceMock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (f *PermissionsComponent) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	f.HTTPServer.Addr = bindAddr
	f.HTTPServer.Handler = router
	return f.HTTPServer
}

// DoGetMongoDB returns a MongoDB
func (f *PermissionsComponent) DoGetMongoDB(ctx context.Context, cfg *config.Config) (service.PermissionsStore, error) {
	return f.MongoClient, nil
}
