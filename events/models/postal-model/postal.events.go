package postalEvenetModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type PostalCreate struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UUID      string             `json:"uuid" bson:"uuid"`
	Name      string             `json:"name" bson:"name"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	Version   int                `json:"version" bson:"version"`
}

type PostalUpdate struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UUID      string             `json:"uuid" bson:"uuid"`
	Name      string             `json:"name" bson:"name"`
	IsActive  bool               `json:"isActive" bson:"isActive"`
	UpdatedAt primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
	Version   int                `json:"version" bson:"version"`
}

type PostalDelete struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	UUID    string             `json:"uuid" bson:"uuid"`
	Name    string             `json:"name" bson:"name"`
	Version int                `json:"version" bson:"version"`
}
