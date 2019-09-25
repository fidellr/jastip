package models

import (
	"time"

	"github.com/fidellr/jastip/backend/uranus/models"
	"github.com/globalsign/mgo/bson"
)

type Screen struct {
	CreatedAt  time.Time       `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt   time.Time       `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	ID         bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty"`
	ScreenName string          `json:"screen_name" bson:"screen_name" validate:"required"`
	Role       models.UserRole `json:"role" bson:"role"`
	Items      []Item          `json:"items" bson:"items"`
}

type Item struct {
	ID      bson.ObjectId          `json:"id,omitempty" bson:"_id,omitempty"`
	Type    string                 `json:"type" bson:"type" validate:"required"`
	Content map[string]interface{} `json:"content" bson:"content"`
	Layout  string                 `json:"layout" bson:"layout" validate:"required"`
}
