package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"cardflex-backend/middleware"
	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPaymentTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Create tables
	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Account{}, &models.Transaction{}); err != nil {
		t.Fatalf("failed to migrate tables: %v", err)
	}

	return db
}

func TestRecordPayment_Success(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	// Create test data
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	user := models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
		TenantID: tenant.ID,
	}
	db.Create(&user)

	account := models.Account{
		ID:               1,
		UserID:           user.ID,
		TenantID:         tenant.ID,
		MaskedCardNumber: "****1234",
		CreditLimit:      5000,
		AvailableBalance: 2000,
		Currency:         "USD",
	}
	db.Create(&account)

	// Create request
	payload := paymentRequest{Amount: 500}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	// Set context values
	claims := &utils.Claims{UserID: "1"}
	c.Set("tenant", tenant)
	c.Set("claims", claims)

	// Call controller
	controller.RecordPayment(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response paymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "payment recorded successfully", response.Message)
	assert.Equal(t, 1500.0, response.UpdatedBalance)
	assert.Equal(t, 500.0, response.Amount)
	assert.NotZero(t, response.TransactionID)
	assert.NotEmpty(t, response.Timestamp)
	_, err = time.Parse("2006-01-02T15:04:05Z", response.Timestamp)
	assert.NoError(t, err)

	// Verify database state
	var updatedAccount models.Account
	db.First(&updatedAccount, account.ID)
	assert.Equal(t, 1500.0, updatedAccount.AvailableBalance)

	var transaction models.Transaction
	db.First(&transaction)
	assert.Equal(t, response.TransactionID, transaction.ID)
	assert.Equal(t, account.ID, transaction.AccountID)
	assert.Equal(t, user.ID, transaction.UserID)
	assert.Equal(t, tenant.ID, transaction.TenantID)
	assert.Equal(t, "Card Payment", transaction.Merchant)
	assert.Equal(t, -500.0, transaction.Amount)
	assert.Equal(t, "Posted", transaction.Status)
}

func TestRecordPayment_InvalidAmount(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	// Create test data
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	user := models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
		TenantID: tenant.ID,
	}
	db.Create(&user)

	account := models.Account{
		ID:               1,
		UserID:           user.ID,
		TenantID:         tenant.ID,
		MaskedCardNumber: "****1234",
		CreditLimit:      5000,
		AvailableBalance: 2000,
		Currency:         "USD",
	}
	db.Create(&account)

	// Create request with negative amount
	payload := paymentRequest{Amount: -100}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	claims := &utils.Claims{UserID: "1"}
	c.Set("tenant", tenant)
	c.Set("claims", claims)

	// Call controller
	controller.RecordPayment(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"amount must be a positive number"}`, w.Body.String())

	var updatedAccount models.Account
	db.First(&updatedAccount, account.ID)
	assert.Equal(t, 2000.0, updatedAccount.AvailableBalance)

	var transactionCount int64
	db.Model(&models.Transaction{}).Count(&transactionCount)
	assert.Equal(t, int64(0), transactionCount)
}

func TestRecordPayment_InvalidJSONBody(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBufferString(`{"amount":`))
	c.Request.Header.Set("Content-Type", "application/json")

	claims := &utils.Claims{UserID: "1"}
	c.Set("tenant", tenant)
	c.Set("claims", claims)

	controller.RecordPayment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"invalid request body"}`, w.Body.String())
}

func TestRecordPayment_ExceedsBalance(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	// Create test data
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	user := models.User{
		ID:       1,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashed_password",
		TenantID: tenant.ID,
	}
	db.Create(&user)

	account := models.Account{
		ID:               1,
		UserID:           user.ID,
		TenantID:         tenant.ID,
		MaskedCardNumber: "****1234",
		CreditLimit:      5000,
		AvailableBalance: 2000,
		Currency:         "USD",
	}
	db.Create(&account)

	// Create request with amount exceeding balance
	payload := paymentRequest{Amount: 3000}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	claims := &utils.Claims{UserID: "1"}
	c.Set("tenant", tenant)
	c.Set("claims", claims)

	// Call controller
	controller.RecordPayment(c)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error":"payment amount exceeds available balance"}`, w.Body.String())

	var updatedAccount models.Account
	db.First(&updatedAccount, account.ID)
	assert.Equal(t, 2000.0, updatedAccount.AvailableBalance)

	var transactionCount int64
	db.Model(&models.Transaction{}).Count(&transactionCount)
	assert.Equal(t, int64(0), transactionCount)
}

func TestRecordPayment_AccountNotFound(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	// Create test data
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	// Create request
	payload := paymentRequest{Amount: 500}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	claims := &utils.Claims{UserID: "999"} // Non-existent user
	c.Set("tenant", tenant)
	c.Set("claims", claims)

	// Call controller
	controller.RecordPayment(c)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"account not found"}`, w.Body.String())
}

func TestRecordPayment_MissingTenant(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	payload := paymentRequest{Amount: 500}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	claims := &utils.Claims{UserID: "1"}
	c.Set("claims", claims)
	// Not setting tenant

	// Call controller
	controller.RecordPayment(c)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"tenant context missing"}`, w.Body.String())
}

func TestRecordPayment_MissingClaims(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	payload := paymentRequest{Amount: 500}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	c.Set("tenant", tenant)
	// Not setting claims

	// Call controller
	controller.RecordPayment(c)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"authentication claims missing"}`, w.Body.String())
}

func TestRecordPayment_MissingJWT(t *testing.T) {
	db := setupPaymentTestDB(t)
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	controller := &PaymentController{DB: db}

	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.POST("/payment", controller.RecordPayment)

	body, _ := json.Marshal(paymentRequest{Amount: 500})
	req := httptest.NewRequest(http.MethodPost, "/payment?company=TEST", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.JSONEq(t, `{"error":"missing or invalid authorization header"}`, w.Body.String())
}

func TestRecordPayment_RejectsTenantMismatch(t *testing.T) {
	db := setupPaymentTestDB(t)
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&tenant)

	otherTenant := models.Tenant{
		ID:          2,
		CompanyCode: "OTHER",
		Name:        "Other Company",
		Features: models.FeatureFlags{
			"payments_enabled": true,
		},
	}
	db.Create(&otherTenant)

	token, err := utils.GenerateToken("1", strconv.FormatUint(uint64(otherTenant.ID), 10), "test-secret")
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	controller := &PaymentController{DB: db}

	protected := r.Group("/")
	protected.Use(middleware.TenantResolver(db), middleware.JWTAuth("test-secret"))
	protected.POST("/payment", controller.RecordPayment)

	body, _ := json.Marshal(paymentRequest{Amount: 500})
	req := httptest.NewRequest(http.MethodPost, "/payment?company=TEST", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.JSONEq(t, `{"error":"token tenant mismatch"}`, w.Body.String())
}

func TestRecordPayment_PaymentsDisabled(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
		Features: models.FeatureFlags{
			"payments_enabled": false,
		},
	}
	db.Create(&tenant)

	payload := paymentRequest{Amount: 500}
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/payment", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	claims := &utils.Claims{UserID: "1"}
	c.Set("tenant", tenant)
	c.Set("claims", claims)

	controller.RecordPayment(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.JSONEq(t, `{"error":"payments are disabled for this tenant"}`, w.Body.String())

	var transactionCount int64
	db.Model(&models.Transaction{}).Count(&transactionCount)
	assert.Equal(t, int64(0), transactionCount)
}
