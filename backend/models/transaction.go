package models

import "time"

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AccountID uint      `gorm:"not null;index" json:"accountId"`
	UserID    uint      `gorm:"not null;index" json:"userId"`
	TenantID  uint      `gorm:"not null;index" json:"tenantId"`
	Date      time.Time `gorm:"not null;index" json:"date"`
	Merchant  string    `gorm:"not null" json:"merchant"`
	Amount    float64   `gorm:"not null" json:"amount"`
	Status    string    `gorm:"not null" json:"status"`
	Account   Account   `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
