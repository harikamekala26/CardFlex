package main

import (
	"context"
	"log"
	"time"

	"cardflex-backend/config"
	"cardflex-backend/controllers"
	"cardflex-backend/models"
	"cardflex-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	env := config.LoadEnv()
	db := config.ConnectMongo(env)

	database := db.Client.Database(env.MongoDB)
	tenantCollection := database.Collection("tenants")
	userCollection := database.Collection("users")

	if err := ensureIndexesAndSeed(tenantCollection, userCollection); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authController := &controllers.AuthController{Users: userCollection, JWTSecret: env.JWTSecret}
	dashboardController := &controllers.DashboardController{}

	routes.RegisterRoutes(r, tenantCollection, authController, dashboardController, env.JWTSecret)

	log.Printf("CardFlex backend running on http://localhost:%s", env.Port)
	if err := r.Run(":" + env.Port); err != nil {
		log.Fatal(err)
	}
}

func ensureIndexesAndSeed(tenantCollection, userCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := tenantCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "companyCode", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_company_code"),
	})
	if err != nil {
		return err
	}

	_, err = userCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "tenantId", Value: 1}, {Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_tenant_email"),
	})
	if err != nil {
		return err
	}

	sampleTenants := []models.Tenant{
		{Name: "Acme Card", CompanyCode: "acme", ThemeColor: "#0B6E4F"},
		{Name: "Nova Finance", CompanyCode: "nova", ThemeColor: "#C84B31"},
		{Name: "Prime Credit", CompanyCode: "prime", ThemeColor: "#00539C"},
	}

	for _, t := range sampleTenants {
		_, err = tenantCollection.UpdateOne(
			ctx,
			bson.M{"companyCode": t.CompanyCode},
			bson.M{"$set": bson.M{"name": t.Name, "themeColor": t.ThemeColor, "companyCode": t.CompanyCode}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
