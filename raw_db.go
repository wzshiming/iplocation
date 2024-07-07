package iplocation

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/wzshiming/bsbf"
)

var ErrNotFound = errors.New("not found")

type RawDB struct {
	db      *bsbf.BSBF
	ipToKey func(addr netip.Addr) ([]byte, error)
}

type RangeKind uint

const (
	_ RangeKind = iota
	// CIDR TODO: not supported
	CIDR
	// RawRange TODO: not supported
	RawRange

	DecRange
	HexRange
)

func NewRawDB(path string, sep byte, rangeKind RangeKind, cacheLevel int) (*RawDB, error) {
	var ipToKey func(addr netip.Addr) ([]byte, error)
	switch rangeKind {
	default:
		return nil, fmt.Errorf("unknown IP range kind %d", rangeKind)
	case DecRange:
		ipToKey = ip2NumberDec
	case HexRange:
		ipToKey = ip2NumberHex
	}
	bs := bsbf.NewBSBF(path,
		bsbf.WithKeySepFunc(keySepAfter2(sep)),
		bsbf.WithCmpFunc(cmpInBetween(sep)),
		bsbf.WithCacheLevel(cacheLevel),
	)
	d := &RawDB{
		db:      bs,
		ipToKey: ipToKey,
	}
	return d, nil
}

func (d *RawDB) Lookup(addr netip.Addr) ([]byte, error) {
	ip, err := d.ipToKey(addr)
	if err != nil {
		return nil, err
	}

	_, _, got, ok, err := d.db.Search(ip)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%w: ip %s", ErrNotFound, addr)
	}
	return got, nil
}

func (d *RawDB) Reload() error {
	return d.db.Reload()
}
