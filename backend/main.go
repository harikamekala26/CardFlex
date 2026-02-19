package main

import (
	"log"

	"cardflex-backend/config"
	"cardflex-backend/controllers"
	"cardflex-backend/migrations"
	"cardflex-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	env := config.LoadEnv()
	db := config.ConnectSQL(env)
	if err := runMigrations(db.Client); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	authController := &controllers.AuthController{DB: db.Client, JWTSecret: env.JWTSecret}
	dashboardController := &controllers.DashboardController{}

	routes.RegisterRoutes(r, db.Client, authController, dashboardController, env.JWTSecret)

	log.Printf("CardFlex backend running on http://localhost:%s", env.Port)
	if err := r.Run(":" + env.Port); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(db *gorm.DB) error {
	if err := migrations.MigrateTenants(db); err != nil {
		return err
	}

	return migrations.MigrateUsers(db)
}
