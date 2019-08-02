package models

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

type UserAccount struct {
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    string        `json:"first_name" bson:"first_name" validate:"required"`
	LastName     string        `json:"last_name" bson:"last_name" validate:"required"`
	EmailAddress string        `json:"email_address" bson:"email_address" validaet:"required"`
	Role         UserRole      `json:"role" bson:"role" validate:"required"`
	IsBanned     bool          `json:"is_banned" bson:"is_banned"`
	Password     string        `json:"password,omitempty" bson:"password,omitempty"`
	Needs        string        `json:"needs" bson:"needs"`
}

type UserRole struct {
	RoleName string `json:"role_name" bson:"role_name"`
}
