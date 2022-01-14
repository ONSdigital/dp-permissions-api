package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ONSdigital/dp-permissions-api/config"

	dpMongodb "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	log.Info(ctx, "loaded config", log.Data{"config": cfg})

	mongoConnection, err := dpMongodb.Open(&cfg.MongoDB)
	if err != nil {
		log.Error(ctx, "error initialising mongo", err)
		os.Exit(1)
	}

	importRoles(ctx, mongoConnection.Collection(cfg.Collections[config.RolesCollection]))
	importPolicies(ctx, mongoConnection.Collection(cfg.Collections[config.PoliciesCollection]))
}

func importRoles(ctx context.Context, mongoCollection *dpMongodb.Collection) {
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

		_, err = mongoCollection.UpsertById(ctx, role.ID, bson.M{"$set": role})
		if err != nil {
			log.Error(ctx, "failed to upsert role document", err, logData)
			os.Exit(1)
		}

		log.Info(ctx, "successfully put role into mongo", logData)
	}
}

func importPolicies(ctx context.Context, mongoCollection *dpMongodb.Collection) {
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
		logData := log.Data{"role": policy}

		_, err = mongoCollection.UpsertById(ctx, policy.ID, bson.M{"$set": policy})
		if err != nil {
			log.Error(ctx, "failed to upsert policy document", err, logData)
			os.Exit(1)
		}

		log.Info(ctx, "successfully put policy into mongo")
	}
}
