package migrations

import (
	"cardflex-backend/models"
	"gorm.io/gorm"
)

func MigrateTenants(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Tenant{}); err != nil {
		return err
	}

	sampleTenants := []models.Tenant{
		{Name: "Acme Card", CompanyCode: "acme", ThemeColor: "#0B6E4F", Features: defaultTenantFeatures()},
		{Name: "Nova Finance", CompanyCode: "nova", ThemeColor: "#C84B31", Features: defaultTenantFeatures()},
		{Name: "Prime Credit", CompanyCode: "prime", ThemeColor: "#00539C", Features: defaultTenantFeatures()},
		{Name: "Chase Bank", CompanyCode: "chase-bank", ThemeColor: "#0A2A66", Features: models.FeatureFlags{
			"payments_enabled": true,
			"profile_enabled":  true,
		}},
		{Name: "Wells Fargo", CompanyCode: "wells-fargo", ThemeColor: "#B31B1B", Features: models.FeatureFlags{
			"payments_enabled": false,
			"profile_enabled":  true,
		}},
		{Name: "Capital One", CompanyCode: "capital-one", ThemeColor: "#003B95", Features: models.FeatureFlags{
			"payments_enabled": true,
			"profile_enabled":  false,
		}},
	}

	for _, tenant := range sampleTenants {
		if err := db.
			Where(models.Tenant{CompanyCode: tenant.CompanyCode}).
			Assign(models.Tenant{Name: tenant.Name, ThemeColor: tenant.ThemeColor, Features: tenant.Features}).
			FirstOrCreate(&tenant).Error; err != nil {
			return err
		}
	}

	return nil
}

func defaultTenantFeatures() models.FeatureFlags {
	return models.FeatureFlags{
		"payments_enabled": true,
		"profile_enabled":  true,
	}
}
