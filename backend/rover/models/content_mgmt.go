package models

import (
	"time"

	"github.com/fidellr/jastip/backend/uranus/models"
	"github.com/globalsign/mgo/bson"
)

type Content struct {
	CreatedAt time.Time       `json:"created_at,omiempty" bson:"created_at,omiempty"`
	UpdateAt  time.Time       `json:"updated_at,omiempty" bson:"updated_at,omiempty"`
	ID        bson.ObjectId   `json:"id,omitempty" bson:"_id,omitempty"`
	Screen    string          `json:"screen" bson:"screen" validate:"required"`
	Role      models.UserRole `json:"role" bson:"role"`
	Items     []Item          `json:"items" bson:"items"`
}

type Item struct {
	ID      bson.ObjectId          `json:"id,omiempty" bson:"_id,omiempty"`
	Type    string                 `json:"type" bson:"type" validate:"required"`
	Content map[string]interface{} `json:"content" bson:"content"`
	Layout  string                 `json:"layout" bson:"layout" validate:"required"`
}

type HomeContent struct {
}
