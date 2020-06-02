package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"testing"
)

func TestInitEtcdClient(t *testing.T) {

	tests := []struct {
		name  string
		addrs []string
		key   string
		value string
	}{
		{
			name:  "10.231.50.26",
			addrs: []string{"10.231.50.30:22379"},
			key:   "/log/10.231.23.85",
			value: "{ \"appKey\":\"demo\",\"logPath\":\"/Users/yihongzhi/log/current.log\" }",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.addrs)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)
				return
			}
			kv := clientv3.NewKV(client.Client)

			_, err = kv.Put(context.TODO(), tt.key, tt.value)
			if err != nil {
				t.Errorf("Put error = %v", err)
				return
			}

			getResponse, err := kv.Get(context.TODO(), tt.key)
			if err != nil {
				t.Errorf("Get error = %v", err)
				return
			}
			keyValue := getResponse.Kvs[0]
			if tt.key != string(keyValue.Key) || tt.value != string(keyValue.Value) {
				t.Errorf("Get error key=%s value=%s", keyValue.Key, keyValue.Value)
			} else {
				t.Logf("Get key=%s value=%s", keyValue.Key, keyValue.Value)
			}

			/*deleteResponse, err := kv.Delete(context.TODO(), tt.key)
			if err != nil {
				t.Errorf("Get error = %v", err)
				return
			}
			if deleteResponse.Deleted != 1 {
				t.Errorf("delete error key=%s value=%s", tt.key, tt.value)
			} else {
				t.Logf("delete key=%s success count %d", tt.key, deleteResponse.Deleted)
			}*/

		})
	}
}
