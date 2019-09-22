package models

import (
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/fidellr/jastip/backend/uranus/models"
)

type Image struct {
	ID         bson.ObjectId    `json:"id,omitempty" bson:"_id,omitempty"`
	PersonName string           `json:"person_name" bson:"person_name"`
	CreatedAt  time.Time        `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Needs      string           `json:"needs" bson:"needs" validate:"required"`
	Role       *models.UserRole `json:"role_name,omitempty" bson:"role_name,omitempty"`
	MIME       Mime             `json:"mime,omitempty" bson:"mime,omitempty"`
	FileLink   string           `json:"file,omitempty" bson:"file,omitempty"`
	Width      int              `json:"width,omitempty" bson:"width,omitempty"`
	Height     int              `json:"height,omitempty" bson:"height,omitempty"`
}

type Mime struct {
	Type string `json:"type,omitempty" bson:"type,omitempty"`
}
