package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"cardflex-backend/controllers"
	"cardflex-backend/middleware"
	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetDashboardReturnsTenantScopedData(t *testing.T) {
	db, tenant, user, account := setupControllerTenantDB(t)

	transactions := []models.Transaction{
		{
			AccountID: account.ID,
			UserID:    user.ID,
			TenantID:  tenant.ID,
			Date:      time.Date(2026, time.February, 14, 0, 0, 0, 0, time.UTC),
			Merchant:  "Grocery Mart",
			Amount:    -82.41,
			Status:    "Posted",
		},
		{
			AccountID: account.ID,
			UserID:    user.ID,
			TenantID:  tenant.ID,
			Date:      time.Date(2026, time.February, 10, 0, 0, 0, 0, time.UTC),
			Merchant:  "Card Payment",
			Amount:    500.00,
			Status:    "Completed",
		},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("failed to seed transactions: %v", err)
	}

	token, err := utils.GenerateToken(
		strconv.FormatUint(uint64(user.ID), 10),
		strconv.FormatUint(uint64(tenant.ID), 10),
		"test-secret",
	)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	dashboardController := &controllers.DashboardController{DB: db}

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
	if response.Card.MaskedCardNumber != account.MaskedCardNumber {
		t.Fatalf("expected masked card number %q, got %q", account.MaskedCardNumber, response.Card.MaskedCardNumber)
	}
	if response.Card.CreditLimit != account.CreditLimit {
		t.Fatalf("expected credit limit %v, got %v", account.CreditLimit, response.Card.CreditLimit)
	}
	if response.Card.AvailableBalance != account.AvailableBalance {
		t.Fatalf("expected available balance %v, got %v", account.AvailableBalance, response.Card.AvailableBalance)
	}
	if response.Card.Currency != account.Currency {
		t.Fatalf("expected currency %q, got %q", account.Currency, response.Card.Currency)
	}
	if len(response.Transactions) != len(transactions) {
		t.Fatalf("expected %d dashboard transactions, got %d", len(transactions), len(response.Transactions))
	}
	if response.Transactions[0].Merchant != "Grocery Mart" {
		t.Fatalf("expected most recent transaction to be Grocery Mart, got %q", response.Transactions[0].Merchant)
	}
	if response.Transactions[0].Date != "2026-02-14" {
		t.Fatalf("expected formatted transaction date 2026-02-14, got %q", response.Transactions[0].Date)
	}
}

func setupControllerTenantDB(t *testing.T) (*gorm.DB, models.Tenant, models.User, models.Account) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Account{}, &models.Transaction{}); err != nil {
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

	user := models.User{
		Name:     "Test User",
		Email:    "test@acme.local",
		Password: "hashed-password",
		TenantID: tenant.ID,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	account := models.Account{
		UserID:           user.ID,
		TenantID:         tenant.ID,
		MaskedCardNumber: "**** **** **** 4821",
		CreditLimit:      12000,
		AvailableBalance: 8250,
		Currency:         "USD",
	}
	if err := db.Create(&account).Error; err != nil {
		t.Fatalf("failed to seed account: %v", err)
	}

	return db, tenant, user, account
}
