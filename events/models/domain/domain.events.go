package domainEvenetHandler

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type Department struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}
