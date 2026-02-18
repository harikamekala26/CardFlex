package routes

import (
	"cardflex-backend/controllers"
	"cardflex-backend/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(
	r *gin.Engine,
	tenantCollection *mongo.Collection,
	authController *controllers.AuthController,
	dashboardController *controllers.DashboardController,
	jwtSecret string,
) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	tenantAware := r.Group("/")
	tenantAware.Use(middleware.TenantResolver(tenantCollection))
	{
		tenantAware.POST("/register", authController.Register)
		tenantAware.POST("/login", authController.Login)
	}

	protected := r.Group("/")
	protected.Use(
		middleware.TenantResolver(tenantCollection),
		middleware.JWTAuth(jwtSecret),
	)
	{
		protected.GET("/dashboard", dashboardController.GetDashboard)
	}
}
