package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v2/pkg/health"
	dpMongodb "github.com/ONSdigital/dp-mongodb/v2/pkg/mongodb"
	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/log"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	connectTimeoutInSeconds = 5
	queryTimeoutInSeconds   = 15
)

//Mongo represents a simplistic MongoDB configuration, with session and health client
type Mongo struct {
	URI          string
	Database     string
	Collection   string
	Connection   *dpMongodb.MongoConnection
	Username     string
	Password     string
	healthClient *dpMongoHealth.CheckMongoClient
	IsSSL        bool
}

func (m *Mongo) getConnectionConfig(shouldEnableReadConcern, shouldEnableWriteConcern bool) *dpMongodb.MongoConnectionConfig {
	return &dpMongodb.MongoConnectionConfig{
		IsSSL:                   m.IsSSL,
		ConnectTimeoutInSeconds: connectTimeoutInSeconds,
		QueryTimeoutInSeconds:   queryTimeoutInSeconds,

		Username:                      m.Username,
		Password:                      m.Password,
		ClusterEndpoint:               m.URI,
		Database:                      m.Database,
		Collection:                    m.Collection,
		IsWriteConcernMajorityEnabled: shouldEnableWriteConcern,
		IsStrongReadConcernEnabled:    shouldEnableReadConcern,
	}
}

//Init creates a new mongoConnection with a strong consistency and a write mode of "majority"
func (m *Mongo) Init(ctx context.Context, shouldEnableReadConcern, shouldEnableWriteConcern bool) (err error) {
	if m.Connection != nil {
		return errors.New("datastore connection already exists")
	}

	mongoConnection, err := dpMongodb.Open(m.getConnectionConfig(shouldEnableReadConcern, shouldEnableWriteConcern))
	if err != nil {
		return err
	}

	m.Connection = mongoConnection
	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.Collection)}

	client := dpMongoHealth.NewClientWithCollections(mongoConnection, databaseCollectionBuilder)

	m.healthClient = &dpMongoHealth.CheckMongoClient{
		Client:      *client,
		Healthcheck: client.Healthcheck,
	}

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	if m.Connection == nil {
		return errors.New("cannot close a empty connection")
	}
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

//GetRole retrieves a role document by its ID
func (m *Mongo) GetRole(ctx context.Context, id string) (*models.Role, error) {
	log.Event(ctx, "getting role by ID", log.INFO, log.Data{"id": id})

	var role models.Role
	err := m.Connection.GetConfiguredCollection().FindOne(ctx, bson.M{"_id": id}, &role)
	if err != nil {
		if dpMongodb.IsErrNoDocumentFound(err) {
			return nil, apierrors.ErrRoleNotFound
		}
		return nil, err
	}

	return &role, nil
}

// GetRoles retrieves all role documents from Mongo, according to the provided limit and offset.
//Offset and limit need to  be positive or zero. Zero limit is equivalent to no limit and all items starting at the offset will be returned.
func (m *Mongo) GetRoles(ctx context.Context, offset, limit int) (*models.Roles, error) {
	if offset < 0 || limit < 0 {
		return nil, apierrors.ErrLimitAndOffset
	}

	log.Event(ctx, "querying document store for list of roles", log.INFO)

	roles := m.Connection.GetConfiguredCollection().Find(nil)
	totalCount, err := roles.Count(ctx)
	if err != nil {
		if dpMongodb.IsErrNoDocumentFound(err) {
			return nil, apierrors.ErrRoleNotFound
		}
		return nil, err
	}

	results := []models.Role{}
	iter := roles.Skip(offset).Limit(limit).Iter()
	if err := iter.All(ctx, &results); err != nil {
		return nil, err
	}

	return &models.Roles{
		Items:      results,
		Count:      len(results),
		TotalCount: totalCount,
		Offset:     offset,
		Limit:      limit,
	}, nil

}
