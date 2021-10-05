package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gofrs/uuid"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/globalsign/mgo"
)

func main() {

	var (
		mongoURL string
	)

	flag.StringVar(&mongoURL, "mongo-url", "localhost:27017", "mongoDB URL")
	flag.Parse()

	ctx := context.Background()

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Error(ctx, "unable to create mongo session", err)
		os.Exit(1)
	}
	defer session.Close()

	importRoles(ctx, session)
	importPolicies(ctx, session)
}

func importRoles(ctx context.Context, session *mgo.Session) {
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

		if err = session.DB("permissions").C("roles").Insert(role); err != nil {
			log.Error(ctx, "failed to insert new edition document, data lost in mongo but exists in this log", err, logData)
			os.Exit(1)
		}

		log.Info(ctx, "successfully put role into mongo", logData)
	}
}

func importPolicies(ctx context.Context, session *mgo.Session) {
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

		if err = session.DB("permissions").C("policies").Insert(policy); err != nil {
			log.Error(ctx, "failed to insert new policy document, data lost in mongo but exists in this log", err)
			os.Exit(1)
		}

		log.Info(ctx, "successfully put policy into mongo")
	}
}
