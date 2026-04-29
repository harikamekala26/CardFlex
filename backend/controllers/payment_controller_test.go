package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

	// Verify database state
	var updatedAccount models.Account
	db.First(&updatedAccount, account.ID)
	assert.Equal(t, 1500.0, updatedAccount.AvailableBalance)

	var transaction models.Transaction
	db.First(&transaction)
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
}

func TestRecordPayment_ExceedsBalance(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	// Create test data
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
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
}

func TestRecordPayment_AccountNotFound(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	// Create test data
	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
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
}

func TestRecordPayment_MissingClaims(t *testing.T) {
	db := setupPaymentTestDB(t)
	controller := &PaymentController{DB: db}

	tenant := models.Tenant{
		ID:          1,
		CompanyCode: "TEST",
		Name:        "Test Company",
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
}
