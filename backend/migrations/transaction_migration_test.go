package migrations_test

import (
	"testing"

	"cardflex-backend/migrations"
	"cardflex-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigrateTransactionsSeedsSampleDashboardData(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := migrations.MigrateTenants(db); err != nil {
		t.Fatalf("failed to migrate tenants: %v", err)
	}
	if err := migrations.MigrateUsers(db); err != nil {
		t.Fatalf("failed to migrate users: %v", err)
	}
	if err := migrations.MigrateAccounts(db); err != nil {
		t.Fatalf("failed to migrate accounts: %v", err)
	}
	if err := migrations.MigrateTransactions(db); err != nil {
		t.Fatalf("failed to migrate transactions: %v", err)
	}
	if err := migrations.MigrateTransactions(db); err != nil {
		t.Fatalf("failed to rerun transaction migration: %v", err)
	}

	var tenant models.Tenant
	if err := db.Where("company_code = ?", "capital-one").First(&tenant).Error; err != nil {
		t.Fatalf("failed to load seeded tenant: %v", err)
	}

	var user models.User
	if err := db.Where("tenant_id = ?", tenant.ID).First(&user).Error; err != nil {
		t.Fatalf("failed to load seeded user: %v", err)
	}

	var account models.Account
	if err := db.Where("tenant_id = ? AND user_id = ?", tenant.ID, user.ID).First(&account).Error; err != nil {
		t.Fatalf("failed to load seeded account: %v", err)
	}

	var transactions []models.Transaction
	if err := db.Where("tenant_id = ? AND user_id = ? AND account_id = ?", tenant.ID, user.ID, account.ID).
		Order("date desc").
		Find(&transactions).Error; err != nil {
		t.Fatalf("failed to load seeded transactions: %v", err)
	}

	if account.MaskedCardNumber == "" {
		t.Fatal("expected seeded account to include masked card number")
	}
	if len(transactions) != 3 {
		t.Fatalf("expected 3 seeded transactions for capital-one, got %d", len(transactions))
	}
	if transactions[0].Merchant == "" || transactions[0].Status == "" {
		t.Fatal("expected seeded transactions to include merchant and status")
	}
}
