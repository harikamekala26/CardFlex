package routes

import (
	"cardflex-backend/controllers"
	"cardflex-backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	r *gin.Engine,
	db *gorm.DB,
	authController *controllers.AuthController,
	dashboardController *controllers.DashboardController,
	jwtSecret string,
) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	tenantAware := r.Group("/")
	tenantAware.Use(middleware.TenantResolver(db))
	{
		tenantAware.POST("/register", authController.Register)
		tenantAware.POST("/login", authController.Login)
	}

	protected := r.Group("/")
	protected.Use(
		middleware.TenantResolver(db),
		middleware.JWTAuth(jwtSecret),
	)
	{
		protected.GET("/dashboard", dashboardController.GetDashboard)
	}
}
