package migrations

import (
	"cardflex-backend/models"
	"gorm.io/gorm"
)

func MigrateAccounts(db *gorm.DB) error {
	return db.AutoMigrate(&models.Account{})
}
