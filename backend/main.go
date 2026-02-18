package main

import (
	"context"
	"log"
	"time"

	"cardflex-backend/config"
	"cardflex-backend/controllers"
	"cardflex-backend/migrations"
	"cardflex-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	env := config.LoadEnv()
	db := config.ConnectMongo(env)

	database := db.Client.Database(env.MongoDB)
	tenantCollection := database.Collection("tenants")
	userCollection := database.Collection("users")

	if err := runMigrations(tenantCollection, userCollection); err != nil {
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

func runMigrations(tenantCollection, userCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrations.MigrateTenants(ctx, tenantCollection); err != nil {
		return err
	}

	return migrations.MigrateUsers(ctx, userCollection)
}
