package iplocation

import (
	"fmt"
	"net/netip"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

type LazyDB[T any] struct {
	dbPath *url.URL

	sep        byte
	rangeKind  RangeKind
	cacheLevel int
	cachePath  string

	db              *DB[T]
	updating        bool
	updated         time.Time
	updatedInterval time.Duration

	mut sync.Mutex
}

func MustLazyDB[T any](dbPath string, sep byte, rangeKind RangeKind, cacheLevel int, updateInterval time.Duration) *LazyDB[T] {
	db, err := NewLazyDB[T](dbPath, sep, rangeKind, cacheLevel, updateInterval)
	if err != nil {
		panic(err)
	}
	return db
}

func NewLazyDB[T any](dbPath string, sep byte, rangeKind RangeKind, cacheLevel int, updateInterval time.Duration) (*LazyDB[T], error) {
	u, err := url.Parse(dbPath)
	if err != nil {
		return nil, err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &LazyDB[T]{
		dbPath:          u,
		sep:             sep,
		rangeKind:       rangeKind,
		cacheLevel:      cacheLevel,
		cachePath:       path.Join(homeDir, ".iplocation"),
		updatedInterval: updateInterval,
	}, nil
}

func (d *LazyDB[T]) updateFromRemote() {
	_, modTime, err := downloadOrCache(d.cachePath, d.dbPath, true)
	d.mut.Lock()
	defer d.mut.Unlock()
	if err == nil {
		_ = d.db.Reload()
		d.updated = modTime
	}
	d.updating = false
}

func (d *LazyDB[T]) init() error {
	d.mut.Lock()
	defer d.mut.Unlock()
	if d.updating {
		return nil
	}

	if !d.updated.IsZero() && time.Since(d.updated) < d.updatedInterval {
		return nil
	}

	dbPath := d.dbPath.Path
	switch d.dbPath.Scheme {
	case "http", "https":
		if d.db == nil {
			path, modTime, err := downloadOrCache(d.cachePath, d.dbPath, false)
			if err != nil {
				return fmt.Errorf("unable to download %s file: %w", d.dbPath, err)
			}
			db, err := NewDB[T](path, d.sep, d.rangeKind, d.cacheLevel)
			if err != nil {
				return err
			}
			d.db = db
			d.updated = modTime
		} else {

			d.updating = true
			go d.updateFromRemote()
		}
	default:
		if d.db == nil {
			db, err := NewDB[T](dbPath, d.sep, d.rangeKind, d.cacheLevel)
			if err != nil {
				return err
			}
			d.db = db
		}
		fi, err := os.Stat(dbPath)
		if err == nil {
			modTime := fi.ModTime()
			if d.updated.IsZero() {
				d.updated = modTime
			} else {
				if d.updated.Before(modTime) {
					return d.db.Reload()
				}
			}
		}
	}
	return nil
}

func (d *LazyDB[T]) Lookup(addr netip.Addr) (*T, error) {
	err := d.init()
	if err != nil {
		return nil, err
	}

	return d.db.Lookup(addr)
}
