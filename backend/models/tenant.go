package models

type Tenant struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	CompanyCode string `gorm:"not null;uniqueIndex" json:"companyCode"`
	ThemeColor  string `json:"themeColor,omitempty"`
}
