package controllers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"cardflex-backend/controllers"
	"cardflex-backend/middleware"
	"cardflex-backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterInsertsUserForTenant(t *testing.T) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		t.Skip("MONGO_URI is not set; skipping integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer disconnectCancel()
		_ = client.Disconnect(disconnectCtx)
	}()

	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("mongodb is unreachable for integration test: %v", err)
	}

	dbName := fmt.Sprintf("cardflex_test_%d", time.Now().UnixNano())
	database := client.Database(dbName)
	defer func() {
		dropCtx, dropCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer dropCancel()
		_ = database.Drop(dropCtx)
	}()

	tenantCollection := database.Collection("tenants")
	userCollection := database.Collection("users")

	tenant := models.Tenant{
		ID:          primitive.NewObjectID(),
		Name:        "Acme Card",
		CompanyCode: "acme",
		ThemeColor:  "#0B6E4F",
	}

	if _, err := tenantCollection.InsertOne(ctx, tenant); err != nil {
		t.Fatalf("failed to seed tenant: %v", err)
	}
	t.Logf("seeded tenant: id=%s company=%s name=%s theme=%s", tenant.ID.Hex(), tenant.CompanyCode, tenant.Name, tenant.ThemeColor)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authController := &controllers.AuthController{
		Users:     userCollection,
		JWTSecret: "test-secret",
	}

	tenantAware := r.Group("/")
	tenantAware.Use(middleware.TenantResolver(tenantCollection))
	tenantAware.POST("/register", authController.Register)

	body := `{"name":"Alice","email":"alice@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register?company=acme", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, res.Code, res.Body.String())
	}
	t.Logf("register response: status=%d body=%s", res.Code, strings.TrimSpace(res.Body.String()))

	var registerResponse struct {
		Message string      `json:"message"`
		UserID  interface{} `json:"userId"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &registerResponse); err == nil {
		t.Logf("register parsed: message=%q userId=%v", registerResponse.Message, registerResponse.UserID)
	}

	queryCtx, queryCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer queryCancel()

	var inserted models.User
	err = userCollection.FindOne(queryCtx, bson.M{
		"tenantId": tenant.ID,
		"email":    "alice@example.com",
	}).Decode(&inserted)
	if err != nil {
		t.Fatalf("expected inserted user, find failed: %v", err)
	}

	if inserted.Name != "Alice" {
		t.Fatalf("expected name Alice, got %q", inserted.Name)
	}
	if inserted.TenantID != tenant.ID {
		t.Fatalf("expected tenant id %s, got %s", tenant.ID.Hex(), inserted.TenantID.Hex())
	}
	t.Logf("inserted user: id=%s name=%s email=%s tenantId=%s", inserted.ID.Hex(), inserted.Name, inserted.Email, inserted.TenantID.Hex())
	if inserted.Password == "secret123" {
		t.Fatal("expected password to be hashed, but plain text was stored")
	}
	t.Logf("stored password hash length: %d", len(inserted.Password))
	if err := bcrypt.CompareHashAndPassword([]byte(inserted.Password), []byte("secret123")); err != nil {
		t.Fatalf("stored password is not a valid hash for original password: %v", err)
	}
}
