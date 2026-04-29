package controllers

import (
	"net/http"
	"strconv"
	"time"

	"cardflex-backend/models"
	"cardflex-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentController struct {
	DB *gorm.DB
}

type paymentRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type paymentResponse struct {
	Message        string  `json:"message"`
	UpdatedBalance float64 `json:"updatedBalance"`
	TransactionID  uint    `json:"transactionId"`
	Amount         float64 `json:"amount"`
	Timestamp      string  `json:"timestamp"`
}

func (p *PaymentController) RecordPayment(c *gin.Context) {
	if p.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "payment database is not configured"})
		return
	}

	// Get tenant from context
	tenantRaw, ok := c.Get("tenant")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context missing"})
		return
	}

	tenant, ok := tenantRaw.(models.Tenant)
	if !ok || tenant.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid tenant context"})
		return
	}

	// Get claims from context
	claimsRaw, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication claims missing"})
		return
	}

	claims, ok := claimsRaw.(*utils.Claims)
	if !ok || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication claims"})
		return
	}

	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token user"})
		return
	}
	if !tenant.Features.Enabled("payments_enabled") {
		c.JSON(http.StatusForbidden, gin.H{"error": "payments are disabled for this tenant"})
		return
	}

	// Parse request body
	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate amount
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be a positive number"})
		return
	}

	// Fetch account
	var account models.Account
	if err := p.DB.
		Where("tenant_id = ? AND user_id = ?", tenant.ID, uint(userID)).
		First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load account"})
		return
	}

	// Validate amount does not exceed current balance
	if req.Amount > account.AvailableBalance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment amount exceeds available balance"})
		return
	}

	// Begin transaction
	tx := p.DB.Begin()

	// Update account balance
	newBalance := account.AvailableBalance - req.Amount
	if err := tx.Model(&account).Update("available_balance", newBalance).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update account balance"})
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		AccountID: account.ID,
		UserID:    uint(userID),
		TenantID:  tenant.ID,
		Date:      time.Now(),
		Merchant:  "Card Payment",
		Amount:    -req.Amount, // Negative amount for payment/credit
		Status:    "Posted",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record transaction"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit payment"})
		return
	}

	// Return response
	c.JSON(http.StatusOK, paymentResponse{
		Message:        "payment recorded successfully",
		UpdatedBalance: newBalance,
		TransactionID:  transaction.ID,
		Amount:         req.Amount,
		Timestamp:      transaction.Date.Format("2006-01-02T15:04:05Z"),
	})
}
