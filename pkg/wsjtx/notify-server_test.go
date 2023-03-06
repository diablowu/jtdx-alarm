package wsjtx

import (
	"fmt"
	"net"
	"testing"
)

func TestMakeServerGiven(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", "server.natappfree.cc", 57241))
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("addr: %v", addr)
	}
}

func TestParseIP(t *testing.T) {
}
