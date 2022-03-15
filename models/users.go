package models

import (
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//User is the structure we work with
type User struct {
	ID            primitive.ObjectID  `bson:"_id" json:"_id"`
	Username      string              `json:"username" gorm:"unique" bson:"username,omitempty"`
	Email         string              `json:"email" gorm:"unique" bson:"email,omitempty"`
	Password      string              `json:"password" bson:"password"`
	CreatedAt     primitive.Timestamp `json:"createdat" bson:"createat"`
	DeactivatedAt primitive.Timestamp `json:"updatedat" bson:"updatedat"`
}

// NewUser create new instance of User
func NewUser(username string, email string, password string, createdAt time.Time, deactivatedAt time.Time) User {
	user := User{}
	user.ID = primitive.NewObjectID()
	user.Username = username
	user.Email = email
	user.Password = base64.StdEncoding.EncodeToString([]byte(password))
	user.CreatedAt = primitive.Timestamp{T: uint32(createdAt.Unix())}
	user.DeactivatedAt = primitive.Timestamp{T: uint32(deactivatedAt.Unix())}
	return user
}

func AddDate(years int, months int, days int) time.Time {
	return time.Now().AddDate(years, months, days)
}

//type User struct {
//	ID            primitive.ObjectID  `bson:"_id" json:"_id" default:primitive.NewObjectID()`
//	Username      string              `json:"username" gorm:"unique" bson:"username,omitempty"`
//	Email         string              `json:"email" gorm:"unique" bson:"email,omitempty"`
//	Password      string              `json:"password" bson:"password"`
//	CreatedAt     primitive.Timestamp `json:"createdat" bson:"createat"`
//	DeactivatedAt primitive.Timestamp `json:"updatedat" bson:"updatedat"`
//}
