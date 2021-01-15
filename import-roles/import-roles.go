package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
)

var (
	mongoURL string
)

//Role represents a role that will be stored in mongo
type Role struct {
	ID          int      `bson:"id" json:"id"`
	Name        string   `bson:"name" json:"name"`
	Permissions []string `bson:"permissions" json:"permissions"`
}

func main() {
	flag.StringVar(&mongoURL, "mongo-url", mongoURL, "mongoDB URL")
	flag.Parse()

	ctx := context.Background()

	if mongoURL == "" {
		log.Event(ctx, "missing mongo-url flag", log.ERROR)
		return
	}

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Event(ctx, "unable to create mongo session", log.ERROR, log.Error(err))
		return
	}
	defer session.Close()

	filename := "roles.json"
	fileLocation := "./" + filename
	f, err := os.Open(fileLocation)
	if err != nil {
		log.Event(ctx, "failed to open roles json file", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Event(ctx, "failed to read json file as a byte array", log.ERROR, log.Error(err))
		return
	}

	res := []Role{}
	if err := json.Unmarshal(b, &res); err != nil {
		logData := log.Data{"json file": res}
		log.Event(ctx, "failed to unmarshal json", log.ERROR, log.Error(err), logData)
		return
	}

	for _, role := range res {

		roleToBeAdded := Role{
			role.ID,
			role.Name,
			role.Permissions,
		}

		logData := log.Data{"role": roleToBeAdded}

		if err = session.DB("permissions").C("permissions").Insert(roleToBeAdded); err != nil {
			log.Event(ctx, "failed to insert new edition document, data lost in mongo but exists in this log", log.ERROR, log.Error(err), logData)
			return
		}

		log.Event(ctx, "successfully put role into mongo", log.INFO, logData)
	}

}
