package migrations

import (
	"fmt"
	"time"

	"cardflex-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type sampleAccountSeed struct {
	MaskedCardNumber string
	CreditLimit      float64
	AvailableBalance float64
	Currency         string
}

type sampleTransactionSeed struct {
	Date     time.Time
	Merchant string
	Amount   float64
	Status   string
}

func MigrateTransactions(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		return err
	}

	return seedSampleDashboardData(db)
}

func seedSampleDashboardData(db *gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	accountSeeds := map[string]sampleAccountSeed{
		"acme":        {MaskedCardNumber: "**** **** **** 4821", CreditLimit: 12000, AvailableBalance: 8250, Currency: "USD"},
		"nova":        {MaskedCardNumber: "**** **** **** 7134", CreditLimit: 9800, AvailableBalance: 6410, Currency: "USD"},
		"prime":       {MaskedCardNumber: "**** **** **** 9150", CreditLimit: 15000, AvailableBalance: 11025, Currency: "USD"},
		"chase-bank":  {MaskedCardNumber: "**** **** **** 2468", CreditLimit: 18000, AvailableBalance: 12760, Currency: "USD"},
		"wells-fargo": {MaskedCardNumber: "**** **** **** 5531", CreditLimit: 10500, AvailableBalance: 7025, Currency: "USD"},
		"capital-one": {MaskedCardNumber: "**** **** **** 8842", CreditLimit: 13500, AvailableBalance: 9360, Currency: "USD"},
	}

	transactionSeeds := map[string][]sampleTransactionSeed{
		"acme": {
			{Date: time.Date(2026, time.February, 14, 0, 0, 0, 0, time.UTC), Merchant: "Grocery Mart", Amount: -82.41, Status: "Posted"},
			{Date: time.Date(2026, time.February, 12, 0, 0, 0, 0, time.UTC), Merchant: "Fuel Station", Amount: -47.10, Status: "Posted"},
			{Date: time.Date(2026, time.February, 10, 0, 0, 0, 0, time.UTC), Merchant: "Card Payment", Amount: 500.00, Status: "Completed"},
		},
		"nova": {
			{Date: time.Date(2026, time.February, 15, 0, 0, 0, 0, time.UTC), Merchant: "Fresh Market", Amount: -64.89, Status: "Posted"},
			{Date: time.Date(2026, time.February, 11, 0, 0, 0, 0, time.UTC), Merchant: "Streaming Plus", Amount: -18.99, Status: "Posted"},
			{Date: time.Date(2026, time.February, 8, 0, 0, 0, 0, time.UTC), Merchant: "Card Payment", Amount: 350.00, Status: "Completed"},
		},
		"prime": {
			{Date: time.Date(2026, time.February, 13, 0, 0, 0, 0, time.UTC), Merchant: "Airline Booking", Amount: -426.55, Status: "Posted"},
			{Date: time.Date(2026, time.February, 9, 0, 0, 0, 0, time.UTC), Merchant: "Hotel Stay", Amount: -219.40, Status: "Posted"},
			{Date: time.Date(2026, time.February, 4, 0, 0, 0, 0, time.UTC), Merchant: "Card Payment", Amount: 900.00, Status: "Completed"},
		},
		"chase-bank": {
			{Date: time.Date(2026, time.February, 16, 0, 0, 0, 0, time.UTC), Merchant: "Online Retail", Amount: -129.99, Status: "Posted"},
			{Date: time.Date(2026, time.February, 12, 0, 0, 0, 0, time.UTC), Merchant: "Dining Hub", Amount: -58.22, Status: "Posted"},
			{Date: time.Date(2026, time.February, 7, 0, 0, 0, 0, time.UTC), Merchant: "Card Payment", Amount: 620.00, Status: "Completed"},
		},
		"wells-fargo": {
			{Date: time.Date(2026, time.February, 14, 0, 0, 0, 0, time.UTC), Merchant: "Pharmacy Care", Amount: -36.70, Status: "Posted"},
			{Date: time.Date(2026, time.February, 10, 0, 0, 0, 0, time.UTC), Merchant: "Utility Bill", Amount: -112.08, Status: "Posted"},
			{Date: time.Date(2026, time.February, 3, 0, 0, 0, 0, time.UTC), Merchant: "Card Payment", Amount: 410.00, Status: "Completed"},
		},
		"capital-one": {
			{Date: time.Date(2026, time.February, 14, 0, 0, 0, 0, time.UTC), Merchant: "Grocery Mart", Amount: -82.41, Status: "Posted"},
			{Date: time.Date(2026, time.February, 12, 0, 0, 0, 0, time.UTC), Merchant: "Fuel Station", Amount: -47.10, Status: "Posted"},
			{Date: time.Date(2026, time.February, 10, 0, 0, 0, 0, time.UTC), Merchant: "Card Payment", Amount: 500.00, Status: "Completed"},
		},
	}

	var tenants []models.Tenant
	if err := db.Find(&tenants).Error; err != nil {
		return err
	}

	for _, tenant := range tenants {
		accountSeed, ok := accountSeeds[tenant.CompanyCode]
		if !ok {
			accountSeed = sampleAccountSeed{
				MaskedCardNumber: "**** **** **** 0001",
				CreditLimit:      10000,
				AvailableBalance: 7500,
				Currency:         "USD",
			}
		}

		user := models.User{
			Name:     fmt.Sprintf("%s Demo User", tenant.Name),
			Email:    fmt.Sprintf("demo+%s@cardflex.local", tenant.CompanyCode),
			Password: string(passwordHash),
			TenantID: tenant.ID,
		}
		if err := db.
			Where(models.User{TenantID: tenant.ID, Email: user.Email}).
			Attrs(models.User{Name: user.Name, Password: user.Password}).
			FirstOrCreate(&user).Error; err != nil {
			return err
		}

		account := models.Account{
			UserID:           user.ID,
			TenantID:         tenant.ID,
			MaskedCardNumber: accountSeed.MaskedCardNumber,
			CreditLimit:      accountSeed.CreditLimit,
			AvailableBalance: accountSeed.AvailableBalance,
			Currency:         accountSeed.Currency,
		}
		if err := db.
			Where(models.Account{UserID: user.ID, TenantID: tenant.ID}).
			Assign(models.Account{
				MaskedCardNumber: account.MaskedCardNumber,
				CreditLimit:      account.CreditLimit,
				AvailableBalance: account.AvailableBalance,
				Currency:         account.Currency,
			}).
			FirstOrCreate(&account).Error; err != nil {
			return err
		}

		var existingCount int64
		if err := db.Model(&models.Transaction{}).
			Where("account_id = ? AND user_id = ? AND tenant_id = ?", account.ID, user.ID, tenant.ID).
			Count(&existingCount).Error; err != nil {
			return err
		}
		if existingCount > 0 {
			continue
		}

		seeds := transactionSeeds[tenant.CompanyCode]
		if len(seeds) == 0 {
			seeds = []sampleTransactionSeed{
				{Date: time.Date(2026, time.February, 14, 0, 0, 0, 0, time.UTC), Merchant: "Welcome Purchase", Amount: -45.50, Status: "Posted"},
				{Date: time.Date(2026, time.February, 10, 0, 0, 0, 0, time.UTC), Merchant: "Card Payment", Amount: 250.00, Status: "Completed"},
			}
		}

		transactions := make([]models.Transaction, 0, len(seeds))
		for _, seed := range seeds {
			transactions = append(transactions, models.Transaction{
				AccountID: account.ID,
				UserID:    user.ID,
				TenantID:  tenant.ID,
				Date:      seed.Date,
				Merchant:  seed.Merchant,
				Amount:    seed.Amount,
				Status:    seed.Status,
			})
		}

		if err := db.Create(&transactions).Error; err != nil {
			return err
		}
	}

	return nil
}
