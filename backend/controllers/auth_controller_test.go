package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"cardflex-backend/controllers"
	"cardflex-backend/middleware"
	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRegisterReturnsDummyResponseForTenant(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:register_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Account{}); err != nil {
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

	r.POST("/register", authController.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123","tenantId":"acme"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, res.Code, res.Body.String())
	}

	var registerResponse struct {
		Message string `json:"message"`
		Dummy   struct {
			Name        string `json:"name"`
			Email       string `json:"email"`
			UserID      uint   `json:"userId"`
			TenantID    uint   `json:"tenantId"`
			CompanyCode string `json:"companyCode"`
		} `json:"dummy"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &registerResponse); err != nil {
		t.Fatalf("failed to parse register response: %v", err)
	}

	if registerResponse.Message != "user registered" {
		t.Fatalf("unexpected register message: %q", registerResponse.Message)
	}
	if registerResponse.Dummy.Name != "Alice" {
		t.Fatalf("expected dummy name Alice, got %q", registerResponse.Dummy.Name)
	}
	if registerResponse.Dummy.Email != "alice@example.com" {
		t.Fatalf("expected normalized email alice@example.com, got %q", registerResponse.Dummy.Email)
	}
	if registerResponse.Dummy.UserID == 0 {
		t.Fatal("expected non-zero user id")
	}
	if registerResponse.Dummy.TenantID != tenant.ID {
		t.Fatalf("expected tenant id %d, got %d", tenant.ID, registerResponse.Dummy.TenantID)
	}
	if registerResponse.Dummy.CompanyCode != "acme" {
		t.Fatalf("expected company code acme, got %q", registerResponse.Dummy.CompanyCode)
	}

	var inserted models.User
	if err := db.Where("tenant_id = ? AND email = ?", tenant.ID, "alice@example.com").First(&inserted).Error; err != nil {
		t.Fatalf("expected inserted user, got error: %v", err)
	}
	if inserted.Password == "secret123" {
		t.Fatal("expected stored password to be hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(inserted.Password), []byte("secret123")); err != nil {
		t.Fatalf("stored hash does not match original password: %v", err)
	}

	var account models.Account
	if err := db.Where("tenant_id = ? AND user_id = ?", tenant.ID, inserted.ID).First(&account).Error; err != nil {
		t.Fatalf("expected registered user account, got error: %v", err)
	}
	if account.MaskedCardNumber != "**** **** **** 0000" {
		t.Fatalf("expected default masked card number, got %q", account.MaskedCardNumber)
	}
	if account.CreditLimit != 5000 || account.AvailableBalance != 5000 {
		t.Fatalf("expected default account balances to be 5000/5000, got %v/%v", account.CreditLimit, account.AvailableBalance)
	}
}

func TestRegisterSupportsCompanyQueryFallback(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:register_query_fallback_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Account{}); err != nil {
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
	r.POST("/register", authController.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register?company=acme", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, res.Code, res.Body.String())
	}
}

func TestRegisterRejectsDuplicateEmailForTenant(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:register_duplicate_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Account{}); err != nil {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if err := db.Create(&models.User{
		Name:     "Existing User",
		Email:    "alice@example.com",
		Password: string(hashedPassword),
		TenantID: tenant.ID,
	}).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authController := &controllers.AuthController{
		DB:        db,
		JWTSecret: "test-secret",
	}
	r.POST("/register", authController.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123","tenantId":"acme"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusConflict, res.Code, res.Body.String())
	}
}

func TestRegisterRequiresTenantIdentifier(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:register_missing_tenant_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Account{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authController := &controllers.AuthController{
		DB:        db,
		JWTSecret: "test-secret",
	}
	r.POST("/register", authController.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, res.Code, res.Body.String())
	}
}

func TestLoginReturnsDummyResponseForTenant(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:login_test?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Account{}); err != nil {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	user := models.User{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: string(hashedPassword),
		TenantID: tenant.ID,
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authController := &controllers.AuthController{
		DB:        db,
		JWTSecret: "test-secret",
	}

	tenantAware := r.Group("/")
	tenantAware.Use(middleware.TenantResolver(db))
	tenantAware.POST("/login", authController.Login)

	body := `{"email":"alice@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/login?company=acme", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, res.Code, res.Body.String())
	}

	var loginResponse struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &loginResponse); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}

	if loginResponse.Message != "user logged in" {
		t.Fatalf("unexpected login message: %q", loginResponse.Message)
	}
	if loginResponse.Token == "" {
		t.Fatal("expected JWT token in login response")
	}

	claims, err := utils.ParseToken(loginResponse.Token, "test-secret")
	if err != nil {
		t.Fatalf("expected parsable JWT token, got error: %v", err)
	}

	if claims.UserID != strconv.FormatUint(uint64(user.ID), 10) {
		t.Fatalf("expected user id claim %d, got %s", user.ID, claims.UserID)
	}
	if claims.TenantID != strconv.FormatUint(uint64(tenant.ID), 10) {
		t.Fatalf("expected tenant id claim %d, got %s", tenant.ID, claims.TenantID)
	}

	var account models.Account
	if err := db.Where("tenant_id = ? AND user_id = ?", tenant.ID, user.ID).First(&account).Error; err != nil {
		t.Fatalf("expected login to initialize missing account, got error: %v", err)
	}
	if account.AvailableBalance != 5000 {
		t.Fatalf("expected default account available balance 5000, got %v", account.AvailableBalance)
	}
}

func TestLoginValidationRejectsInvalidInputs(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:login_validation_test?mode=memory&cache=shared"), &gorm.Config{})
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if err := db.Create(&models.User{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: string(hashedPassword),
		TenantID: tenant.ID,
	}).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authController := &controllers.AuthController{
		DB:        db,
		JWTSecret: "test-secret",
	}

	tenantAware := r.Group("/")
	tenantAware.Use(middleware.TenantResolver(db))
	tenantAware.POST("/login", authController.Login)

	invalidEmailBody := `{"email":"alice@localhost","password":"secret123"}`
	invalidEmailReq := httptest.NewRequest(http.MethodPost, "/login?company=acme", strings.NewReader(invalidEmailBody))
	invalidEmailReq.Header.Set("Content-Type", "application/json")
	invalidEmailRes := httptest.NewRecorder()
	r.ServeHTTP(invalidEmailRes, invalidEmailReq)

	if invalidEmailRes.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for invalid email, got %d, body: %s", http.StatusBadRequest, invalidEmailRes.Code, invalidEmailRes.Body.String())
	}

	shortPasswordBody := `{"email":"alice@example.com","password":"123"}`
	shortPasswordReq := httptest.NewRequest(http.MethodPost, "/login?company=acme", strings.NewReader(shortPasswordBody))
	shortPasswordReq.Header.Set("Content-Type", "application/json")
	shortPasswordRes := httptest.NewRecorder()
	r.ServeHTTP(shortPasswordRes, shortPasswordReq)

	if shortPasswordRes.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for short password, got %d, body: %s", http.StatusBadRequest, shortPasswordRes.Code, shortPasswordRes.Body.String())
	}

	wrongPasswordBody := `{"email":"alice@example.com","password":"wrongpass"}`
	wrongPasswordReq := httptest.NewRequest(http.MethodPost, "/login?company=acme", strings.NewReader(wrongPasswordBody))
	wrongPasswordReq.Header.Set("Content-Type", "application/json")
	wrongPasswordRes := httptest.NewRecorder()
	r.ServeHTTP(wrongPasswordRes, wrongPasswordReq)

	if wrongPasswordRes.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d for wrong password, got %d, body: %s", http.StatusUnauthorized, wrongPasswordRes.Code, wrongPasswordRes.Body.String())
	}
}
