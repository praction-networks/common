package authevent

type AuthUserLoginEvent struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	Version    int    `json:"version" bson:"version"`
}
