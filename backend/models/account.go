package models

type Account struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	UserID           uint    `gorm:"not null;uniqueIndex:idx_account_user_tenant" json:"userId"`
	TenantID         uint    `gorm:"not null;uniqueIndex:idx_account_user_tenant;index" json:"tenantId"`
	MaskedCardNumber string  `gorm:"not null" json:"maskedCardNumber"`
	CreditLimit      float64 `gorm:"not null" json:"creditLimit"`
	AvailableBalance float64 `gorm:"not null" json:"availableBalance"`
	Currency         string  `gorm:"not null;default:USD" json:"currency"`
	User             User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Tenant           Tenant  `gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
