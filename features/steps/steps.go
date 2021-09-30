package steps

import (
	"context"
	"encoding/json"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/mongodb"
	"github.com/ONSdigital/dp-permissions-api/models"
	"time"

	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *PermissionsComponent) iHaveTheseRoles(rolesWriteJson *godog.DocString) error {
	ctx := context.Background()
	roles := []models.Role{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(rolesWriteJson.Content), &roles)
	if err != nil {
		return err
	}

	for _, roleDoc := range roles {
		if err := f.putRolesInDatabase(ctx, m.Connection, roleDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *PermissionsComponent) putRolesInDatabase(ctx context.Context, mongoConnection *dpMongoDriver.MongoConnection, roleDoc models.Role) error {
	update := bson.M{
		"$set": roleDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoConnection.GetConfiguredCollection().UpsertId(ctx, roleDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *PermissionsComponent) iHaveThesePolicies(jsonInput *godog.DocString) error {
	ctx := context.Background()
	policies := []models.Policy{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(jsonInput.Content), &policies)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err := f.putPolicyInDatabase(ctx, m.Connection, policy, f.Config.MongoConfig.PoliciesCollection); err != nil {
			return err
		}
	}

	return nil
}

func (f *PermissionsComponent) putPolicyInDatabase(
	ctx context.Context,
	mongoConnection *dpMongoDriver.MongoConnection,
	policy models.Policy,
	collection string) error {

	update := bson.M{
		"$set": policy,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoConnection.C(collection).UpsertId(ctx, policy.ID, update)
	if err != nil {
		return err
	}
	return nil
}
