package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"cardflex-backend/controllers"
	"cardflex-backend/middleware"
	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestGetProfileReturnsAuthenticatedUser(t *testing.T) {
	db, tenant, user, _ := setupControllerTenantDB(t)
	res := performProfileRequest(t, db, tenant, user.ID, "acme", "test-secret", true)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, res.Code, res.Body.String())
	}

	var response struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse profile response: %v", err)
	}

	if response.Name != user.Name {
		t.Fatalf("expected name %q, got %q", user.Name, response.Name)
	}
	if response.Email != user.Email {
		t.Fatalf("expected email %q, got %q", user.Email, response.Email)
	}
}

func TestGetProfileRequiresCompanyQuery(t *testing.T) {
	db, tenant, user, _ := setupControllerTenantDB(t)
	res := performProfileRequest(t, db, tenant, user.ID, "", "test-secret", true)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, res.Code, res.Body.String())
	}
}

func TestGetProfileRequiresValidJWT(t *testing.T) {
	db, _, _, _ := setupControllerTenantDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	profileController := &controllers.ProfileController{DB: db}

	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.GET("/profile", profileController.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/profile?company=acme", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusUnauthorized, res.Code, res.Body.String())
	}
}

func TestGetProfileRejectsInvalidJWT(t *testing.T) {
	db, _, _, _ := setupControllerTenantDB(t)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	profileController := &controllers.ProfileController{DB: db}

	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.GET("/profile", profileController.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/profile?company=acme", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-token")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusUnauthorized, res.Code, res.Body.String())
	}
}

func TestGetProfileRejectsTenantMismatch(t *testing.T) {
	db, tenant, user, _ := setupControllerTenantDB(t)
	otherTenant := models.Tenant{
		Name:        "Globex Card",
		CompanyCode: "globex",
		ThemeColor:  "#123456",
	}
	if err := db.Create(&otherTenant).Error; err != nil {
		t.Fatalf("failed to seed other tenant: %v", err)
	}

	res := performProfileRequest(t, db, tenant, user.ID, "globex", "test-secret", true)

	if res.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusForbidden, res.Code, res.Body.String())
	}
}

func TestGetProfileReturnsNotFoundWhenTenantMissing(t *testing.T) {
	db, tenant, user, _ := setupControllerTenantDB(t)
	res := performProfileRequest(t, db, tenant, user.ID, "missing", "test-secret", true)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusNotFound, res.Code, res.Body.String())
	}
}

func TestGetProfileReturnsNotFoundWhenUserMissing(t *testing.T) {
	db, tenant, _, _ := setupControllerTenantDB(t)
	res := performProfileRequest(t, db, tenant, 9999, "acme", "test-secret", true)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusNotFound, res.Code, res.Body.String())
	}
	if res.Body.String() != "{\"error\":\"user not found\"}" {
		t.Fatalf("expected user not found body, got %s", res.Body.String())
	}
}

func performProfileRequest(
	t *testing.T,
	db *gorm.DB,
	tokenTenant models.Tenant,
	userID uint,
	company string,
	secret string,
	includeToken bool,
) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	profileController := &controllers.ProfileController{DB: db}

	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth(secret))
	protected.GET("/profile", profileController.GetProfile)

	target := "/profile"
	if company != "" {
		target += "?company=" + company
	}
	req := httptest.NewRequest(http.MethodGet, target, nil)
	if includeToken {
		token, err := utils.GenerateToken(
			strconv.FormatUint(uint64(userID), 10),
			strconv.FormatUint(uint64(tokenTenant.ID), 10),
			secret,
		)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	return res
}
