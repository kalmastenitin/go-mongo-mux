package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Email     string             `json:"email,omitempty" validate:"required"`
	Company   string             `json:"company,omitempty"`
	IsActive  bool               `json:"isactive"`
	TsCreated time.Time          `json:"created_on"`
	TsUpdated time.Time          `json:"updated_on"`
}
