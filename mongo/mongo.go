package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v3/health"
	dpMongodb "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

type Mongo struct {
	dpMongodb.MongoConnectionConfig
	PoliciesCollection string

	Connection   *dpMongodb.MongoConnection
	healthClient *dpMongoHealth.CheckMongoClient
}

//Init creates a new mongoConnection with a strong consistency and a write mode of "majority"
func (m *Mongo) Init(mongoConf config.MongoDB) (err error) {
	m.MongoConnectionConfig = mongoConf.MongoConnectionConfig
	m.PoliciesCollection = mongoConf.PoliciesCollection
	m.Connection, err = dpMongodb.Open(&m.MongoConnectionConfig)
	if err != nil {
		return err
	}

	databaseCollectionBuilder := map[dpMongoHealth.Database][]dpMongoHealth.Collection{
		(dpMongoHealth.Database)(m.Database): {
			(dpMongoHealth.Collection)(m.Collection),
			(dpMongoHealth.Collection)(m.PoliciesCollection),
		},
	}
	m.healthClient = dpMongoHealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

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
	log.Info(ctx, "getting role by ID", log.Data{"id": id})

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

	log.Info(ctx, "querying document store for list of roles")

	roles := m.Connection.GetConfiguredCollection().Find(bson.D{})
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

// GetAllRoles returns all role documents, without pagination
func (m *Mongo) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	query := m.Connection.GetConfiguredCollection().Find(bson.D{})

	var roles []*models.Role
	if err := query.IterAll(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetAllBundlePolicies returns all policy documents for a permissions bundle, without pagination
func (m *Mongo) GetAllBundlePolicies(ctx context.Context) ([]*models.BundlePolicy, error) {

	query := m.Connection.C(m.PoliciesCollection).Find(bson.D{})

	var policies []*models.BundlePolicy
	if err := query.IterAll(ctx, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

//AddPolicy inserts new policy to data store
func (m *Mongo) AddPolicy(ctx context.Context, policy *models.Policy) (*models.Policy, error) {

	var documents []interface{}
	documents = append(documents, policy)
	if _, err := m.Connection.C(m.PoliciesCollection).InsertMany(ctx, documents); err != nil {
		return nil, err
	}
	return policy, nil
}

func (m *Mongo) GetPolicy(ctx context.Context, id string) (*models.Policy, error) {
	log.Info(ctx, "getting policy by id", log.Data{"id": id})

	var policy models.Policy
	err := m.Connection.C(m.PoliciesCollection).FindOne(ctx, bson.M{"_id": id}, &policy)
	if err != nil {
		if dpMongodb.IsErrNoDocumentFound(err) {
			return nil, apierrors.ErrPolicyNotFound
		}
		return nil, err
	}

	return &policy, nil

}

func (m *Mongo) UpdatePolicy(ctx context.Context, policy *models.Policy) (*models.UpdateResult, error) {
	log.Info(ctx, "update policy by id", log.Data{"id": policy.ID})
	updatePolicy := bson.M{
		"$set": policy,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}

	upsertResult, err := m.Connection.C(m.PoliciesCollection).UpsertById(ctx, policy.ID, updatePolicy)
	if err != nil {
		return nil, err
	}
	return &models.UpdateResult{ModifiedCount: upsertResult.ModifiedCount, UpsertedCount: upsertResult.UpsertedCount}, nil

}

func (m *Mongo) DeletePolicy(ctx context.Context, id string) error {
	log.Info(ctx, "deleting policy by id", log.Data{"id": id})

	var collectionDeleteResult *dpMongodb.CollectionDeleteResult

	collectionDeleteResult, err := m.Connection.C(m.PoliciesCollection).DeleteById(ctx, id)
	if err != nil {
		return err
	}

	if collectionDeleteResult.DeletedCount == 0 {
		return apierrors.ErrPolicyNotFound
	}

	return nil

}
