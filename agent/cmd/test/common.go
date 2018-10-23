package test

import (
	"net"
	"strconv"
	"strings"
	"testing"
)

func AddrToPort(addr net.Addr) (int32, error) {
	b := strings.Split(addr.String(), ":")

	p, err := strconv.ParseInt(b[1], 10, 32)
	if err != nil {
		return -1, err
	}
	return int32(p), nil
}

type convert func() error

func WaitFor(s convert, t chan error) {
	t <- s()
}

func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	t.Fatalf("%s, %v != %v", message, a, b)
}

func AssertContain(t *testing.T, s string, substr string, message string) {
	if strings.Contains(s, substr) {
		return
	}
	t.Fatalf("%s, %v does not include %v", message, s, substr)
}
