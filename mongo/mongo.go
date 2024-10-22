package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-permissions-api/apierrors"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/models"
	"github.com/ONSdigital/log.go/v2/log"

	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type Mongo struct {
	mongodriver.MongoDriverConfig

	Connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
}

// NewMongoStore creates a new Mongo object encapsulating a connection to the mongo server/cluster with the given configuration,
// and a health client to check the health of the mongo server/cluster
func NewMongoStore(_ context.Context, cfg config.MongoDB) (m *Mongo, err error) {
	m = &Mongo{MongoDriverConfig: cfg}

	m.Connection, err = mongodriver.Open(&m.MongoDriverConfig)
	if err != nil {
		return nil, err
	}

	databaseCollectionBuilder := map[mongohealth.Database][]mongohealth.Collection{
		mongohealth.Database(m.Database): {
			mongohealth.Collection(m.ActualCollectionName(config.RolesCollection)),
			mongohealth.Collection(m.ActualCollectionName(config.PoliciesCollection)),
		},
	}
	m.healthClient = mongohealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return m, nil
}

// Close the mongo session and returns any error
// It is an error to call m.Close if m.Init() returned an error, and there is no open connection
func (m *Mongo) Close(ctx context.Context) error {
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// GetRole retrieves a role document by its ID
func (m *Mongo) GetRole(ctx context.Context, id string) (*models.Role, error) {
	log.Info(ctx, "getting role by ID", log.Data{"id": id})

	var role models.Role
	err := m.Connection.Collection(m.ActualCollectionName(config.RolesCollection)).FindOne(ctx, bson.M{"_id": id}, &role)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocumentFound) {
			return nil, apierrors.ErrRoleNotFound
		}
		return nil, err
	}

	return &role, nil
}

// GetRoles retrieves all role documents from Mongo, according to the provided limit and offset.
// Offset and limit need to  be positive or zero.
func (m *Mongo) GetRoles(ctx context.Context, offset, limit int) (*models.Roles, error) {
	if offset < 0 || limit < 0 {
		return nil, apierrors.ErrLimitAndOffset
	}
	log.Info(ctx, "querying document store for list of roles")

	results := []models.Role{}
	totalCount, err := m.Connection.Collection(m.ActualCollectionName(config.RolesCollection)).Find(ctx, bson.D{}, &results,
		mongodriver.Sort(bson.M{"_id": 1}), mongodriver.Offset(offset), mongodriver.Limit(limit))
	if err != nil {
		return nil, err
	}
	if totalCount == 0 {
		return nil, apierrors.ErrRoleNotFound
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
	var roles []*models.Role
	if _, err := m.Connection.Collection(m.ActualCollectionName(config.RolesCollection)).Find(ctx, bson.D{}, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetAllBundlePolicies returns all policy documents for a permissions bundle, without pagination
func (m *Mongo) GetAllBundlePolicies(ctx context.Context) ([]*models.BundlePolicy, error) {
	var policies []*models.BundlePolicy
	if _, err := m.Connection.Collection(m.ActualCollectionName(config.PoliciesCollection)).Find(ctx, bson.D{}, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

// AddPolicy inserts new policy to data store
func (m *Mongo) AddPolicy(ctx context.Context, policy *models.Policy) (*models.Policy, error) {
	if _, err := m.Connection.Collection(m.ActualCollectionName(config.PoliciesCollection)).Insert(ctx, policy); err != nil {
		return nil, err
	}

	return policy, nil
}

// GetPolicy returns a policy given its id
func (m *Mongo) GetPolicy(ctx context.Context, id string) (*models.Policy, error) {
	log.Info(ctx, "getting policy by id", log.Data{"id": id})

	var policy models.Policy
	err := m.Connection.Collection(m.ActualCollectionName(config.PoliciesCollection)).FindOne(ctx, bson.M{"_id": id}, &policy)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocumentFound) {
			return nil, apierrors.ErrPolicyNotFound
		}
		return nil, err
	}

	return &policy, nil
}

// UpdatePolicy updates the given policy, or inserts/creates the given policy if it does not exist
func (m *Mongo) UpdatePolicy(ctx context.Context, policy *models.Policy) (*models.UpdateResult, error) {
	log.Info(ctx, "update policy by id", log.Data{"id": policy.ID})
	updatePolicy := bson.M{
		"$set": policy,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}

	upsertResult, err := m.Connection.Collection(m.ActualCollectionName(config.PoliciesCollection)).UpsertById(ctx, policy.ID, updatePolicy)
	if err != nil {
		return nil, err
	}

	return &models.UpdateResult{ModifiedCount: upsertResult.ModifiedCount, UpsertedCount: upsertResult.UpsertedCount}, nil
}

// DeletePolicy deletes a policy given its id
func (m *Mongo) DeletePolicy(ctx context.Context, id string) error {
	log.Info(ctx, "deleting policy by id", log.Data{"id": id})

	var collectionDeleteResult *mongodriver.CollectionDeleteResult

	collectionDeleteResult, err := m.Connection.Collection(m.ActualCollectionName(config.PoliciesCollection)).DeleteById(ctx, id)
	if err != nil {
		return err
	}

	if collectionDeleteResult.DeletedCount == 0 {
		return apierrors.ErrPolicyNotFound
	}

	return nil
}
