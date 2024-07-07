package iplocation

import (
	"net/netip"
)

type DB[T any] struct {
	db            *RawDB
	bytesToStruct func(b []byte, v any) error
}

func NewDB[T any](dbPath string, sep byte, rangeKind RangeKind, cacheLevel int) (*DB[T], error) {
	rdb, err := NewRawDB(dbPath, sep, rangeKind, cacheLevel)
	if err != nil {
		return nil, err
	}
	return &DB[T]{
		db:            rdb,
		bytesToStruct: bytesToStruct(sep),
	}, nil
}

func (d *DB[T]) Lookup(addr netip.Addr) (*T, error) {
	data, err := d.db.Lookup(addr)
	if err != nil {
		return nil, err
	}

	var t T
	err = d.bytesToStruct(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (d *DB[T]) Reload() error {
	return d.db.Reload()
}
