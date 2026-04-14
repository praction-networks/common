package provider

// radiusVendorAdtran holds the RADIUS attribute dictionary for ADTRAN.
// It is registered into RadiusDictionaryRegistry in radius.attribute.registry.go.
var radiusVendorAdtran = VendorDictionaryInfo{
	Value:       "ADTRAN",
	Label:       "ADTRAN",
	Description: "ADTRAN RADIUS VSA (Vendor ID: 664). Used with ADTRAN Total Access and SDX OLT platforms in US ISP DSL and fibre deployments.",
	VendorID:    664,
	Attributes: []RadiusAttributeSchema{
		{
			Key:         "ADTRAN-User-Group",
			AttributeID:    1,
			Label:       "User Group",
			Description: "User group name for admin management access control",
			DataType:    RadiusDataTypeString,
			Category:    RadiusCategoryAuthorization,
			MaxLength:   intPtr(64),
			Examples:    []string{"admin-group", "read-only-group"},
		},
		{
			Key:         "ADTRAN-Privilege-Level",
			AttributeID:    6,
			Label:       "Privilege Level",
			Description: "CLI privilege level granted to the administrator (0 = lowest, 15 = highest)",
			DataType:    RadiusDataTypeInteger,
			Category:    RadiusCategoryAuthorization,
			MinValue:    float64Ptr(0),
			MaxValue:    float64Ptr(15),
			Examples:    []string{"15", "7", "1"},
		},
		{
			Key:         "ADTRAN-Rate-Limit-Up",
			AttributeID:    4,
			Label:       "Rate Limit Up",
			Description: "Maximum upstream bandwidth for the subscriber in kbps",
			DataType:    RadiusDataTypeInteger,
			Category:    RadiusCategoryQoS,
			MinValue:    float64Ptr(0),
			Examples:    []string{"10240", "51200"},
		},
		{
			Key:         "ADTRAN-Rate-Limit-Down",
			AttributeID:    3,
			Label:       "Rate Limit Down",
			Description: "Maximum downstream bandwidth for the subscriber in kbps",
			DataType:    RadiusDataTypeInteger,
			Category:    RadiusCategoryQoS,
			MinValue:    float64Ptr(0),
			Examples:    []string{"51200", "204800"},
		},
		{
			Key:         "ADTRAN-VLAN-ID",
			AttributeID:    5,
			Label:       "VLAN ID",
			Description: "VLAN ID assigned to the subscriber's access port",
			DataType:    RadiusDataTypeInteger,
			Category:    RadiusCategoryTunneling,
			MinValue:    float64Ptr(1),
			MaxValue:    float64Ptr(4094),
			Examples:    []string{"100", "200"},
		},
		{
			Key:         "ADTRAN-IP-Pool",
			AttributeID:    7,
			Label:       "IP Pool",
			Description: "IP address pool name for subscriber address allocation",
			DataType:    RadiusDataTypeString,
			Category:    RadiusCategoryAuthorization,
			MaxLength:   intPtr(64),
			Examples:    []string{"pool-dsl", "pool-fiber"},
		},
		{
			Key:         "ADTRAN-Service-Profile",
			AttributeID:    2,
			Label:       "Service Profile",
			Description: "Service profile name defining the subscriber access parameters",
			DataType:    RadiusDataTypeString,
			Category:    RadiusCategoryAuthorization,
			MaxLength:   intPtr(64),
			Examples:    []string{"100M-Residential", "1G-Business"},
		},
	},
}

// -----------------------------------------------------------------------

