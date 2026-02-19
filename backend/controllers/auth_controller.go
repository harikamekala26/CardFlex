package controllers

import (
	"net/http"
	"regexp"
	"strings"

	"cardflex-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB        *gorm.DB
	JWTSecret string
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

var strictEmailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]{2,}$`)

func (a *AuthController) Register(c *gin.Context) {
	tenantRaw, _ := c.Get("tenant")
	tenant := tenantRaw.(models.Tenant)

	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)
	if !strictEmailRegex.MatchString(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}
	var count int64
	if err := a.DB.Model(&models.User{}).
		Where("tenant_id = ? AND email = ?", tenant.ID, email).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate existing user"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists for this tenant"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		TenantID: tenant.ID,
	}

	if err := a.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user registered",
		"dummy": gin.H{
			"name":        name,
			"email":       email,
			"userId":      user.ID,
			"tenantId":    user.TenantID,
			"companyCode": tenant.CompanyCode,
		},
	})
}

func (a *AuthController) Login(c *gin.Context) {
	tenantRaw, _ := c.Get("tenant")
	tenant := tenantRaw.(models.Tenant)

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if !strictEmailRegex.MatchString(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}
	var user models.User
	if err := a.DB.Where("tenant_id = ? AND email = ?", tenant.ID, email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user logged in",
		"dummy": gin.H{
			"email":       email,
			"userId":      user.ID,
			"tenantId":    user.TenantID,
			"companyCode": tenant.CompanyCode,
		},
	})
}
