package kafka

import (
	"github.com/Shopify/sarama"
	"testing"
)

//测试kafka连接
func TestInitKafkaClient(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "10.231.50.30",
			args: []string{"10.231.50.30:5463"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := InitKafkaClient(tt.args)
			if err != nil {
				t.Errorf("InitKafkaClient() error = %v", err)
				return
			}
			message := sarama.ProducerMessage{
				Topic: "default",
				Value: sarama.StringEncoder("测试消息"),
			}
			_, _, err = gotClient.Client.SendMessage(&message)
			if err != nil {
				t.Errorf("SendMessage error = %v", err)
				return
			}
		})
	}
}
