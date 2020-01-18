package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"testing"
)

func TestInitEtcdClient(t *testing.T) {

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "10.231.50.30",
			args: []string{"10.231.50.30:5460"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEtcdClient, err := InitEtcdClient(tt.args)
			if err != nil {
				t.Errorf("InitEtcdClient() error = %v", err)
				return
			}
			kv := clientv3.NewKV(gotEtcdClient.Client)

			var demoKey = "/demo"
			var demoValue = "demo"

			_, err = kv.Put(context.TODO(), demoKey, demoValue)

			if err != nil {
				t.Errorf("Put error = %v", err)
				return
			}

			getResponse, err := kv.Get(context.Background(), demoKey)

			if err != nil {
				t.Errorf("Get error = %v", err)
				return
			}

			keyValue := getResponse.Kvs[0]

			if demoKey != string(keyValue.Key) || demoValue != string(keyValue.Value) {
				t.Errorf("Get error key=%s value=%s", keyValue.Key, keyValue.Value)
			}

		})
	}
}
