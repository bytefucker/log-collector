package kafka

import (
	"logTransfer/elasticsearch"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

// kafka消费者
type KafkaClient struct {
	client sarama.Consumer
	// topic列表
	topics []string
}

var (
	// 正在运行的topic任务
	runTopics   []string
	kafkaClient *KafkaClient
)

// 初始化kafka消费者
func InitKafka(kafkaAddress []string, topics []string) (err error) {
	consumer, err := sarama.NewConsumer(kafkaAddress, nil)
	if err != nil {
		return
	}
	kafkaClient = &KafkaClient{
		client: consumer,
		topics: topics,
	}
	// 遍历topics列表创建topic任务
	for _, topic := range topics {
		go createTopicTask(topic)
	}
	return
}

// 判断指定topic是否在运行
func isTopicExists(topic string) (status bool) {
	status = false
	for _, t := range runTopics {
		if t == topic {
			status = true
			break
		}
	}
	return
}

// 创建kafka topic消费者任务
func createTopicTask(topic string) {
	// 等待1秒创建任务
	time.Sleep(time.Second)
	// 判断指定topic是否在运行
	if isTopicExists(topic) == true {
		return
	}
	var wg sync.WaitGroup
	logs.Info("create topic [%s] kafka consumer", topic)
	// 开始消费kafka topic
	partitionList, err := kafkaClient.client.Partitions(topic)
	if err != nil {
		logs.Error("get topic: [%s] partitions failed, err: %s", topic, err)
		return
	}

	for partition := range partitionList {
		pc, err := kafkaClient.client.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logs.Warn("topic: [%s] start cousumer partition failed, err: %s", topic, err)
			continue
		}
		defer pc.AsyncClose()
		wg.Add(1)
		// 启动goroute去消费topic中的每一个块
		go func(pc sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				esMsg := elasticsearch.EsMsg{
					Msg:   msg.Value,
					Topic: topic,
				}
				// 获取到的数据发送到eschan
				elasticsearch.EsChan <- esMsg
			}
			wg.Done()
		}(pc)
	}

	// 启动成功的topic任务添加到运行列表中去
	runTopics = append(runTopics, topic)

	logs.Info("create topic [%s] kafka consumer success", topic)
	wg.Wait()
}

// 更新topic任务
func UpdateTopicTask(topics []string) (err error) {
	for _, newTopic := range topics {
		var topicStauts = false
		for _, oldTopic := range kafkaClient.topics {
			if newTopic == oldTopic {
				topicStauts = true
				break
			}
		}
		if topicStauts == false {
			go createTopicTask(newTopic)
		}
	}
	return
}
