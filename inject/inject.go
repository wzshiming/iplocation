package inject

import (
	"bufio"
	"io"
	"net/netip"
	"regexp"
)

var (
	ip = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|([a-fA-F0-9]{0,4}:){1,7}[a-fA-F0-9]{0,4}`)
)

type ReplaceIPFunc func(addr netip.Addr) ([]byte, bool)

func replace(reg *regexp.Regexp, src []byte, repl ReplaceIPFunc) []byte {
	return reg.ReplaceAllFunc(src, func(bytes []byte) []byte {
		addr, err := netip.ParseAddr(string(bytes))
		if err != nil {
			return bytes
		}
		n, ok := repl(addr)
		if !ok {
			return bytes
		}
		return n
	})
}

func Replace(src []byte, repl ReplaceIPFunc) []byte {
	src = replace(ip, src, repl)
	return src
}

func NewReader(r io.Reader, repl ReplaceIPFunc) io.Reader {
	return &reader{
		repl: repl,
		r:    bufio.NewReader(r),
	}
}

type reader struct {
	repl ReplaceIPFunc
	r    *bufio.Reader
	buf  []byte
	err  error
}

func (r *reader) Read(p []byte) (n int, err error) {
	if len(r.buf) == 0 {
		if r.err != nil {
			return 0, r.err
		}
		line, err := r.r.ReadSlice('\n')
		r.err = err
		if len(line) == 0 {
			return 0, err
		}
		r.buf = Replace(line, r.repl)
	}

	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}
