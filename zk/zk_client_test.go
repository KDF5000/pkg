package zk

import (
	"testing"
)

func TestZK(t *testing.T) {
	var addr []string
	for i := 0; i < 1; i++ {
		addr = append(addr, "127.0.0.1:2181")
	}

	cli, err := NewZKClient(addr)
	if err != nil {
		t.Fatal(err)
	}

	// [ms_topo_config_ rootservers_ cluster_uuid chunkservers_ cs_topo_config_ deployed cs_ids_config_]
	ret, _, err := cli.Children("/4")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", ret)

	data, stat, err := cli.Get("/hello")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("data: %s, stat: %+v", data, stat)
}
