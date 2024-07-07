package sapics

import (
	"fmt"
	"time"

	"github.com/wzshiming/iplocation"
)

var (
	RouteViewsAsnDbipASNIPv4 = NewDB[ASN](sourceURL("asn", "ipv4"), 8, 12*time.Hour)
	RouteViewsAsnDbipASNIPv6 = NewDB[ASN](sourceURL("asn", "ipv6"), 8, 12*time.Hour)
	IptoasnASNIPv4           = NewDB[ASN](sourceURL("iptoasn-asn", "ipv4"), 8, 12*time.Hour)
	IptoasnASNIPv6           = NewDB[ASN](sourceURL("iptoasn-asn", "ipv6"), 8, 12*time.Hour)
	DbipASNIPv4              = NewDB[ASN](sourceURL("dbip-asn", "ipv4"), 8, 15*24*time.Hour)
	DbipASNIPv6              = NewDB[ASN](sourceURL("dbip-asn", "ipv6"), 8, 15*24*time.Hour)
	GeoLite2ASNIPv4          = NewDB[ASN](sourceURL("geolite2-asn", "ipv4"), 8, 2*24*time.Hour)
	GeoLite2ASNIPv6          = NewDB[ASN](sourceURL("geolite2-asn", "ipv6"), 8, 2*24*time.Hour)

	AsnCountryIPv4         = NewDB[Country](sourceURL("asn-country", "ipv4"), 8, 12*time.Hour)
	AsnCountryIPv6         = NewDB[Country](sourceURL("asn-country", "ipv6"), 8, 12*time.Hour)
	GeoAsnCountryIPv4      = NewDB[Country](sourceURL("geo-asn-country", "ipv4"), 8, 12*time.Hour)
	GeoAsnCountryIPv6      = NewDB[Country](sourceURL("geo-asn-country", "ipv6"), 8, 12*time.Hour)
	GeoWhoisAsnCountryIPv4 = NewDB[Country](sourceURL("geo-whois-asn-country", "ipv4"), 8, 12*time.Hour)
	GeoWhoisAsnCountryIPv6 = NewDB[Country](sourceURL("geo-whois-asn-country", "ipv6"), 8, 12*time.Hour)
	IptoasnCountryIPv4     = NewDB[Country](sourceURL("iptoasn-country", "ipv4"), 8, 12*time.Hour)
	IptoasnCountryIPv6     = NewDB[Country](sourceURL("iptoasn-country", "ipv6"), 8, 12*time.Hour)
	DbipCountryIPv4        = NewDB[Country](sourceURL("dbip-country", "ipv4"), 8, 15*24*time.Hour)
	DbipCountryIPv6        = NewDB[Country](sourceURL("dbip-country", "ipv6"), 8, 15*24*time.Hour)
	GeoLite2CountryIPv4    = NewDB[Country](sourceURL("geolite2-country", "ipv4"), 8, 2*24*time.Hour)
	GeoLite2CountryIPv6    = NewDB[Country](sourceURL("geolite2-country", "ipv6"), 8, 2*24*time.Hour)

	DbipCityIPv4     = NewDB[City](sourceURL("dbip-city", "ipv4"), 12, 15*24*time.Hour)
	DbipCityIPv6     = NewDB[City](sourceURL("dbip-city", "ipv6"), 12, 15*24*time.Hour)
	GeoLite2CityIPv4 = NewDB[City](sourceURL("geolite2-city", "ipv4"), 12, 2*24*time.Hour)
	GeoLite2CityIPv6 = NewDB[City](sourceURL("geolite2-city", "ipv6"), 12, 2*24*time.Hour)
)

func NewDB[T any](dbPath string, cacheLevel int, updateInterval time.Duration) *iplocation.LazyDB[T] {
	return iplocation.MustLazyDB[T](dbPath, ',', iplocation.DecRange, cacheLevel, updateInterval)
}

// sourceURL build the url of db file
// https://github.com/sapics/ip-location-db
func sourceURL(source string, v string) string {
	var ext = "csv"

	switch source {
	case "dbip-city", "geolite2-city":
		ext = "csv.gz"
	}

	return fmt.Sprintf("https://raw.githubusercontent.com/sapics/ip-location-db/main/%s/%s-%s-num.%s",
		source, source, v, ext,
	)
}
