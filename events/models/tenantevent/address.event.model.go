package tenantevent

// AddressModel represents an address structure
// @swagger:model AddressModel
type AddressModel struct {
	HouseNumber string `json:"houseNumber" bson:"houseNumber"`
	Address1    string `json:"address1" bson:"address1"`
	Address2    string `json:"address2,omitempty" bson:"address2,omitempty"`
	Pincode     string `json:"pincode" bson:"pincode"`
	City        string `json:"city" bson:"city"`
	State       string `json:"state" bson:"state"`
	Country     string `json:"country" bson:"country"`
	LandMark    string `json:"landMark,omitempty" bson:"landMark,omitempty"`
}
