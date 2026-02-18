package migrations

import (
	"context"

	"cardflex-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MigrateTenants(ctx context.Context, tenantCollection *mongo.Collection) error {
	_, err := tenantCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "companyCode", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_company_code"),
	})
	if err != nil {
		return err
	}

	sampleTenants := []models.Tenant{
		{Name: "Acme Card", CompanyCode: "acme", ThemeColor: "#0B6E4F"},
		{Name: "Nova Finance", CompanyCode: "nova", ThemeColor: "#C84B31"},
		{Name: "Prime Credit", CompanyCode: "prime", ThemeColor: "#00539C"},
	}

	for _, tenant := range sampleTenants {
		_, err = tenantCollection.UpdateOne(
			ctx,
			bson.M{"companyCode": tenant.CompanyCode},
			bson.M{"$set": bson.M{
				"name":        tenant.Name,
				"companyCode": tenant.CompanyCode,
				"themeColor":  tenant.ThemeColor,
			}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
