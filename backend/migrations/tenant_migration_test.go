package migrations_test

import (
	"testing"

	"cardflex-backend/migrations"
	"cardflex-backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigrateTenantsSeedsFeatureFlags(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}

	if err := migrations.MigrateTenants(db); err != nil {
		t.Fatalf("failed to migrate tenants: %v", err)
	}

	expected := map[string]models.FeatureFlags{
		"chase-bank": {
			"payments_enabled": true,
			"profile_enabled":  true,
		},
		"capital-one": {
			"payments_enabled": true,
			"profile_enabled":  false,
		},
		"wells-fargo": {
			"payments_enabled": false,
			"profile_enabled":  true,
		},
	}

	for companyCode, features := range expected {
		var tenant models.Tenant
		if err := db.Where("company_code = ?", companyCode).First(&tenant).Error; err != nil {
			t.Fatalf("failed to load seeded tenant %q: %v", companyCode, err)
		}

		for feature, enabled := range features {
			if tenant.Features[feature] != enabled {
				t.Fatalf("expected %s %s to be %t, got %t", companyCode, feature, enabled, tenant.Features[feature])
			}
		}
	}
}
