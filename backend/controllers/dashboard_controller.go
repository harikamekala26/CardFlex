package controllers

import (
	"net/http"
	"strconv"

	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardController struct {
	DB *gorm.DB
}

func (d *DashboardController) GetDashboard(c *gin.Context) {
	tenantRaw, _ := c.Get("tenant")
	tenant := tenantRaw.(models.Tenant)

	claimsRaw, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication claims missing"})
		return
	}

	claims := claimsRaw.(*utils.Claims)
	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user"})
		return
	}

	var account models.Account
	if err := d.DB.
		Where("tenant_id = ? AND user_id = ?", tenant.ID, uint(userID)).
		First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load account"})
		return
	}

	var transactions []models.Transaction
	if err := d.DB.
		Where("tenant_id = ? AND user_id = ? AND account_id = ?", tenant.ID, uint(userID), account.ID).
		Order("date DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load transactions"})
		return
	}

	transactionResponse := make([]gin.H, 0, len(transactions))
	for _, transaction := range transactions {
		transactionResponse = append(transactionResponse, gin.H{
			"date":     transaction.Date.Format("2006-01-02"),
			"merchant": transaction.Merchant,
			"amount":   transaction.Amount,
			"status":   transaction.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tenant": gin.H{
			"name":        tenant.Name,
			"companyCode": tenant.CompanyCode,
			"themeColor":  tenant.ThemeColor,
		},
		"card": gin.H{
			"maskedCardNumber": account.MaskedCardNumber,
			"creditLimit":      account.CreditLimit,
			"availableBalance": account.AvailableBalance,
			"currency":         account.Currency,
		},
		"transactions": transactionResponse,
	})
}
