package tenantevent

// AddressModel represents an address structure
// @swagger:model AddressModel
type AddressModel struct {
	HouseNumber string `json:"houseNumber" bson:"houseNumber" validate:"required" example:"123"`
	Address1    string `json:"address1" bson:"address1" validate:"required" example:"Street 1, Block A"`
	Address2    string `json:"address2,omitempty" bson:"address2,omitempty" example:"Near Park"`
	Pincode     int    `json:"pincode" bson:"pincode" validate:"required,pincode" example:"201017"`
	City        string `json:"city" bson:"city" validate:"required" example:"Ghaziabad"`
	State       string `json:"state" bson:"state" validate:"required" example:"Uttar Pradesh"`
	Country     string `json:"country" bson:"country" validate:"required" example:"India"`
	LandMark    string `json:"landMark,omitempty" bson:"landMark,omitempty" example:"Near NG Sons Road"`
}
