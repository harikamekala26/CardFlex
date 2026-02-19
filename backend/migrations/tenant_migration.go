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
		{Name: "Acme Card", CompanyCode: "acme", ThemeColor: "#0B6E4F"},
		{Name: "Nova Finance", CompanyCode: "nova", ThemeColor: "#C84B31"},
		{Name: "Prime Credit", CompanyCode: "prime", ThemeColor: "#00539C"},
		{Name: "Chase Bank", CompanyCode: "chase-bank", ThemeColor: "#0A2A66"},
		{Name: "Wells Fargo", CompanyCode: "wells-fargo", ThemeColor: "#B31B1B"},
		{Name: "Capital One", CompanyCode: "capital-one", ThemeColor: "#003B95"},
	}

	for _, tenant := range sampleTenants {
		if err := db.
			Where(models.Tenant{CompanyCode: tenant.CompanyCode}).
			Assign(models.Tenant{Name: tenant.Name, ThemeColor: tenant.ThemeColor}).
			FirstOrCreate(&tenant).Error; err != nil {
			return err
		}
	}

	return nil
}
