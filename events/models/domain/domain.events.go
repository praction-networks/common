package domainEvenetModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type DomainCreate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type DomainUpdate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type DomainDelete struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

type DepartmentCreate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type DepartmentUpdate struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UUID        string             `json:"uuid" bson:"uuid,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	SystemName  string             `json:"systemName" bson:"systemName,omitempty"`
	ParentRefID primitive.ObjectID `json:"parentRefId,omitempty" bson:"parentRefId,omitempty"` // Nullable parent reference
	Version     int                `json:"version" bson:"version"`
}

type DepartmentDelete struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}
