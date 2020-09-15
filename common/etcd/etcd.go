package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

// etcd客户端对象
type EtcdClient struct {
	*clientv3.Client
}

// 初始化etcd
func NewClient(addrs []string) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	etcdClient := &EtcdClient{client}
	return etcdClient, err
}

// 初始化etcd key监控
/*func initEtcdWatch() {
	for _, key := range etcdClient.collectKeys {
		go etcdWatch(key)
	}
}*/

// 	etcd key监控处理
/*func etcdWatch(key string) {
	logs.Debug("start watch key: %s", key)
	for true {
		rech := etcdClient.Client.Watch(context.Background(), key)
		var colConfig []model.CollectTask
		var getConfStatus = true
		for wresp := range rech {
			for _, ev := range wresp.Events {
				// key 删除
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key [%s] is deleted", key)
					continue
				}
				// key 更新
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &colConfig)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%s", key, err)
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
		logs.Info("Update task config")
		// 更新tailf任务
		err := task.UpdateTailfTask(colConfig)
		if err != nil {
			logs.Error("Update task task failed, connect: %s, err: %s", colConfig, err)
		}
	}
}
*/
