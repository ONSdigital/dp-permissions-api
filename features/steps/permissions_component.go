package steps

import (
	"context"
	"net/http"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-component-test/utils"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/dp-permissions-api/mongo"
	"github.com/ONSdigital/dp-permissions-api/service"
	serviceMock "github.com/ONSdigital/dp-permissions-api/service/mock"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/ONSdigital/dp-authorisation/v2/authorisation"
	"github.com/ONSdigital/dp-authorisation/v2/authorisationtest"

	permsdk "github.com/ONSdigital/dp-permissions-api/sdk"

	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	publicSigningkey = map[string]string{
		"NeKb65194Jo=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyehkd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdgcKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbcmwIDAQAB",
	}
)

// PermissionsComponent holds the initialized http server, mongo client and configs required for running component tests.
type PermissionsComponent struct {
	ErrorFeature            componenttest.ErrorFeature
	svc                     *service.Service
	errorChan               chan error
	MongoClient             *mongo.Mongo
	Config                  *config.Config
	HTTPServer              *http.Server
	ServiceRunning          bool
	ApiFeature              *componenttest.APIFeature
	AuthorisationMiddleware authorisation.Middleware
}

func setupFakePermissionsAPI() *authorisationtest.FakePermissionsAPI {
	fakePermissionsAPI := authorisationtest.NewFakePermissionsAPI()
	bundle := getPermissionsBundle()
	fakePermissionsAPI.Reset()
	fakePermissionsAPI.UpdatePermissionsBundleResponse(bundle)
	return fakePermissionsAPI
}

// getPermissionsBundle seed's the PermissionsComponent bundle on startup
func getPermissionsBundle() *permsdk.Bundle {
	return &permsdk.Bundle{
		models.PoliciesCreate: { // role
			"groups/role-admin": { // groups
				permsdk.Policy{
					ID:        "policy1",
					Condition: permsdk.Condition{},
				},
			},
		},
		models.PoliciesRead: { // role
			"groups/role-admin": { // groups
				permsdk.Policy{
					ID:        "policy1",
					Condition: permsdk.Condition{},
				},
			},
			"groups/role-publisher": { // groups
				permsdk.Policy{
					ID:        "policy2",
					Condition: permsdk.Condition{},
				},
			},
			"groups/role-viewer": { // groups
				permsdk.Policy{
					ID:        "policy2",
					Condition: permsdk.Condition{},
				},
			},
		},
		models.PoliciesUpdate: { // role
			"groups/role-admin": { // groups
				permsdk.Policy{
					ID:        "policy3",
					Condition: permsdk.Condition{},
				},
			},
		},
		models.PoliciesDelete: { // role
			"groups/role-admin": { // groups
				permsdk.Policy{
					ID:        "policy1",
					Condition: permsdk.Condition{},
				},
			},
		},
		models.RolesRead: { // role
			"groups/role-admin": { // groups
				permsdk.Policy{
					ID:        "policy1",
					Condition: permsdk.Condition{},
				},
			},
		},
	}
}

// NewPermissionsComponent initializes mock server and in-memory mongodb used for running component tests.
func NewPermissionsComponent(mongoURI string) (*PermissionsComponent, error) {
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

	f.Config.MongoDB.ClusterEndpoint = mongoURI
	f.Config.MongoDB.Database = utils.RandomDatabase()
	// The following is to reset the Username and Password that have been set is Config from the previous
	// config.Get()
	f.Config.Username, f.Config.Password = "", ""
	f.Config.MongoDB.Username, f.Config.Password = createCredsInDB(f.Config.MongoDB)

	if f.MongoClient, err = mongo.NewMongoStore(context.Background(), f.Config.MongoDB); err != nil {
		return nil, err
	}

	f.ApiFeature = componenttest.NewAPIFeature(f.InitialiseService)

	fakePermissionsAPI := setupFakePermissionsAPI()
	f.Config.AuthorisationConfig.JWTVerificationPublicKeys = rsaJWKS
	f.Config.AuthorisationConfig.PermissionsAPIURL = fakePermissionsAPI.URL()

	return f, nil
}

func createCredsInDB(mongoConfig dpMongoDriver.MongoDriverConfig) (string, string) {
	mongoConnection, err := dpMongoDriver.Open(&mongoConfig)
	if err != nil {
		panic("expected db connection to be opened")
	}

	username := "admin"
	password, _ := uuid.NewV4()
	createCollectionResponse := mongoConnection.RunCommand(context.TODO(), bson.D{
		{Key: "create", Value: "test"},
	})
	if createCollectionResponse != nil {
		panic("expected test collection to be created")
	}
	userCreationResponse := mongoConnection.RunCommand(context.TODO(), bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: password.String()},
		{Key: "roles", Value: []bson.M{
			{"role": "root", "db": "admin"},
		}},
	})
	if userCreationResponse != nil {
		panic("expected admin user to be created")
	}

	return username, password.String()
}

func (f *PermissionsComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have these roles:$`, f.iHaveTheseRoles)
	ctx.Step(`^I have these policies:$`, f.iHaveThesePolicies)
	ctx.Step(`^I am an admin user$`, f.adminJWTToken)
	ctx.Step(`^I am a publisher user$`, f.publisherJWTToken)
	ctx.Step(`^I am a viewer user$`, f.viewerJWTToken)
	ctx.Step(`^I am a basic user$`, f.basicUserJWTToken)
	ctx.Step(`^I am a publisher user with invalid auth token$`, f.publisherWithNoJWTToken)
}

func (f *PermissionsComponent) Close() error {
	ctx := context.Background()
	err := f.MongoClient.Connection.DropDatabase(ctx)
	if err != nil {
		log.Warn(ctx, "error dropping database on Close()", log.Data{"err": err.Error()})
	}
	if f.svc != nil && f.ServiceRunning {
		err = f.svc.Close(ctx)
		if err != nil {
			log.Warn(ctx, "error closing service on Close()", log.Data{"err": err.Error()})
		}
		f.ServiceRunning = false
	}

	return nil
}

func (f *PermissionsComponent) InitialiseService() (http.Handler, error) {
	initMock := &serviceMock.InitialiserMock{
		DoGetMongoDBFunc:                 f.DoGetMongoDB,
		DoGetHealthCheckFunc:             f.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:              f.DoGetHTTPServer,
		DoGetAuthorisationMiddlewareFunc: f.DoGetAuthorisationMiddleware,
	}

	if svc, err := service.Run(context.Background(), f.Config, service.NewServiceList(initMock), "1", "", "", f.errorChan); err != nil {
		return nil, err
	} else {
		f.svc = svc
	}
	f.ServiceRunning = true

	return f.HTTPServer.Handler, nil
}

func (f *PermissionsComponent) DoGetHealthcheckOk(_ *config.Config, _ string, _ string, _ string) (service.HealthChecker, error) {
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
func (f *PermissionsComponent) DoGetMongoDB(_ context.Context, _ *config.Config) (service.PermissionsStore, error) {
	return f.MongoClient, nil
}

// DoGetAuthorisationMiddleware returns an authorisationMock.Middleware object
func (f *PermissionsComponent) DoGetAuthorisationMiddleware(ctx context.Context, cfg *authorisation.Config) (authorisation.Middleware, error) {
	middleware, err := authorisation.NewMiddlewareFromConfig(ctx, cfg, publicSigningkey)

	if err != nil {
		return nil, err
	}

	f.AuthorisationMiddleware = middleware
	return f.AuthorisationMiddleware, nil
}
