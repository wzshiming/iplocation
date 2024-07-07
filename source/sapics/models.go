package sapics

import "encoding/json"

type Country struct {
	// country_code is the two-letter code defined in ISO 3166-1 alpha-2.
	// You can get the country name, capital, continent, currency, languages,
	// etc. from the country_code by Countries Database in JSON, CSV, SQL format.
	// https://wikipedia.org/wiki/ISO_3166-1_alpha-2
	// https://github.com/annexare/Countries
	CountryCode string `json:"countryCode"`
}

type City struct {
	CountryCode string      `json:"countryCode"`
	State1      string      `json:"state1"`
	State2      string      `json:"state2,omitempty"`
	City        string      `json:"city"`
	Postcode    string      `json:"postcode,omitempty"`
	Latitude    json.Number `json:"latitude"`
	Longitude   json.Number `json:"longitude"`
	Timezone    string      `json:"timezone,omitempty"`
}

type ASN struct {
	// as_number is a unique number assigned to an Autonomous System (AS) by the IANA.
	// https://www.iana.org/
	// https://wikipedia.org/wiki/Autonomous_system_(Internet)
	AsNumber       string `json:"asNumber"`
	AsOrganization string `json:"asOrganization"`
}
