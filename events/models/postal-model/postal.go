package postalmodel

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostalCreate struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UUID      string             `json:"uuid" bson:"uuid"`
	Name      string             `json:"name" bson:"name"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	Version   int                `json:"version" bson:"version"`
}
