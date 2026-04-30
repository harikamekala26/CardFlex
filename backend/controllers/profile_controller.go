package controllers

import (
	"net/http"
	"strconv"

	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileController struct {
	DB *gorm.DB
}

func (p *ProfileController) GetProfile(c *gin.Context) {
	if p.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "profile database is not configured"})
		return
	}

	tenantRaw, ok := c.Get("tenant")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context missing"})
		return
	}

	tenant, ok := tenantRaw.(models.Tenant)
	if !ok || tenant.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid tenant context"})
		return
	}
	if !tenant.Features.Enabled("profile_enabled") {
		c.JSON(http.StatusForbidden, gin.H{"error": "profiles are disabled for this tenant"})
		return
	}

	claimsRaw, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication claims missing"})
		return
	}

	claims, ok := claimsRaw.(*utils.Claims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication claims"})
		return
	}

	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user"})
		return
	}

	var user models.User
	if err := p.DB.
		Where("tenant_id = ? AND id = ?", tenant.ID, uint(userID)).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"email": user.Email,
	})
}
