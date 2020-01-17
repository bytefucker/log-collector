package etcd

import (
	"context"
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
			_, err = gotEtcdClient.Client.KV.Get(context.TODO(), "/demo")
		})
	}
}
