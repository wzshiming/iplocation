package inject

import (
	"bytes"
	"fmt"
	"io"
	"net/netip"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	data := []byte(`111.111.111.111
0011:2233:4455:6677:8899:aabb:ccdd:eeff
1.1.1.1`)
	netip.IPv6Loopback()
	got, _ := io.ReadAll(NewReader(bytes.NewReader(data), func(addr netip.Addr) ([]byte, bool) {
		return []byte(fmt.Sprintf(">%s<", addr.String())), true
	}))
	want := `>111.111.111.111<
>11:2233:4455:6677:8899:aabb:ccdd:eeff<
>1.1.1.1<`

	if !strings.EqualFold(string(got), want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
