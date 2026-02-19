package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"not null;index:idx_tenant_email,unique" json:"email"`
	Password string `gorm:"not null" json:"-"`
	TenantID uint   `gorm:"not null;index:idx_tenant_email,unique" json:"tenantId"`
}
