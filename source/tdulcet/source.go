package tdulcet

import (
	"fmt"
	"time"

	"github.com/wzshiming/iplocation"
)

var (
	GeoWhoisAsnCountryIPv4 = NewDB[Country](sourceURL("geo-whois-asn-country", "ipv4"), 8, 12*time.Hour)
	GeoWhoisAsnCountryIPv6 = NewDB[Country](sourceURL("geo-whois-asn-country", "ipv6"), 8, 12*time.Hour)
	IptoasnCountryIPv4     = NewDB[Country](sourceURL("iptoasn-country", "ipv4"), 8, 12*time.Hour)
	IptoasnCountryIPv6     = NewDB[Country](sourceURL("iptoasn-country", "ipv6"), 8, 12*time.Hour)
	IpinfoCountryIPv4      = NewDB[Country](sourceURL("ipinfo-country", "ipv4"), 8, 12*time.Hour)
	IpinfoCountryIPv6      = NewDB[Country](sourceURL("ipinfo-country", "ipv6"), 8, 12*time.Hour)
	DbipCountryIPv4        = NewDB[Country](sourceURL("dbip-country", "ipv4"), 8, 15*24*time.Hour)
	DbipCountryIPv6        = NewDB[Country](sourceURL("dbip-country", "ipv6"), 8, 15*24*time.Hour)
	Ip2locationCountryIPv4 = NewDB[Country](sourceURL("ip2location-country", "ipv4"), 8, 30*24*time.Hour)
	Ip2locationCountryIPv6 = NewDB[Country](sourceURL("ip2location-country", "ipv6"), 8, 30*24*time.Hour)
	Geolite2CountryIPv4    = NewDB[Country](sourceURL("geolite2-country", "ipv4"), 8, 4*24*time.Hour)
	Geolite2CountryIPv6    = NewDB[Country](sourceURL("geolite2-country", "ipv6"), 8, 4*24*time.Hour)

	DbipCityIPv4         = NewDB[City1](sourceURL("dbip-city", "ipv4"), 12, 15*24*time.Hour)
	DbipCityIPv6         = NewDB[City1](sourceURL("dbip-city", "ipv6"), 12, 15*24*time.Hour)
	Ip2locationCityIPv4  = NewDB[City1](sourceURL("ip2location-city", "ipv4"), 12, 30*24*time.Hour)
	Ip2locationCityIPv6  = NewDB[City1](sourceURL("ip2location-city", "ipv6"), 12, 30*24*time.Hour)
	Geolite2CityIPv4     = NewDB[City2](sourceURL("geolite2-city", "ipv4-en"), 12, 4*24*time.Hour)
	Geolite2CityIPv6     = NewDB[City2](sourceURL("geolite2-city", "ipv6-en"), 12, 4*24*time.Hour)
	Geolite2CityIPv4DE   = NewDB[City2](sourceURL("geolite2-city", "ipv4-de"), 12, 4*24*time.Hour)
	Geolite2CityIPv6DE   = NewDB[City2](sourceURL("geolite2-city", "ipv6-de"), 12, 4*24*time.Hour)
	Geolite2CityIPv4ES   = NewDB[City2](sourceURL("geolite2-city", "ipv4-es"), 12, 4*24*time.Hour)
	Geolite2CityIPv6ES   = NewDB[City2](sourceURL("geolite2-city", "ipv6-es"), 12, 4*24*time.Hour)
	Geolite2CityIPv4FR   = NewDB[City2](sourceURL("geolite2-city", "ipv4-fr"), 12, 4*24*time.Hour)
	Geolite2CityIPv6FR   = NewDB[City2](sourceURL("geolite2-city", "ipv6-fr"), 12, 4*24*time.Hour)
	Geolite2CityIPv4JA   = NewDB[City2](sourceURL("geolite2-city", "ipv4-ja"), 12, 4*24*time.Hour)
	Geolite2CityIPv6JA   = NewDB[City2](sourceURL("geolite2-city", "ipv6-ja"), 12, 4*24*time.Hour)
	Geolite2CityIPv4PTBR = NewDB[City2](sourceURL("geolite2-city", "ipv4-pt-BR"), 12, 4*24*time.Hour)
	Geolite2CityIPv6PTBR = NewDB[City2](sourceURL("geolite2-city", "ipv6-pt-BR"), 12, 4*24*time.Hour)
	Geolite2CityIPv4RU   = NewDB[City2](sourceURL("geolite2-city", "ipv4-ru"), 12, 4*24*time.Hour)
	Geolite2CityIPv6RU   = NewDB[City2](sourceURL("geolite2-city", "ipv6-ru"), 12, 4*24*time.Hour)
	Geolite2CityIPv4ZHCN = NewDB[City2](sourceURL("geolite2-city", "ipv4-zh-CN"), 12, 4*24*time.Hour)
	Geolite2CityIPv6ZHCN = NewDB[City2](sourceURL("geolite2-city", "ipv6-zh-CN"), 12, 4*24*time.Hour)
)

func NewDB[T any](dbPath string, cacheLevel int, updateInterval time.Duration) *iplocation.LazyDB[T] {
	return iplocation.MustLazyDB[T](dbPath, '\t', iplocation.HexRange, cacheLevel, updateInterval)
}

// sourceURL build the url of db file
// https://github.com/tdulcet/ip-geolocation-dbs
func sourceURL(source string, v string) string {
	var ext = "tsv"
	return fmt.Sprintf("https://gitlab.com/tdulcet/ip-geolocation-dbs/-/raw/main/%s/%s.%s",
		source, v, ext,
	)
}
