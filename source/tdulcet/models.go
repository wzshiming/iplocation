package tdulcet

import (
	"encoding/json"
)

type Country struct {
	// country_code is the two-letter code defined in ISO 3166-1 alpha-2.
	// https://wikipedia.org/wiki/ISO_3166-1_alpha-2
	CountryCode string `json:"countryCode"`
}

type City1 struct {
	CountryCode string      `json:"countryCode"`
	State1      string      `json:"state1"`
	City        string      `json:"city"`
	Latitude    json.Number `json:"latitude"`
	Longitude   json.Number `json:"longitude"`
}

type City2 struct {
	CountryCode string      `json:"countryCode"`
	State1      string      `json:"state1"`
	State2      string      `json:"State2,omitempty"`
	City        string      `json:"city"`
	Latitude    json.Number `json:"latitude"`
	Longitude   json.Number `json:"longitude"`
}
