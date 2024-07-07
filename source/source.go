package source

import (
	"net/netip"

	"github.com/wzshiming/iplocation"
	"github.com/wzshiming/iplocation/source/sapics"
	"github.com/wzshiming/iplocation/source/tdulcet"
)

type Source interface {
	Lookup(addr netip.Addr) (any, error)
}

type sourceWrapper[T any] struct {
	ipv4 *iplocation.LazyDB[T]
	ipv6 *iplocation.LazyDB[T]
}

func newSourceWrapper[T any](ipv4, ipv6 *iplocation.LazyDB[T]) Source {
	return &sourceWrapper[T]{
		ipv4: ipv4,
		ipv6: ipv6,
	}
}

func (s sourceWrapper[T]) Lookup(addr netip.Addr) (any, error) {
	if addr.Is4() {
		return s.ipv4.Lookup(addr)
	}
	return s.ipv6.Lookup(addr)
}

var Sources = map[string]Source{
	"tdulcet-geowhoisasn-country": newSourceWrapper(tdulcet.GeoWhoisAsnCountryIPv4, tdulcet.GeoWhoisAsnCountryIPv6),
	"tdulcet-iptoasn-country":     newSourceWrapper(tdulcet.IptoasnCountryIPv4, tdulcet.IptoasnCountryIPv6),
	"tdulcet-ipinfo-country":      newSourceWrapper(tdulcet.IpinfoCountryIPv4, tdulcet.IpinfoCountryIPv6),

	"tdulcet-dbip-country": newSourceWrapper(tdulcet.DbipCountryIPv4, tdulcet.DbipCountryIPv6),
	"tdulcet-dbip-city":    newSourceWrapper(tdulcet.DbipCityIPv4, tdulcet.DbipCityIPv6),

	"tdulcet-ip2location-country": newSourceWrapper(tdulcet.Ip2locationCountryIPv4, tdulcet.Ip2locationCountryIPv6),
	"tdulcet-ip2location-city":    newSourceWrapper(tdulcet.Ip2locationCityIPv4, tdulcet.Ip2locationCityIPv6),

	"tdulcet-geolite2-country": newSourceWrapper(tdulcet.Geolite2CountryIPv4, tdulcet.Geolite2CountryIPv6),
	"tdulcet-geolite2-city":    newSourceWrapper(tdulcet.Geolite2CityIPv4, tdulcet.Geolite2CityIPv6),

	"sapics-routeviewsasndbip-asn": newSourceWrapper(sapics.RouteViewsAsnDbipASNIPv4, sapics.RouteViewsAsnDbipASNIPv6),
	"sapics-asn-country":           newSourceWrapper(sapics.AsnCountryIPv4, sapics.AsnCountryIPv6),
	"sapics-geoasn-country":        newSourceWrapper(sapics.GeoAsnCountryIPv4, sapics.GeoAsnCountryIPv6),
	"sapics-geowhoisasn-country":   newSourceWrapper(sapics.GeoWhoisAsnCountryIPv4, sapics.GeoWhoisAsnCountryIPv6),

	"sapics-iptoasn-asn":     newSourceWrapper(sapics.IptoasnASNIPv4, sapics.IptoasnASNIPv6),
	"sapics-iptoasn-country": newSourceWrapper(sapics.IptoasnCountryIPv4, sapics.IptoasnCountryIPv6),

	"sapics-geolite2-asn":     newSourceWrapper(sapics.GeoLite2ASNIPv4, sapics.GeoLite2ASNIPv6),
	"sapics-geolite2-country": newSourceWrapper(sapics.GeoLite2CountryIPv4, sapics.GeoLite2CountryIPv6),
	"sapics-geolite2-city":    newSourceWrapper(sapics.GeoLite2CityIPv4, sapics.GeoLite2CityIPv6),

	"sapics-dbip-asn":     newSourceWrapper(sapics.DbipASNIPv4, sapics.DbipASNIPv6),
	"sapics-dbip-country": newSourceWrapper(sapics.DbipCountryIPv4, sapics.DbipCountryIPv6),
	"sapics-dbip-city":    newSourceWrapper(sapics.DbipCityIPv4, sapics.DbipCityIPv6),
}
