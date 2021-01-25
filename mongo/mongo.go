package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/health"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

//Mongo represents a simplistic MongoDB configuration, with session and health client
type Mongo struct {
	Session      *mgo.Session
	healthClient *dpMongoHealth.CheckMongoClient
	Database     string
	Collection   string
}

var (
	errRoleNotFound = errors.New("role not found")
)

//Init creates a new mgo.Session with a strong consistency and a write mode of "majority"
func (m *Mongo) Init(mongoConf config.MongoConfiguration) (err error) {
	if m.Session != nil {
		return errors.New("session already exists")
	}

	//Create Session
	if m.Session, err = mgo.Dial(mongoConf.BindAddr); err != nil {
		return err
	}
	m.Session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	m.Session.SetMode(mgo.Strong, true)

	m.Database = mongoConf.Database
	m.Collection = mongoConf.Collection

	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(mongoConf.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(mongoConf.Collection)}

	client := dpMongoHealth.NewClientWithCollections(m.Session, databaseCollectionBuilder)

	m.healthClient = &dpMongoHealth.CheckMongoClient{
		Client:      *client,
		Healthcheck: client.Healthcheck,
	}

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	if m.Session == nil {
		return errors.New("cannot close a mongoDB connection without a valid session")
	}
	return dpMongodb.Close(ctx, m.Session)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

//GetRole retrieves a role document by its ID
func (m *Mongo) GetRole(ctx context.Context, id string) (*models.Role, error) {
	s := m.Session.Copy()
	defer s.Close()
	log.Event(ctx, "getting role by ID", log.INFO, log.Data{"id": id})

	var role models.Role
	err := s.DB(m.Database).C(m.Collection).Find(bson.M{"_id": id}).One(&role)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errRoleNotFound
		}
		return nil, err
	}

	return &role, nil
}
