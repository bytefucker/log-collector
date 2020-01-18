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
			name:  "10.231.50.30",
			addrs: []string{"10.231.50.30:5460"},
			key:   "/logagent/10.231.22.76",
			value: "{\"appKey\":\"DemoApp\",\"logPath\":\"/Users/yihongzhi/Desktop/megvii.YL01201912-GDC-42-1223dj-dj-police_net-biz.biz-biz.sng_biz_config-1/current.log\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEtcdClient, err := InitEtcdClient(tt.addrs)
			if err != nil {
				t.Errorf("InitEtcdClient() error = %v", err)
				return
			}
			kv := clientv3.NewKV(gotEtcdClient.Client)

			_, err = kv.Put(context.TODO(), tt.key, tt.value)

			if err != nil {
				t.Errorf("Put error = %v", err)
				return
			}

			getResponse, err := kv.Get(context.Background(), tt.key)

			if err != nil {
				t.Errorf("Get error = %v", err)
				return
			}

			keyValue := getResponse.Kvs[0]

			if tt.key != string(keyValue.Key) || tt.value != string(keyValue.Value) {
				t.Errorf("Get error key=%s value=%s", keyValue.Key, keyValue.Value)
			}

			/*_, err = kv.Delete(context.TODO(), tt.key)

			if err != nil {
				t.Errorf("Delete error = %v", err)
				return
			}*/

		})
	}
}
