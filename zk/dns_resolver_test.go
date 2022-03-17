package zk

import "testing"

func TestDNSHostProvider(t *testing.T) {
	rs := &DNSHostProvider{}
	err := rs.Init([]string{"zk.test.com:2128"})
	if err != nil {
		t.Fatal(err)
	}

	if rs.Len() < 0 {
		t.Fatalf("len: %d", rs.Len())
	}

	addr, retryStart := rs.Next()
	for !retryStart {
		t.Logf("%s", addr)
		addr, retryStart = rs.Next()
	}
}
