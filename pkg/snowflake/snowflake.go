//package snowflake
//
//import "github.com/bwmarrin/snowflake"
//
//var node *snowflake.Node
//
//func Init() error {
//	// 单机固定 0 号节点，先跑通功能
//	n, err := snowflake.NewNode(0)
//	if err != nil {
//		return err
//	}
//	node = n
//	return nil
//}
//
//func Generate() string {
//	return node.Generate().Base58()
//}

package snowflake

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var node *snowflake.Node

func Init() error {
	// 连接到 Etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"host.docker.internal:2379"},
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("etcd dial: %w", err)
	}
	defer cli.Close()

	// 用 hostname + 时间戳抢占机器 ID 0-1023
	key := fmt.Sprintf("/shorturl/node/%s", "local")
	resp, err := cli.Grant(context.Background(), 30)
	if err != nil {
		return fmt.Errorf("etcd grant: %w", err)
	}
	leaseID := resp.ID

	// 尝试从 0 开始抢
	for id := int64(0); id < 1024; id++ {
		k := fmt.Sprintf("%s/%d", key, id)
		tx := cli.Txn(context.Background()).
			If(clientv3.Compare(clientv3.Version(k), "=", 0)).
			Then(clientv3.OpPut(k, "1", clientv3.WithLease(leaseID)))
		if _, err := tx.Commit(); err == nil {
			// 抢到 ID，初始化 Snowflake
			node, _ = snowflake.NewNode(id)
			return nil
		}
	}
	return fmt.Errorf("no available node id")
}

func Generate() string {
	return node.Generate().Base58()
}
