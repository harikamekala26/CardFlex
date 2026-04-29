package models

import "strings"

type FeatureFlags map[string]bool

type Tenant struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"not null" json:"name"`
	CompanyCode string       `gorm:"not null;uniqueIndex" json:"companyCode"`
	ThemeColor  string       `json:"themeColor,omitempty"`
	Features    FeatureFlags `gorm:"type:json;serializer:json" json:"features"`
}

func (f FeatureFlags) Enabled(name string) bool {
	enabled, ok := f[name]
	return ok && enabled
}

func (f FeatureFlags) ToCamelCaseMap() map[string]bool {
	features := make(map[string]bool, len(f))
	for name, enabled := range f {
		features[snakeToLowerCamel(name)] = enabled
	}

	return features
}

func snakeToLowerCamel(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) == 0 {
		return value
	}

	result := parts[0]
	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		result += strings.ToUpper(part[:1]) + part[1:]
	}

	return result
}
