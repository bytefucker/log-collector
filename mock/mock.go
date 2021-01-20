package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/yihongzhi/log-collector/agent/task"
	"github.com/yihongzhi/log-collector/common/utils"

	"github.com/yihongzhi/log-collector/common/etcd"
)

func mockTask() {
	var err error
	var key = "/log-collector/10.231.50.28/"
	etcdClient, err := etcd.NewClient([]string{"10.231.50.28:5460"})

	deleteResponse, err := etcdClient.Delete(context.TODO(), key, clientv3.WithPrefix())
	fmt.Println(deleteResponse, err)

	ops := []clientv3.Op{
		clientv3.OpPut(key+"vsr-camera", taskData("vsr-camera", "/var/log/vsr/vsr-camera/application.log")),
		clientv3.OpPut(key+"vsr-offline-task", taskData("vsr-offline-task", "/var/log/vsr/vsr-offline-task/application.log")),
		clientv3.OpPut(key+"vsr-plantask", taskData("vsr-plantask", "/var/log/vsr/vsr-plantask/application.log")),
	}
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

func mockLocalTask() {
	var err error
	var key = "/log-collector/" + utils.LocalIpArray[0] + "/"
	etcdClient, err := etcd.NewClient([]string{"10.231.50.28:5460"})

	deleteResponse, err := etcdClient.Delete(context.TODO(), key, clientv3.WithPrefix())
	fmt.Println(deleteResponse, err)

	ops := []clientv3.Op{
		clientv3.OpPut(key+"demo", taskData("demo", "/Users/yihongzhi/Downloads/megvii.SJ01202101-XXX-11-1234vsr-vsr-police_net-biz.biz-biz.vsr_offline_task-1/current.log")),
	}
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

func taskData(name string, path string) string {
	details := task.TailTaskDetails{
		AppKey:  name,
		LogPath: path,
	}
	json, _ := json.Marshal(details)
	return string(json)
}
