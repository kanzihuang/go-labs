package main

import (
	"context"
	"flag"
	etcdcli "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func main() {
	const ping, pong = "/test/ping", "pong"
	var address string
	flag.StringVar(&address, "address", "localhost:2379", "address")
	flag.Parse()
	log.Println("Connecting to", address)
	cli, err := etcdcli.New(etcdcli.Config{Endpoints: []string{address}})
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, ping)
	if err != nil {
		log.Fatal("读取数据失败：", err)
	}
	for _, kv := range resp.Kvs {
		if string(kv.Key) != ping || string(kv.Value) != pong {
			log.Fatalf("got %s from etcd, want %s", kv.Value, pong)
		}
	}
	<-ctx.Done()
}
