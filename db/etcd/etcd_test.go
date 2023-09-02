package main

import (
	"context"
	etcdcli "go.etcd.io/etcd/client/v3"
	"testing"
)

func TestEtcdGet(t *testing.T) {
	const ping, pong = "/test/ping", "pong"
	cli, err := etcdcli.New(etcdcli.Config{Endpoints: []string{"192.168.1.189:2379"}})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	resp, err := cli.Get(ctx, ping)
	if err != nil {
		t.Error("读取数据失败：", err)
	}
	for _, kv := range resp.Kvs {
		if string(kv.Key) != ping || string(kv.Value) != pong {
			t.Errorf("got %s from etcd, want %s", kv.Value, pong)
		}
	}
}
