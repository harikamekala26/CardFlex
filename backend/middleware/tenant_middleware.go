package middleware

import (
	"context"
	"net/http"
	"time"

	"cardflex-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TenantResolver(tenantCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		company := c.Query("company")
		if company == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "company query parameter is required"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var tenant models.Tenant
		err := tenantCollection.FindOne(ctx, bson.M{"companyCode": company}).Decode(&tenant)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve tenant"})
			c.Abort()
			return
		}

		c.Set("tenant", tenant)
		c.Next()
	}
}
