package authevent

type AuthUserUpdateEvent struct {
	ID      string `json:"id,omitempty" bson:"_id,omitempty"`
	Version int    `json:"version" bson:"version"`
}
