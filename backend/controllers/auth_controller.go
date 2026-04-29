package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB        *gorm.DB
	JWTSecret string
}

type registerRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	TenantID    string `json:"tenantId"`
	CompanyCode string `json:"companyCode"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

var strictEmailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]{2,}$`)

func (a *AuthController) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := a.resolveRegisterTenant(c, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		case err.Error() == "tenant identifier is required":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve tenant"})
		}
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

	if err := a.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		account := defaultAccountForUser(user.ID, tenant.ID)
		return tx.Create(&account).Error
	}); err != nil {
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

func defaultAccountForUser(userID uint, tenantID uint) models.Account {
	return models.Account{
		UserID:           userID,
		TenantID:         tenantID,
		MaskedCardNumber: "**** **** **** 0000",
		CreditLimit:      5000,
		AvailableBalance: 5000,
		Currency:         "USD",
	}
}

func (a *AuthController) resolveRegisterTenant(c *gin.Context, req registerRequest) (models.Tenant, error) {
	if tenantRaw, ok := c.Get("tenant"); ok {
		tenant, ok := tenantRaw.(models.Tenant)
		if ok {
			return tenant, nil
		}
	}

	tenantIdentifier := strings.TrimSpace(req.TenantID)
	if tenantIdentifier == "" {
		tenantIdentifier = strings.TrimSpace(req.CompanyCode)
	}
	if tenantIdentifier == "" {
		tenantIdentifier = strings.TrimSpace(c.Query("company"))
	}
	if tenantIdentifier == "" {
		return models.Tenant{}, errors.New("tenant identifier is required")
	}

	var tenant models.Tenant
	err := a.DB.Where("company_code = ?", strings.ToLower(tenantIdentifier)).First(&tenant).Error
	if err != nil {
		return models.Tenant{}, err
	}

	return tenant, nil
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
	if err := a.ensureAccountForUser(user.ID, tenant.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize account"})
		return
	}

	token, err := utils.GenerateToken(
		strconv.FormatUint(uint64(user.ID), 10),
		strconv.FormatUint(uint64(user.TenantID), 10),
		a.JWTSecret,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user logged in",
		"token":   token,
	})
}

func (a *AuthController) ensureAccountForUser(userID uint, tenantID uint) error {
	account := defaultAccountForUser(userID, tenantID)
	return a.DB.
		Where(models.Account{UserID: userID, TenantID: tenantID}).
		Attrs(account).
		FirstOrCreate(&account).Error
}
