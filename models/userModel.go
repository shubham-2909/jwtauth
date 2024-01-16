package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json:"name" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required"`
	Email         *string            `json:"email" validate:"email,required"`
	Token         *string            `json:"token"`
	UserType      *string            `json:"usertype" validate:"required,eq=ADMIN|eq=USER"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}
