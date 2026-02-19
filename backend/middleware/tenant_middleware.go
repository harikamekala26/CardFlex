package middleware

import (
	"net/http"

	"cardflex-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TenantResolver(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		company := c.Query("company")
		if company == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "company query parameter is required"})
			c.Abort()
			return
		}

		var tenant models.Tenant
		err := db.Where("company_code = ?", company).First(&tenant).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
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
