package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MigrateUsers(ctx context.Context, userCollection *mongo.Collection) error {
	_, err := userCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenantId", Value: 1}, {Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_tenant_email"),
	})

	return err
}
