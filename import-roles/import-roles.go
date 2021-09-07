package main

import (
	"context"
	"encoding/json"
	"flag"
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
