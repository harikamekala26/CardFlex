package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tenant struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	CompanyCode string             `bson:"companyCode" json:"companyCode"`
	ThemeColor  string             `bson:"themeColor" json:"themeColor"`
}
