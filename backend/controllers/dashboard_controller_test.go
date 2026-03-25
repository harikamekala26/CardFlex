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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetDashboardReturnsTenantScopedData(t *testing.T) {
	db, tenant := setupControllerTenantDB(t)

	token, err := utils.GenerateToken(
		"42",
		strconv.FormatUint(uint64(tenant.ID), 10),
		"test-secret",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	dashboardController := &controllers.DashboardController{}

	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.GET("/dashboard", dashboardController.GetDashboard)

	req := httptest.NewRequest(http.MethodGet, "/dashboard?company=acme", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, res.Code, res.Body.String())
	}

	var response struct {
		Tenant struct {
			Name        string `json:"name"`
			CompanyCode string `json:"companyCode"`
			ThemeColor  string `json:"themeColor"`
		} `json:"tenant"`
		Card struct {
			MaskedCardNumber string  `json:"maskedCardNumber"`
			CreditLimit      float64 `json:"creditLimit"`
			AvailableBalance float64 `json:"availableBalance"`
			Currency         string  `json:"currency"`
		} `json:"card"`
		Transactions []struct {
			Date     string  `json:"date"`
			Merchant string  `json:"merchant"`
			Amount   float64 `json:"amount"`
			Status   string  `json:"status"`
		} `json:"transactions"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse dashboard response: %v", err)
	}

	if response.Tenant.CompanyCode != "acme" {
		t.Fatalf("expected company code acme, got %q", response.Tenant.CompanyCode)
	}
	if response.Tenant.Name != tenant.Name {
		t.Fatalf("expected tenant name %q, got %q", tenant.Name, response.Tenant.Name)
	}
	if response.Card.MaskedCardNumber == "" {
		t.Fatal("expected masked card number in dashboard response")
	}
	if len(response.Transactions) == 0 {
		t.Fatal("expected dashboard transactions")
	}
}

func setupControllerTenantDB(t *testing.T) (*gorm.DB, models.Tenant) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
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

	return db, tenant
}
