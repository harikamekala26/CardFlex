package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"cardflex-backend/middleware"
	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestJWTAuthAllowsProtectedRouteWithValidToken(t *testing.T) {
	db, tenant := setupTenantDB(t)

	token, err := utils.GenerateToken(
		strconv.FormatUint(uint64(1), 10),
		strconv.FormatUint(uint64(tenant.ID), 10),
		"test-secret",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authorized"})
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard?company=acme", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, res.Code, res.Body.String())
	}
}

func TestJWTAuthRejectsExpiredToken(t *testing.T) {
	db, tenant := setupTenantDB(t)

	expiredClaims := utils.Claims{
		UserID:   "1",
		TenantID: strconv.FormatUint(uint64(tenant.ID), 10),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authorized"})
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard?company=acme", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusUnauthorized, res.Code, res.Body.String())
	}
}

func TestJWTAuthRejectsTokenForDifferentTenant(t *testing.T) {
	db, _ := setupTenantDB(t)

	otherTenantID := uint(999)
	token, err := utils.GenerateToken(
		"1",
		strconv.FormatUint(uint64(otherTenantID), 10),
		"test-secret",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "authorized"})
	})

	req := httptest.NewRequest(http.MethodGet, "/dashboard?company=acme", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusForbidden, res.Code, res.Body.String())
	}
}

func setupTenantDB(t *testing.T) (*gorm.DB, models.Tenant) {
	t.Helper()

	dbName := strings.ReplaceAll(strings.ToLower(t.Name()), "/", "_")
	db, err := gorm.Open(sqlite.Open("file:"+dbName+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	tenant := models.Tenant{
		Name:        "Acme Card",
		CompanyCode: "acme",
		ThemeColor:  "#0B6E4F",
	}
	if err := db.Create(&tenant).Error; err != nil {
		t.Fatalf("failed to seed tenant: %v", err)
	}

	return db, tenant
}
