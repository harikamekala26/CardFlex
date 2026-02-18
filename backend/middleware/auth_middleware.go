package middleware

import (
	"net/http"
	"strings"

	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenString, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		tenantRaw, ok := c.Get("tenant")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context missing"})
			c.Abort()
			return
		}

		tenant := tenantRaw.(models.Tenant)
		if claims.TenantID != tenant.ID.Hex() {
			c.JSON(http.StatusForbidden, gin.H{"error": "token tenant mismatch"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
