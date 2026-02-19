package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cardflex-backend/controllers"
	"cardflex-backend/middleware"
	"cardflex-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterInsertsUserForTenant(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}); err != nil {
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

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authController := &controllers.AuthController{
		DB:        db,
		JWTSecret: "test-secret",
	}

	tenantAware := r.Group("/")
	tenantAware.Use(middleware.TenantResolver(db))
	tenantAware.POST("/register", authController.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register?company=acme", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, res.Code, res.Body.String())
	}

	var registerResponse struct {
		Message string `json:"message"`
		UserID  uint   `json:"userId"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &registerResponse); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}

	if registerResponse.Message != "user registered successfully" {
		t.Fatalf("unexpected register message: %q", registerResponse.Message)
	}
	if registerResponse.UserID == 0 {
		t.Fatal("expected non-zero user id")
	}

	var inserted models.User
	err = db.Where("tenant_id = ? AND email = ?", tenant.ID, "alice@example.com").First(&inserted).Error
	if err != nil {
		t.Fatalf("expected inserted user, find failed: %v", err)
	}

	if inserted.Name != "Alice" {
		t.Fatalf("expected name Alice, got %q", inserted.Name)
	}
	if inserted.TenantID != tenant.ID {
		t.Fatalf("expected tenant id %d, got %d", tenant.ID, inserted.TenantID)
	}
	if inserted.Password == "secret123" {
		t.Fatal("expected password to be hashed, but plain text was stored")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(inserted.Password), []byte("secret123")); err != nil {
		t.Fatalf("stored password is not a valid hash for original password: %v", err)
	}
}
