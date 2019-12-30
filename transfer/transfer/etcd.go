package main

import (
	"context"
	"encoding/json"
	"logTransfer/kafka"
	"time"

	"github.com/astaxie/beego/logs"

	etcd "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// etcd客户端对象
type EtcdClient struct {
	client *etcd.Client
}

const (
	// etcd中topic任务key值
	topicsName = "topics"
)

var (
	// etcd连接对象
	etcdClient *EtcdClient
)

// 初始化etcd
func initEtcd(etcdAddress []string) (err error) {
	client, err := etcd.New(etcd.Config{
		Endpoints:   etcdAddress,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return
	}
	etcdClient = &EtcdClient{
		client: client,
	}
	// 获取etcd中topics列表
	err = getTopics()
	if err != nil {
		return
	}

	// start watch kafka topics from etcd
	go etcdWatch(topicsName)
	return
}

// 获取etcd中topics列表
func getTopics() (err error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.client.Get(ctx, topicsName)
	cancle()
	if err != nil {
		return
	}
	for _, v := range resp.Kvs {
		if string(v.Key) == topicsName {
			err = json.Unmarshal(v.Value, &transferConfig.Topics)
			if err != nil {
				return
			}
		}
	}
	return
}

// 监听etcd topics key的变化
func etcdWatch(key string) {
	logs.Debug("start watch topics")
	for true {
		rech := etcdClient.client.Watch(context.Background(), key)
		var topics []string
		var getConfStatus = true
		for wresp := range rech {
			for _, ev := range wresp.Events {
				// 删除
				if ev.Type == mvccpb.DELETE {
					logs.Error("topics list is deleted")
					continue
				}
				// 更新
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &topics)
					if err != nil {
						logs.Error("topics, Unmarshal failed, err:%s", err)
						getConfStatus = false
						continue
					}
				}
				logs.Debug("get etcd config, %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfStatus {
				break
			}
		}
		logs.Info("Update topics config")
		// 更新topic任务
		err := kafka.UpdateTopicTask(topics)
		if err != nil {
			logs.Error("Update kafka sonsumer task failed, connect: %s, err: %s", topics, err)
		}
	}
}
