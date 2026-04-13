package migrations_test

import (
	"testing"

	"cardflex-backend/migrations"
	"cardflex-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigrateAccountsCreatesTenantScopedAccountTable(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&models.Tenant{}, &models.User{}); err != nil {
		t.Fatalf("failed to migrate tenant/user schema: %v", err)
	}

	if err := migrations.MigrateAccounts(db); err != nil {
		t.Fatalf("failed to migrate accounts: %v", err)
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
		Name:     "Alice",
		Email:    "alice@example.com",
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
		t.Fatalf("failed to create account: %v", err)
	}

	var stored models.Account
	if err := db.Where("user_id = ? AND tenant_id = ?", user.ID, tenant.ID).First(&stored).Error; err != nil {
		t.Fatalf("failed to load stored account: %v", err)
	}

	if stored.MaskedCardNumber != account.MaskedCardNumber {
		t.Fatalf("expected masked card number %q, got %q", account.MaskedCardNumber, stored.MaskedCardNumber)
	}
	if stored.TenantID != tenant.ID {
		t.Fatalf("expected tenant id %d, got %d", tenant.ID, stored.TenantID)
	}
	if stored.UserID != user.ID {
		t.Fatalf("expected user id %d, got %d", user.ID, stored.UserID)
	}
}
