package mock

import (
	"context"

	"fmt"
	"github.com/coreos/etcd/clientv3"

	"github.com/yihongzhi/log-collector/common/etcd"
)

func mockTask() {
	var err error
	var key = "/log-collector/10.231.20.75"
	etcdClient, err := etcd.NewClient([]string{"10.231.50.28:5460"})

	deleteResponse, err := etcdClient.Delete(context.TODO(), key)
	fmt.Println(deleteResponse, err)

	ops := []clientv3.Op{
		clientv3.OpPut(key+"/1", "1"),
		clientv3.OpPut(key+"/2", "2"),
		clientv3.OpPut(key+"/3", "3")}
	for _, op := range ops {
		if _, err := etcdClient.Do(context.TODO(), op); err != nil {
			fmt.Println(err)
		}
	}
	getRsp, err := etcdClient.Get(context.TODO(), key, clientv3.WithPrefix())
	fmt.Println(getRsp, err)
	for _, kv := range getRsp.Kvs {
		fmt.Println(string(kv.Key), ":", string(kv.Value))
	}
}
