package controllers

import (
	"net/http"

	"cardflex-backend/models"
	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func (d *DashboardController) GetDashboard(c *gin.Context) {
	tenantRaw, _ := c.Get("tenant")
	tenant := tenantRaw.(models.Tenant)

	c.JSON(http.StatusOK, gin.H{
		"tenant": gin.H{
			"name":        tenant.Name,
			"companyCode": tenant.CompanyCode,
			"themeColor":  tenant.ThemeColor,
		},
		"card": gin.H{
			"maskedCardNumber": "**** **** **** 4821",
			"creditLimit":      12000,
			"availableBalance": 8250,
			"currency":         "USD",
		},
		"transactions": []gin.H{
			{"date": "2026-02-14", "merchant": "Grocery Mart", "amount": -82.41, "status": "Posted"},
			{"date": "2026-02-12", "merchant": "Fuel Station", "amount": -47.10, "status": "Posted"},
			{"date": "2026-02-10", "merchant": "Card Payment", "amount": 500.00, "status": "Completed"},
		},
	})
}
