package main

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"os"
	"strings"

	dpMongodb "github.com/ONSdigital/dp-mongodb/v2/mongodb"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	log.Info(ctx, "loaded config", log.Data{"config": cfg})

	mongoConnection, err := dpMongodb.Open(getConnectionConfig(cfg.MongoConfig))
	if err != nil {
		log.Error(ctx, "error initialising mongo", err)
		os.Exit(1)
	}

	importRoles(ctx, mongoConnection)
	importPolicies(ctx, mongoConnection)
}

func getConnectionConfig(mongoConf config.MongoDB) *dpMongodb.MongoConnectionConfig {
	return &dpMongodb.MongoConnectionConfig{
		IsSSL:                   mongoConf.IsSSL,
		ConnectTimeoutInSeconds: 5,
		QueryTimeoutInSeconds:   15,

		Username:                      mongoConf.Username,
		Password:                      mongoConf.Password,
		ClusterEndpoint:               mongoConf.BindAddr,
		Database:                      mongoConf.Database,
		Collection:                    mongoConf.RolesCollection,
		IsWriteConcernMajorityEnabled: mongoConf.EnableWriteConcern,
		IsStrongReadConcernEnabled:    mongoConf.EnableReadConcern,
	}
}

func importRoles(ctx context.Context, mongoConnection *dpMongodb.MongoConnection) {
	filename := "roles.json"
	fileLocation := "./" + filename
	f, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(ctx, "failed to open roles json file", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(ctx, "failed to read json file as a byte array", err)
		os.Exit(1)
	}

	res := []models.Role{}
	if err := json.Unmarshal(b, &res); err != nil {
		logData := log.Data{"json file": res}
		log.Error(ctx, "failed to unmarshal json", err, logData)
		os.Exit(1)
	}

	for _, role := range res {

		role.ID = strings.ToLower(role.Name)
		logData := log.Data{"role": role}

		_, err = mongoConnection.C("roles").InsertOne(ctx, role)
		if err != nil {
			log.Error(ctx, "failed to insert new role document, data lost in mongo but exists in this log", err, logData)
			os.Exit(1)
		}

		log.Info(ctx, "successfully put role into mongo", logData)
	}
}

func importPolicies(ctx context.Context, mongoConnection *dpMongodb.MongoConnection) {
	filename := "policies.json"
	fileLocation := "./" + filename
	f, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(ctx, "failed to open policies json file", err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(ctx, "failed to read policies json file as a byte array", err)
		os.Exit(1)
	}

	res := []models.Policy{}
	if err := json.Unmarshal(b, &res); err != nil {
		logData := log.Data{"json file": res}
		log.Error(ctx, "failed to unmarshal policies json", err, logData)
		os.Exit(1)
	}

	for _, policy := range res {

		uuid, err := uuid.NewV4()
		if err != nil {
			log.Error(ctx, "failed to create a new UUID for policy", err)
			os.Exit(1)
		}
		policy.ID = uuid.String()

		_, err = mongoConnection.C("policies").InsertOne(ctx, policy)
		if err != nil {
			log.Error(ctx, "failed to insert new policy document, data lost in mongo but exists in this log", err)
			os.Exit(1)
		}

		log.Info(ctx, "successfully put policy into mongo")
	}
}
