package iplocation

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/netip"
	"reflect"
	"strconv"
	"unsafe"
)

func keySepAfter2(s byte) func(a []byte) ([]byte, []byte, bool) {
	return func(a []byte) ([]byte, []byte, bool) {
		i := bytes.IndexByte(a, s)
		if i == -1 {
			return nil, nil, false
		}
		j := bytes.IndexByte(a[i+1:], s)
		if j == -1 {
			return nil, nil, false
		}
		return a[:i+j+1], a[i+j+2:], true
	}
}

func anyNumberCompare(a, b []byte) int {
	if len(a) != len(b) {
		if len(a) < len(b) {
			return -1
		}
		return 1
	}

	return bytes.Compare(a, b)
}

func cmpInBetween(s byte) func(a, b []byte) int {
	return func(a, b []byte) int {
		i := bytes.IndexByte(b, s)
		if i == -1 {
			return -1
		}

		ac := anyNumberCompare(a, b[:i])
		if ac == -1 {
			return -1
		}

		bc := anyNumberCompare(a, b[i+1:])
		if bc == 1 {
			return 1
		}

		return 0
	}
}

func ip2NumberDec(addr netip.Addr) ([]byte, error) {
	bi := big.NewInt(0)
	bi.SetBytes(addr.AsSlice())
	return strings2Bytes(bi.String()), nil
}

func ip2NumberHex(addr netip.Addr) ([]byte, error) {
	return strings2Bytes(hex.EncodeToString(addr.AsSlice())), nil
}

func bytesToStruct(s byte) func(b []byte, v any) error {
	return func(b []byte, v any) error {
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Ptr {
			return errors.New("not a pointer")
		}

		val = val.Elem()

		n := val.NumField()

		items := bytes.SplitN(b, []byte{s}, n)
		for i, item := range items {
			field := val.Field(i)
			if field.Kind() != reflect.String {
				return fmt.Errorf("field %d not a string", i)
			}
			if !field.CanSet() {
				return fmt.Errorf("cannot set field %d", i)
			}
			field.SetString(unquote(bytes2String(item)))
		}
		return nil
	}
}

func unquote(s string) string {
	if len(s) < 2 || s[0] != '"' {
		return s
	}

	ns, err := strconv.Unquote(s)
	if err != nil {
		return s
	}
	return ns
}

func bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func strings2Bytes(b string) []byte {
	return *(*[]byte)(unsafe.Pointer(&b))
}
