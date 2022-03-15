package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ObjectID      primitive.ObjectID  `bson:"_id" json:"_id"`
	Username      string              `json:"username" gorm:"unique" bson:"username,omitempty"`
	Email         string              `json:"email" gorm:"unique" bson:"email,omitempty"`
	Password      string              `json:"password" bson:"password"`
	CreatedAt     primitive.Timestamp `json:"createdat" bson:"createat"`
	DeactivatedAt primitive.Timestamp `json:"updatedat" bson:"updatedat"`
}
