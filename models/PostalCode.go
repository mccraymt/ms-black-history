package models

// PostalCode ...
type PostalCode struct {
	ID                string           `json:"id"`
	Value             string           `json:"value"`
	HasRatedLocations bool             `json:"hasRatedLocations"`
	PostalCodeType    *string          `json:"postalCodeType"`
	City              *string          `json:"city"`
	State             *string          `json:"state"`
	CityAbbreviation  *string          `json:"cityAbbreviation,omitempty"`
	RatedLocations    []*RatedLocation `json:"-"`
	Counties          []*County        `json:"counties,omitempty"`
}

// RatedLocation ...
type RatedLocation struct {
	PostalCode       string `json:"postalCode"`
	State            string `json:"state"`
	County           string `json:"county"`
	City             string `json:"city"`
	CityAbbreviation string `json:"cityAbbreviation"`
}

// County ...
type County struct {
	Name           string  `json:"name"`
	FipsCode       string  `json:"-"`
	State          string  `json:"state"`
	CountyNameType *string `json:"-"`
}
