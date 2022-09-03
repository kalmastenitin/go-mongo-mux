package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Email     string             `json:"email,omitempty" validate:"required"`
	Password  string             `json:"-"`
	Company   string             `json:"company,omitempty"`
	Role      string             `json:"role"`
	IsActive  bool               `json:"isactive"`
	TsCreated time.Time          `json:"created_on"`
	TsUpdated time.Time          `json:"updated_on"`
}

type UserSession struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	User         primitive.ObjectID `bson:"user,omitempty"`
	AccessToken  string             `json:"accesstoken"`
	RefreshToken string             `json:"refreshtoken"`
	TsCreated    time.Time          `json:"created_on"`
	TsUpdated    time.Time          `json:"updated_on"`
	UserAgent    string             `bson:"useragent,omitempty"`
}
